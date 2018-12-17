package s3

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/blogctl/pkg/aws"
	"github.com/wcharczuk/blogctl/pkg/fileutil"
)

// New returns a new manager.
func New(cfg *aws.Config) *Manager {
	return &Manager{
		Config:  cfg,
		Session: aws.NewSession(cfg),
	}
}

// Manager is a helper for uploading files to s3.
type Manager struct {
	Log               *logger.Logger
	Config            *aws.Config
	Session           *session.Session
	PutObjectDefaults File
}

// GetKey returns the relative path for a given file.
func (m Manager) GetKey(rootPath, workingPath string) string {
	return strings.TrimPrefix(workingPath, rootPath)
}

// SyncDirectory sync's a directory.
// It returns a list of invalidated keys (i.e. keys to update or remove), and an error.
func (m Manager) SyncDirectory(ctx context.Context, directoryPath, bucket string) ([]string, error) {
	remoteETags := make(map[string]string)
	localKeys := make(map[string]bool)
	invalidated := []string{}

	// walk the s3 bucket, look for files that need to be removed ...
	remoteFiles, err := m.List(ctx, bucket)
	if err != nil {
		return nil, err
	}
	for _, remoteFile := range remoteFiles {
		key := remoteFile.Key
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}
		remoteETags[key] = aws.StripQuotes(remoteFile.ETag)
	}

	// walk the directory ...
	err = filepath.Walk(directoryPath, func(currentPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == directoryPath {
			return nil
		}
		if strings.HasSuffix(currentPath, ".DS_Store") {
			return nil
		}
		if fileInfo.IsDir() {
			return nil
		}

		key := m.GetKey(directoryPath, currentPath)

		localKeys[key] = true

		// process files and directories ...
		remoteETag, ok := remoteETags[key]

		var localETag string
		if ok {
			localETag, err = m.GenerateETag(currentPath)
			if err != nil {
				return err
			}
		}

		if !ok || remoteETag != localETag {
			logger.MaybeInfof(m.Log, "putting %s", key)

			contentType, err := fileutil.DetectContentType(currentPath)
			if err != nil {
				return err
			}

			if err := m.Put(ctx, File{
				FilePath:    currentPath,
				Key:         key,
				Bucket:      bucket,
				ContentType: contentType,
			}); err != nil {
				return err
			}
			// only invalidate files we know to be present (not new files)
			if ok {
				invalidated = append(invalidated, key)
			}
		} else {
			logger.MaybeInfof(m.Log, "skipping %s (unchanged)", key)
		}

		return nil
	})

	if err != nil {
		return nil, exception.New(err)
	}

	for _, remoteFile := range remoteFiles {
		key := remoteFile.Key
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}

		if _, ok := localKeys[key]; !ok {
			logger.MaybeInfof(m.Log, "removing remote %s", remoteFile.Key)
			if err := m.Delete(ctx, bucket, remoteFile.Key); err != nil {
				return nil, err
			}
			invalidated = append(invalidated, key)
		}
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
		err = exception.New(getErr)
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
		return exception.New("invalid put object").WithMessage("must set either the path or the contents")
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
	return exception.New(err)
}

// Delete removes an object with a given key.
func (m Manager) Delete(ctx context.Context, bucket, key string) error {
	_, err := s3.New(m.Session).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.RefStr(bucket),
		Key:    aws.RefStr(key),
	})
	return exception.New(err)
}

// GenerateETag generate an etag for a give file by path.
func (m Manager) GenerateETag(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", exception.New(err)
	}

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", exception.New(err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
