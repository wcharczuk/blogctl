package web

// PostedFile is a file that has been posted to an hc endpoint.
type PostedFile struct {
	Key      string
	FileName string
	Contents []byte
}
