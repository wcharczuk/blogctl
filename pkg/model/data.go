package model

// Data is the site's entire database of posts.
type Data struct {
	Title   string     `json:"title" yaml:"title"`
	Author  string     `json:"author" yaml:"author"`
	BaseURL string     `json:"baseURl" yaml:"baseURL"`
	Posts   []Post     `json:"posts" yaml:"posts"`
	Tags    []TagPosts `json:"tags" yaml:"tags"`
}

// IsZero returns if the object is set.
func (d Data) IsZero() bool {
	return len(d.Posts) == 0
}
