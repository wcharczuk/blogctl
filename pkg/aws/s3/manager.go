package s3

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/blend/go-sdk/async"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"

	"github.com/wcharczuk/blogctl/pkg/aws"
)

// New returns a new manager.
func New(cfg *aws.Config) *Manager {
	return &Manager{
		Config: cfg,
		Ignores: []string{
			".DS_Store",
			".git",
		},
		Session:     aws.NewSession(cfg),
		Parallelism: runtime.NumCPU(),
	}
}

// Manager is a helper for uploading files to s3.
type Manager struct {
	Log               logger.Log
	Ignores           []string
	Config            *aws.Config
	Session           *session.Session
	PutObjectDefaults File
	DryRun            bool
	Parallelism       int
}

// ParallelismOrDefault returns the parallelism or a default.
func (m Manager) ParallelismOrDefault() int {
	if m.Parallelism > 0 {
		return m.Parallelism
	}
	return runtime.NumCPU()
}

// GetKey returns the relative path for a given file.
func (m Manager) GetKey(rootPath, workingPath string) string {
	if !strings.HasPrefix(workingPath, "./") {
		workingPath = "./" + workingPath
	}
	if !strings.HasPrefix(rootPath, "./") {
		rootPath = "./" + rootPath
	}
	return strings.TrimPrefix(workingPath, rootPath)
}

// SyncDirectory sync's a directory.
// It returns a list of invalidated keys (i.e. keys to update or remove), and an error.
func (m Manager) SyncDirectory(ctx context.Context, directoryPath, bucket string) (invalidations []string, err error) {
	if m.DryRun {
		m.Log.Debugf("sync directory (dry run): not realizing changes")
	}
	localFiles := make(chan interface{}, 1024)
	if err := m.DiscoverFiles(ctx, localFiles, directoryPath); err != nil {
		return nil, err
	}
	invalidations, err = m.ProcessFiles(ctx, localFiles, directoryPath, bucket)
	return
}

// DiscoverFiles discovers local files.
func (m Manager) DiscoverFiles(ctx context.Context, localFiles chan interface{}, directoryPath string) (err error) {
	err = filepath.Walk(directoryPath, func(currentPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == directoryPath {
			return nil
		}
		for _, ignore := range m.Ignores {
			if strings.HasSuffix(currentPath, ignore) {
				return nil
			}
		}
		if fileInfo.IsDir() {
			return nil
		}
		localFiles <- currentPath
		return nil
	})
	return
}

// ProcessFiles processes the files list.
func (m Manager) ProcessFiles(ctx context.Context, localFiles chan interface{}, directoryPath, bucket string) (invalidated []string, err error) {
	remoteETags := make(map[string]string)
	localKeys := new(Set)

	remoteFiles, err := m.List(ctx, bucket)
	if err != nil {
		return nil, err
	}

	remoteFileBatch := make(chan interface{}, len(remoteFiles))
	for _, remoteFile := range remoteFiles {
		key := remoteFile.Key
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}
		remoteETags[key] = aws.StripQuotes(remoteFile.ETag)
		remoteFileBatch <- remoteFile
	}

	errors := make(chan error, len(localFiles))

	// create an async batch to process the file list.
	async.NewBatch(func(ctx context.Context, workItem interface{}) error {
		file, fileOK := workItem.(string)
		if !fileOK {
			return ex.New("process files; batch work item was not a string")
		}

		key := m.GetKey(directoryPath, file)
		localKeys.Set(key)

		var localETag string
		remoteETag, hasRemoteFile := remoteETags[key]
		if hasRemoteFile { // if we need to compare against a remote etag
			localETag, err = m.GenerateETag(file)
			if err != nil {
				return err
			}
		}

		if !hasRemoteFile || remoteETag != localETag {
			logger.MaybeInfof(m.Log, "putting %s", key)

			contentType, err := webutil.DetectContentType(file)
			if err != nil {
				return err
			}

			if !m.DryRun {
				if err := m.Put(ctx, File{
					FilePath:    file,
					Key:         key,
					Bucket:      bucket,
					ContentType: contentType,
				}); err != nil {
					return err
				}
			}
			if hasRemoteFile {
				invalidated = append(invalidated, key)
			}
		} else {
			logger.MaybeInfof(m.Log, "skipping %s (unchanged)", key)
		}
		return nil
	}, localFiles, async.OptBatchParallelism(m.ParallelismOrDefault()), async.OptBatchErrors(errors)).Process(ctx)

	// print errors if any were produced by the batch.
	if errorCount := len(errors); errorCount > 0 {
		for x := 0; x < errorCount; x++ {
			logger.MaybeError(m.Log, <-errors)
		}
		return nil, ex.New("process files; issues sending files to s3")
	}

	var invalidatedSync sync.Mutex
	async.NewBatch(func(ctx context.Context, workItem interface{}) error {
		remoteFile, remoteFileOK := workItem.(File)
		if !remoteFileOK {
			return ex.New("process files; remote cleanup batch work item was not a file")
		}

		key := remoteFile.Key
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}

		if !localKeys.Has(key) {
			logger.MaybeInfof(m.Log, "removing remote %s", remoteFile.Key)
			if !m.DryRun {
				if err := m.Delete(ctx, bucket, remoteFile.Key); err != nil {
					return err
				}
			}

			invalidatedSync.Lock()
			invalidated = append(invalidated, key)
			invalidatedSync.Unlock()
		}
		return nil
	}, remoteFileBatch, async.OptBatchParallelism(m.ParallelismOrDefault()), async.OptBatchErrors(errors))

	// print errors if any were produced by the batch.
	if errorCount := len(errors); errorCount > 0 {
		for x := 0; x < errorCount; x++ {
			logger.MaybeError(m.Log, <-errors)
		}
		return nil, ex.New("process files; issues removing files from s3")
	}

	return invalidated, nil
}

