package engine

import (
	"context"
	"os"
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

	post = model.Post{
		Meta: model.Meta{
			Title:  "Mt. Tam",
			Posted: time.Date(2018, 12, 11, 10, 9, 8, 7, time.UTC),
		},
	}
	assert.Equal("2018/12/11/mt-tam", e.CreateSlug(slugTemplate, post))
}

func TestEngineBuild(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(os.Chdir("testdata"))

	cfg, path, err := config.ReadConfig(config.Flags{
		ConfigPath:  ref.String("./config.yml"),
		Parallelism: ref.Int(4),
	})
	assert.Nil(err)
	assert.Equal("./config.yml", path)
	assert.Nil(MustNew(OptConfig(cfg)).Generate(context.TODO()))

	_, err = os.Stat("dist")
	assert.Nil(err)
	_, err = os.Stat("dist/2019/02/10/text-post")
	assert.Nil(err)
	_, err = os.Stat("dist/2019/02/11/image-post")
	assert.Nil(err)
	_, err = os.Stat("dist/2019/02/11/image-post/2048.jpg")
	assert.Nil(err)
	_, err = os.Stat("dist/2019/02/11/image-post/1024.jpg")
	assert.Nil(err)
	_, err = os.Stat("dist/2019/02/11/image-post/512.jpg")
	assert.Nil(err)
}
