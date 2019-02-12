package engine

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/wcharczuk/blogctl/pkg/model"
)

func TestPartitionPosts(t *testing.T) {
	assert := assert.New(t)

	posts := []*model.Post{
		{Meta: model.Meta{Title: "one"}},
		{Meta: model.Meta{Title: "two"}},
		{Meta: model.Meta{Title: "three"}},
		{Meta: model.Meta{Title: "four"}},
		{Meta: model.Meta{Title: "five"}},
		{Meta: model.Meta{Title: "six"}},
		{Meta: model.Meta{Title: "seven"}},
	}

	partition0, err := partition(0, 2, posts)
	assert.Nil(err)
	partition1, err := partition(1, 2, posts)
	assert.Nil(err)

	assert.Len(partition0, 4)
	assert.Len(partition1, 3)

	assert.Equal("one", partition0[0].Meta.Title)
	assert.Equal("three", partition0[1].Meta.Title)
	assert.Equal("five", partition0[2].Meta.Title)
	assert.Equal("seven", partition0[3].Meta.Title)

	assert.Equal("two", partition1[0].Meta.Title)
	assert.Equal("four", partition1[1].Meta.Title)
	assert.Equal("six", partition1[2].Meta.Title)
}
