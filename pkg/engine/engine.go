package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/wcharczuk/photoblog/pkg/fileutil"
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

// InitializeOutputPath creates the output path if it doesn't exist.
func (e Engine) InitializeOutputPath() error {
	if fileutil.Exists(e.Config.OutputPathOrDefault()) {
		if err := exception.New(os.RemoveAll(e.Config.OutputPathOrDefault())); err != nil {
			return err
		}
	}
	return MakeDir(e.Config.OutputPathOrDefault())
}

// GetExistingData gets the existing data file.
func (e Engine) GetExistingData() (*model.Data, error) {
	dataFilePath := filepath.Join(e.Config.OutputPathOrDefault(), constants.FileData)
	if !fileutil.Exists(dataFilePath) {
		return nil, nil
	}

	f, err := os.Open(dataFilePath)
	if err != nil {
		return nil, exception.New(err)
	}
	defer f.Close()

	var data model.Data
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, exception.New(err)
	}
	return &data, nil
}

// GenerateData generates the blog data.
func (e Engine) GenerateData() (*model.Data, error) {
	output := model.Data{
		Title:   e.Config.TitleOrDefault(),
		Author:  e.Config.AuthorOrDefault(),
		BaseURL: e.Config.BaseURLOrDefault(),
	}
	imagesPath := e.Config.PostsPathOrDefault()

	logger.MaybeSyncInfof(e.Log, "searching `%s` for images as posts", imagesPath)
	err := filepath.Walk(imagesPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == imagesPath {
			return nil
		}
		if info.IsDir() {
			logger.MaybeSyncInfof(e.Log, "reading post `%s`", currentPath)
			post, err := e.ReadImage(currentPath)
			if err != nil {
				return err
			}
			output.Posts = append([]model.Post{*post}, output.Posts...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for index := range output.Posts {
		output.Posts[index].Index = index
		if index > 0 {
			output.Posts[index].Previous = &output.Posts[index-1]
		}
		if index < len(output.Posts)-1 {
			output.Posts[index].Next = &output.Posts[index+1]
		}
	}

	return &output, nil
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
		if name == constants.FileMeta {
			if err := ReadYAML(filepath.Join(path, name), &post.Meta); err != nil {
				return nil, err
			}
		} else if HasExtension(name, constants.ImageExtensions...) && post.Image.IsZero() {
			post.OriginalPath = filepath.Join(path, name)
			modTime = fi.ModTime()
			if post.Image, err = ReadImage(post.OriginalPath); err != nil {
				return nil, err
			}
		}
	}
	post.Slug = model.CreateSlug(post)
	if post.Meta.Posted.IsZero() {
		post.Meta.Posted = modTime
	}
	if post.OriginalPath == "" {
		return nil, exception.New("no image found").WithMessage(path)
	}
	return &post, nil
}

// ReadPartials reads all the partials named in the config.
func (e Engine) ReadPartials() ([]string, error) {
	partialsPath := e.Config.PartialsPathOrDefault()

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
func (e Engine) Render(data *model.Data) error {
	outputPath := e.Config.OutputPathOrDefault()

	partials, err := e.ReadPartials()
	if err != nil {
		return err
	}

	pagesPath := e.Config.PagesPathOrDefault()
	pages, err := ListDirectory(pagesPath)
	if err != nil {
		return err
	}
	for index, page := range pages {
		logger.MaybeSyncInfof(e.Log, "rendering page `%s`", page.Name())
		pageTemplate, err := e.CompileTemplate(filepath.Join(pagesPath, page.Name()), partials)
		if err != nil {
			return err
		}
		pageOutputPath := filepath.Join(outputPath, page.Name())
		if err := e.WriteTemplate(pageTemplate, pageOutputPath, ViewModel{
			Config:    e.Config,
			PostIndex: index,
			Post:      model.Posts(data.Posts).First(),
			Posts:     data.Posts,
		}); err != nil {
			return err
		}
	}

	postTemplatePath := e.Config.PostTemplateOrDefault()
	postTemplate, err := e.CompileTemplate(postTemplatePath, partials)
	if err != nil {
		return err
	}

	// foreach post, render the post with single to <slug>/index.html
	for index, post := range data.Posts {
		slugPath := filepath.Join(outputPath, post.Slug)
		logger.MaybeSyncInfof(e.Log, "rendering post `%s`", post.TitleOrDefault())

		if err := MakeDir(slugPath); err != nil {
			return exception.New(err)
		}
		if err := e.WriteTemplate(postTemplate, filepath.Join(slugPath, constants.FileIndex), ViewModel{
			Config:    e.Config,
			Posts:     data.Posts,
			Post:      post,
			PostIndex: index,
		}); err != nil {
			return err
		}
		if err := e.GenerateThumbnails(post.OriginalPath, slugPath); err != nil {
			return err
		}
	}

	staticPath := e.Config.StaticPathOrDefault()
	if err := Copy(staticPath, outputPath); err != nil {
		return err
	}

	if err := e.WriteData(data, filepath.Join(outputPath, constants.FileData)); err != nil {
		return err
	}

	return nil
}

// ShouldGenerateThumbnails returns if we should parse the original file
// so that we can generate thumbnails.
func (e Engine) ShouldGenerateThumbnails(originalPath, destinationPath string) bool {
	existsOriginal := fileutil.Exists(originalPath)
	includeOriginal := e.Config.IncludeOriginalOrDefault()

	if includeOriginal && !existsOriginal {
		return true
	}

	for _, size := range e.Config.ImageSizesOrDefault() {
		imageFilepath := filepath.Join(destinationPath, fmt.Sprintf(constants.ImageSizeFormat, size))
		if !fileutil.Exists(imageFilepath) {
			return true
		}
	}
	return false
}

// GenerateThumbnails generates our main thumbnails for the image.
func (e Engine) GenerateThumbnails(originalPath, destinationPath string) error {
	if !e.ShouldGenerateThumbnails(originalPath, destinationPath) {
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

	existsOriginal := fileutil.Exists(originalPath)
	includeOriginal := e.Config.IncludeOriginalOrDefault()

	if includeOriginal && !existsOriginal {
		logger.MaybeSyncInfof(e.Log, "copying post `%s` original", originalPath)
		if err := WriteFile(filepath.Join(destinationPath, constants.ImageOriginal), originalContents); err != nil {
			return err
		}
	}

	for _, size := range e.Config.ImageSizesOrDefault() {
		imageFilepath := filepath.Join(destinationPath, fmt.Sprintf(constants.ImageSizeFormat, size))
		if fileutil.Exists(imageFilepath) {
			continue
		}
		logger.MaybeSyncInfof(e.Log, "resizing post `%s` %dpx", originalPath, size)
		if err := e.Resize(original, imageFilepath, uint(size)); err != nil {
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

// WriteData writes a data file to disk.
func (e Engine) WriteData(data *model.Data, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return exception.New(err)
	}
	defer f.Close()
	return exception.New(json.NewEncoder(f).Encode(data))
}

// Generate generates the blog to the given output directory.
func (e Engine) Generate() error {
	current, err := e.GenerateData()
	if err != nil {
		return err
	}

	if err := e.InitializeOutputPath(); err != nil {
		return err
	}

	if err := e.Render(current); err != nil {
		return err
	}

	return nil
}
