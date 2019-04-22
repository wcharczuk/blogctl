package model

import "time"

// Data is the site's entire database of posts.
type Data struct {
	Title   string  `json:"title" yaml:"title"`
	Author  string  `json:"author" yaml:"author"`
	BaseURL string  `json:"baseURl" yaml:"baseURL"`
	Posts   []*Post `json:"posts" yaml:"posts"`
	Tags    []Tag   `json:"tags" yaml:"tags"`
}

// IsZero returns if the object is set.
func (d Data) IsZero() bool {
	return len(d.Posts) == 0
}

// NumPosts returns the total number of posts.
func (d Data) NumPosts() (count int) {
	count = len(d.Posts)
	return
}

// NumTags returns the total number of tags.
func (d Data) NumTags() (count int) {
	count = len(d.Tags)
	return
}

// NumImagePosts returns the number of image posts.
func (d Data) NumImagePosts() (count int) {
	for _, post := range d.Posts {
		if post.IsImage() {
			count = count + 1
		}
	}
	return
}

// NumTextPosts returns the number of text posts.
func (d Data) NumTextPosts() (count int) {
	for _, post := range d.Posts {
		if !post.IsImage() {
			count = count + 1
		}
	}
	return
}

// EarliestPost returns the earliest post.
func (d Data) EarliestPost() (date time.Time) {
	if len(d.Posts) == 0 {
		return
	}
	date = d.Posts[0].Meta.Posted

	for _, post := range d.Posts {
		if post.Meta.Posted.Before(date) {
			date = post.Meta.Posted
		}
	}
	return
}

// LatestPost returns the latest post.
func (d Data) LatestPost() (date time.Time) {
	if len(d.Posts) == 0 {
		return
	}
	date = d.Posts[0].Meta.Posted

	for _, post := range d.Posts {
		if post.Meta.Posted.After(date) {
			date = post.Meta.Posted
		}
	}
	return
}
