package engine

import (
	"bytes"
	"html/template"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	sdkTemplate "github.com/blend/go-sdk/template"

	"github.com/wcharczuk/photoblog/pkg/config"
	"github.com/wcharczuk/photoblog/pkg/constants"
	"github.com/wcharczuk/photoblog/pkg/model"
	"github.com/wcharczuk/photoblog/pkg/resize"
)

// New returns a new engine..
func New(cfg config.Config) *Engine {
	return &Engine{
		Config: cfg,
	}
}

// Engine returns a
type Engine struct {
	Config config.Config
	Log    *logger.Logger
}

// WithLogger sets the logger (optional).
func (e *Engine) WithLogger(log *logger.Logger) *Engine {
	e.Log = log
	return e
}

// CreateOutputPath creates the output path if it doesn't exist.
func (e Engine) CreateOutputPath() error {
	if _, err := os.Stat(e.Config.OutputOrDefault()); err != nil {
		return MakeDir(e.Config.OutputOrDefault())
	}
	return nil
}

// DiscoverPosts finds posts and returns an array of posts.
func (e Engine) DiscoverPosts() ([]model.Post, error) {
	imagesPath := e.Config.ImagesOrDefault()

	logger.MaybeSyncInfof(e.Log, "searching `%s` for images as posts", imagesPath)

	var posts []model.Post
	err := filepath.Walk(imagesPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == imagesPath {
			return nil
		}
		if info.IsDir() {
			logger.MaybeSyncInfof(e.Log, "reading `%s` as post", currentPath)
			post, err := e.ReadImage(currentPath)
			if err != nil {
				return err
			}
			posts = append([]model.Post{*post}, posts...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for index := range posts {
		if index > 0 {
			posts[index].Previous = &posts[index-1]
		}
		if index < len(posts)-1 {
			posts[index].Next = &posts[index+1]
		}
	}

	return posts, nil
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
	files, err := ListDirectory(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, exception.New("no child files found").WithMessage(path)
	}

	var post model.Post
	var modTime time.Time
	for _, fi := range files {
		name := fi.Name()
		if name == constants.DiscoveryFileMeta {
			if err := ReadYAML(filepath.Join(path, name), &post.Meta); err != nil {
				return nil, err
			}
		} else if HasExtension(name, constants.ImageExtensions...) && post.Image.IsZero() {
			post.Original = filepath.Join(path, name)
			modTime = fi.ModTime()
			if post.Image, err = ReadImage(post.Original); err != nil {
				return nil, err
			}
		}
	}
	if post.Meta.Posted.IsZero() {
		post.Meta.Posted = modTime
	}
	if post.Original == "" {
		return nil, exception.New("no images found").WithMessage(path)
	}
	return &post, nil
}

// ReadPartials reads all the partials named in the config.
func (e Engine) ReadPartials() ([]string, error) {
	partialsPath := e.Config.Layout.PartialsOrDefault()

	partialFiles, err := ListDirectory(partialsPath)
	if err != nil {
		return nil, err
	}

	var partials []string
	for _, partial := range partialFiles {
		contents, err := ioutil.ReadFile(filepath.Join(partialsPath, partial.Name()))
		if err != nil {
			return nil, exception.New(err)
		}
		partials = append(partials, string(contents))
	}
	return partials, nil
}

// CompileTemplate compiles a template.
func (e Engine) CompileTemplate(templatePath string, partials []string) (*template.Template, error) {
	contents, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, exception.New(err)
	}

	vf := sdkTemplate.ViewFuncs{}.FuncMap()
	tmp := template.New("").Funcs(vf)
	for _, partial := range partials {
		_, err := tmp.Parse(partial)
		if err != nil {
			return nil, exception.New(err)
		}
	}

	final, err := tmp.Parse(string(contents))
	if err != nil {
		return nil, exception.New(err)
	}
	return final, nil
}

// Render writes the templates out for each of the posts.
func (e Engine) Render(posts ...model.Post) error {
	outputPath := e.Config.OutputOrDefault()

	partials, err := e.ReadPartials()
	if err != nil {
		return err
	}

	pagesPath := e.Config.Layout.Pages
	pages, err := ListDirectory(pagesPath)
	if err != nil {
		return err
	}
	for _, page := range pages {
		logger.MaybeSyncInfof(e.Log, "rendering page %s", page.Name())
		pageTemplate, err := e.CompileTemplate(filepath.Join(pagesPath, page.Name()), partials)
		if err != nil {
			return err
		}
		pageOutputPath := filepath.Join(outputPath, page.Name())
		if err := e.WriteTemplate(pageTemplate, pageOutputPath, ViewModel{
			Config: e.Config,
			Posts:  posts,
		}); err != nil {
			return err
		}
	}

	postTemplatePath := e.Config.Layout.PostOrDefault()
	postTemplate, err := e.CompileTemplate(postTemplatePath, partials)
	if err != nil {
		return err
	}

	// foreach post, render the post with single to <slug>/index.html
	for _, post := range posts {
		slugPath := filepath.Join(outputPath, post.Slug())

		logger.MaybeSyncInfof(e.Log, "rendering post %s", post.TitleOrDefault())

		// make the slug directory tree (i.e. `mkdir -p <slug>`)
		if err := MakeDir(slugPath); err != nil {
			return exception.New(err)
		}
		if err := e.WriteTemplate(postTemplate, filepath.Join(slugPath, constants.OutputFileIndex), ViewModel{
			Config: e.Config,
			Posts:  posts,
			Post:   post,
		}); err != nil {
			return err
		}

		if err := e.GenerateThumbnails(post.Original, slugPath); err != nil {
			return err
		}
	}

	return nil
}

// FileExists returns if a given file exists.
func (e Engine) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GenerateThumbnails generates our main thumbnails for the image.
func (e Engine) GenerateThumbnails(originalPath, destinationPath string) error {
	filepath2048 := filepath.Join(destinationPath, constants.Image2048)
	filepath1024 := filepath.Join(destinationPath, constants.Image1024)
	filepath512 := filepath.Join(destinationPath, constants.Image512)
	existsOriginal := e.FileExists(originalPath)
	includeOriginal := e.Config.IncludeOriginalOrDefault()
	exists2048 := e.FileExists(filepath2048)
	exists1024 := e.FileExists(filepath1024)
	exists512 := e.FileExists(filepath512)

	if (!includeOriginal || (includeOriginal && existsOriginal)) && exists2048 && exists1024 && exists512 {
		return nil
	}

	originalContents, err := ioutil.ReadFile(originalPath)
	if err != nil {
		return exception.New(err)
	}

	// decode jpeg into image.Image
	original, err := jpeg.Decode(bytes.NewBuffer(originalContents))
	if err != nil {
		exception.New(err)
	}

	if includeOriginal && !existsOriginal {
		logger.MaybeSyncInfof(e.Log, "copying post %s original", originalPath)
		if err := WriteFile(filepath.Join(destinationPath, constants.ImageOriginal), originalContents); err != nil {
			return err
		}
	}

	if !exists2048 {
		logger.MaybeSyncInfof(e.Log, "resizing post %s 2048px", originalPath)
		if err := e.Resize(original, filepath2048, 2048); err != nil {
			return err
		}
	}

	if !exists1024 {
		logger.MaybeSyncInfof(e.Log, "resizing post %s 1024px", originalPath)
		if err := e.Resize(original, filepath1024, 1024); err != nil {
			return err
		}
	}

	if !exists512 {
		logger.MaybeSyncInfof(e.Log, "resizing post %s 512px", originalPath)
		if err := e.Resize(original, filepath1024, 512); err != nil {
			return err
		}
	}
	return nil
}

// Resize resizes an image to a destination.
func (e Engine) Resize(original image.Image, destination string, maxDimension uint) error {
	resized := resize.Thumbnail(maxDimension, maxDimension, original, resize.Bicubic)
	out, err := os.Create(destination)
	if err != nil {
		return exception.New(err)
	}
	defer out.Close()

	// write new image to file
	if err := jpeg.Encode(out, resized, nil); err != nil {
		return exception.New(err)
	}
	return nil
}

// WriteTemplate writes a template to a given path with a given data viewmodel.
func (e Engine) WriteTemplate(tpl *template.Template, outputPath string, data interface{}) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return exception.New(err)
	}
	defer f.Close()
	if err := tpl.Execute(f, data); err != nil {
		return exception.New(err)
	}
	return nil
}

// CopyStatics copies static files to the output directory.
func (e Engine) CopyStatics() error {
	outputPath := e.Config.OutputOrDefault()
	// copy statics (things like css or js)
	staticPath := e.Config.Layout.StaticOrDefault()
	if err := Copy(staticPath, outputPath); err != nil {
		return err
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

	logger.MaybeSyncInfof(e.Log, "discovered %d posts", len(posts))

	// create the output path if it doesn't exist
	if err := e.CreateOutputPath(); err != nil {
		return err
	}

	// render templates
	if err := e.Render(posts...); err != nil {
		return err
	}

	// copy statics.
	if err := e.CopyStatics(); err != nil {
		return err
	}

	return nil
}
