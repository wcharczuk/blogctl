package engine

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/photoblog/pkg/config"
	"github.com/wcharczuk/photoblog/pkg/constants"
	"github.com/wcharczuk/photoblog/pkg/model"
)

// New returns a new engine..
func New(cfg config.Config) Engine {
	return Engine{
		Config: cfg,
	}
}

// Engine returns a
type Engine struct {
	Log    *logger.Logger
	Config config.Config
}

// CreateOutputPath creates the output path if it doesn't exist.
func (e Engine) CreateOutputPath() error {
	if _, err := os.Stat(e.Config.GetOutput()); err != nil {
		return exception.New(os.MkdirAll(e.Config.GetOutput(), 0666))
	}
	return nil
}

// DiscoverPosts finds posts and returns an array of posts.
func (e Engine) DiscoverPosts() ([]model.Post, error) {
	imagesPath := e.Config.GetImages()

	var posts []model.Post
	err := filepath.Walk(imagesPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			post, err := e.ReadImage(currentPath)
			if err != nil {
				return err
			}
			posts = append(posts, *post)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ReadImage reads post metadata from a folder.
func (e Engine) ReadImage(path string) (*model.Post, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, exception.New("not a directory").WithMessage(path)
	}

	// sniff image file
	// and metadata
	files, err := GetFileInfos(path)
	if err != nil {
		return nil, err
	}

	var post model.Post
	for _, fi := range files {
		name := fi.Name()

		if name == constants.DiscoveryFileMetadata {
			if err := ReadYAML(filepath.Join(path, name), &post.Meta); err != nil {
				return nil, err
			}
		} else if HasExtension(name, constants.ImageExtensions...) && post.Image.IsZero() {
			if err := ReadImage(name, &post.Image); err != nil {
				return nil, err
			}
		}
	}

	return &post, nil
}

// CompileTemplates compiles the templates.
func (e Engine) CompileTemplates() (home, single *template.Template, err error) {
	indexPath := e.Config.Layout.GetIndex()
	homePath := e.Config.Layout.GetHome()
	singlePath := e.Config.Layout.GetSingle()

	home, err = template.New("home").ParseFiles(append([]string{homePath, indexPath}, e.Config.Layout.GetPartials()...)...)
	if err != nil {
		return
	}
	single, err = template.New("single").ParseFiles(append([]string{singlePath, indexPath}, e.Config.Layout.GetPartials()...)...)
	if err != nil {
		return
	}
	return
}

// Render writes the templates out for each of the posts.
func (e Engine) Render(home, single *template.Template, posts ...model.Post) error {
	outputPath := e.Config.GetOutput()

	// render home as "index.html"
	homePath := filepath.Join(outputPath, constants.FileIndex)
	if err := e.WriteTemplate(home, homePath, HomeViewModel{Config: e.Config, Posts: posts}); err != nil {
		return err
	}

	// foreach post, render the post with single to <slug>/index.html
	for index, post := range posts {
		// make the slug directory tree (i.e. `mkdir -p <slug>`)
		if err := os.MkdirAll(post.Slug(), 0666); err != nil {
			return exception.New(err)
		}

		var next, previous model.Post
		if index > 0 {
			previous = posts[index-1]
		}
		if index < len(posts)-1 {
			next = posts[index+1]
		}

		if err := e.WriteTemplate(single, filepath.Join(post.Slug(), constants.FileIndex), SingleViewModel{
			Config:   e.Config,
			Post:     post,
			Previous: previous,
			Next:     next,
		}); err != nil {
			return err
		}

		if err := Copy(post.ImagePath, filepath.Join(post.Slug(), filepath.Base(post.ImagePath))); err != nil {
			return err
		}
	}

	return nil
}

// WriteTemplate writes a template to a given path with a given data viewmodel.
func (e Engine) WriteTemplate(tpl *template.Template, outputPath string, data interface{}) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := tpl.Execute(f, data); err != nil {
		return exception.New(err)
	}
	return nil
}

// CopyStatics copies static files to the output directory.
func (e Engine) CopyStatics() error {
	outputPath := e.Config.GetOutput()

	// copy statics (things like css or js)
	staticPaths := e.Config.Layout.GetStatics()
	for _, staticPath := range staticPaths {
		if err := Copy(staticPath, outputPath); err != nil {
			return err
		}
	}
	return nil
}

// Generate generates the blog to the given output directory.
func (e Engine) Generate() error {
	// discover posts
	posts, err := e.DiscoverPosts()
	if err != nil {
		return err
	}

	// parse template(s)
	list, single, err := e.CompileTemplates()
	if err != nil {
		return err
	}

	// create the output path if it doesn't exist
	if err := e.CreateOutputPath(); err != nil {
		return err
	}

	// render templates
	if err := e.Render(list, single, posts...); err != nil {
		return err
	}

	// copy statics.
	if err := e.CopyStatics(); err != nil {
		return err
	}
	return nil
}

//
// internal helper methods
//

func (e Engine) infof(format string, args ...interface{}) {
	if e.Log == nil {
		return
	}
	e.Log.Infof(format, args...)
}

func (e Engine) warning(err error) {
	if e.Log == nil {
		return
	}
	e.Log.Warning(err)
}

func (e Engine) warningf(format string, args ...interface{}) {
	if e.Log == nil {
		return
	}
	e.Log.Warningf(format, args...)
}

func (e Engine) error(err error) {
	if e.Log == nil {
		return
	}
	e.Log.Error(err)
}

func (e Engine) errorf(format string, args ...interface{}) {
	if e.Log == nil {
		return
	}
	e.Log.Errorf(format, args...)
}