// List lists all files in a bucket.
func (m Manager) List(ctx context.Context, bucket string) ([]File, error) {
	remoteFiles, err := s3.New(m.Session).ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: &bucket,
	})
	if IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var files []File
	for _, file := range remoteFiles.Contents {
		files = append(files, File{
			Bucket: bucket,
			Key:    aws.DerefStr(file.Key),
			ETag:   aws.DerefStr(file.ETag),
		})
	}
	return files, nil
}

// Get fetches a file at a given key
func (m Manager) Get(ctx context.Context, bucket, key string) (file File, contents io.ReadCloser, err error) {
	remoteFile, getErr := s3.New(m.Session).GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if IsNotFound(getErr) {
		return
	}
	if getErr != nil {
		err = ex.New(getErr)
		return
	}

	file = File{
		Bucket:               bucket,
		Key:                  key,
		ContentType:          aws.DerefStr(remoteFile.ContentType),
		ContentDisposition:   aws.DerefStr(remoteFile.ContentDisposition),
		ServerSideEncryption: aws.DerefStr(remoteFile.ServerSideEncryption),
		ETag:                 aws.DerefStr(remoteFile.ETag),
	}
	contents = remoteFile.Body
	return
}

// GetMeta fetches file metadata at a given key
func (m Manager) GetMeta(ctx context.Context, bucket, key string) (meta File, err error) {
	var contents io.ReadCloser
	meta, contents, err = m.Get(ctx, bucket, key)
	if err != nil {
		return
	}
	if contents != nil {
		defer contents.Close()
	}
	return
}

// Put uploads a file to s3.
func (m Manager) Put(ctx context.Context, fileInfo File) error {
	var size int64
	var contentType, contentDisposition, acl, serverSideEncryption string
	var contents io.ReadSeeker

	if fileInfo.FilePath != "" {
		file, err := os.Open(fileInfo.FilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		stats, err := file.Stat()
		if err != nil {
			return err
		}
		contents = file
		size = int64(stats.Size())
	} else if len(fileInfo.Contents) > 0 {
		size = int64(len(fileInfo.Contents))
		contents = bytes.NewReader(fileInfo.Contents)
	} else {
		return ex.New("invalid put object").WithMessage("must set either the path or the contents")
	}

	if fileInfo.ContentType != "" {
		contentType = fileInfo.ContentType
	} else if m.PutObjectDefaults.ContentDisposition != "" {
		contentType = m.PutObjectDefaults.ContentType
	}

	if fileInfo.ContentDisposition != "" {
		contentDisposition = fileInfo.ContentDisposition
	} else if m.PutObjectDefaults.ContentDisposition != "" {
		contentDisposition = m.PutObjectDefaults.ContentDisposition
	}

	if fileInfo.ACL != "" {
		acl = fileInfo.ACL
	} else if m.PutObjectDefaults.ACL != "" {
		acl = m.PutObjectDefaults.ACL
	}

	if fileInfo.ServerSideEncryption != "" {
		serverSideEncryption = fileInfo.ServerSideEncryption
	} else if m.PutObjectDefaults.ServerSideEncryption != "" {
		serverSideEncryption = m.PutObjectDefaults.ServerSideEncryption
	}

	_, err := s3.New(m.Session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.RefStr(fileInfo.Bucket),
		Key:                  aws.RefStr(fileInfo.Key),
		Body:                 contents,
		ContentLength:        &size,
		ContentType:          aws.RefStr(contentType),
		ContentDisposition:   aws.RefStr(contentDisposition),
		ACL:                  aws.RefStr(acl),
		ServerSideEncryption: aws.RefStr(serverSideEncryption),
	})
	return ex.New(err)
}

// Delete removes an object with a given key.
func (m Manager) Delete(ctx context.Context, bucket, key string) error {
	_, err := s3.New(m.Session).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.RefStr(bucket),
		Key:    aws.RefStr(key),
	})
	return ex.New(err)
}

// GenerateETag generate an etag for a give file by path.
func (m Manager) GenerateETag(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", ex.New(err)
	}

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", ex.New(err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
