package engine

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/model"
)

func TestEngineCreateSlugDefaults(t *testing.T) {
	assert := assert.New(t)

	defaults := config.Config{}
	engine := &Engine{Config: defaults}
	slugTemplate, err := engine.ParseSlugTemplate()
	assert.Nil(err)

	post := model.Post{
		Meta: model.Meta{
			Title:  "test slug",
			Posted: time.Date(2018, 12, 11, 10, 9, 8, 7, time.UTC),
		},
	}
	assert.Equal("2018/12/11/test-slug", engine.CreateSlug(slugTemplate, post))
}
