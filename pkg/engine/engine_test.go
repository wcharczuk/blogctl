package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/model"
)

func TestEngineCreateSlugDefaults(t *testing.T) {
	assert := assert.New(t)

	defaults := config.Config{}
	e := &Engine{Config: defaults}
	slugTemplate, err := e.ParseSlugTemplate()
	assert.Nil(err)

	post := model.Post{
		Meta: model.Meta{
			Title:  "test slug",
			Posted: time.Date(2018, 12, 11, 10, 9, 8, 7, time.UTC),
		},
	}
	assert.Equal("2018/12/11/test-slug", e.CreateSlug(slugTemplate, post))
}

func TestEngineBuild(t *testing.T) {
	assert := assert.New(t)

	config, path, err := ReadConfig(&config.PersistentFlags{ConfigPath: ref.String("./testdata/config.yml")})
	assert.Nil(err)
	assert.True(strings.HasSuffix(path, "testdata/config.yml"))
	assert.Nil(New(config).Generate(context.TODO()))
}
