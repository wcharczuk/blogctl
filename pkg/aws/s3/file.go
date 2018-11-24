package s3

// File is info for a file upload.
type File struct {
	ETag                 string
	Contents             []byte
	ACL                  string
	FilePath             string
	Bucket               string
	Key                  string
	ContentType          string
	ContentDisposition   string
	ServerSideEncryption string
}

// IsZero returns if the file is set or not.
func (f File) IsZero() bool {
	return len(f.Key) == 0
}

// ACLOrDefault returns the file ACL or a default.
func (f File) ACLOrDefault() string {
	if f.ACL != "" {
		return f.ACL
	}
	return ACLPrivate
}
