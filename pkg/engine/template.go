package engine

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/blend/go-sdk/exception"
	sdkTemplate "github.com/blend/go-sdk/template"
	"github.com/wcharczuk/blogctl/pkg/model"
)

// ViewFuncs is the template funcs.
func ViewFuncs() template.FuncMap {
	base := sdkTemplate.Funcs.FuncMap()
	base["partition"] = partition
	base["set_title"] = setTitle
	base["render"] = render
	return base
}

// ParseTemplate creates a new template from a string
func ParseTemplate(literal string) (*template.Template, error) {
	tmp := template.New("")
	tmp.Funcs(ViewFuncs())
	return tmp.Parse(literal)
}

// RenderString renders a template to a string for a given viewmodel.
func RenderString(tmp *template.Template, vm interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	if err := tmp.Execute(buffer, vm); err != nil {
		return "", exception.New(err)
	}
	return buffer.String(), nil
}

// Partition Errors
const (
	ErrPartitionCountTooLarge exception.Class = "partition count greater than number of posts"
	ErrPartitionIndexTooLarge exception.Class = "partition index greater than number of partitions"
	ErrPartitionCountInvalid  exception.Class = "partition count invalid; must be greater than 1"
)

func partition(index, partitions int, posts []*model.Post) ([]*model.Post, error) {
	if partitions < 1 {
		return nil, exception.New(ErrPartitionCountInvalid)
	}
	if index < 0 || index >= partitions {
		return nil, exception.New(ErrPartitionIndexTooLarge)
	}
	if partitions == 1 {
		return posts, nil
	}

	var output []*model.Post
	for ; index < len(posts); index += partitions {
		output = append(output, posts[index])
	}
	return output, nil
}

func setTitle(vm *model.ViewModel, title string) error {
	if vm == nil {
		return fmt.Errorf("viewmodel unset")
	}
	vm.Title = title
	return nil
}

func render(post model.Post) (template.HTML, error) {
	if post.Template == nil {
		return "", fmt.Errorf("post has unset template; cannot render. post: %s", post.TitleOrDefault())
	}
	buffer := new(bytes.Buffer)
	err := post.Template.Execute(buffer, post)
	if err != nil {
		return "", err
	}
	return template.HTML(buffer.String()), nil
}
