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
	"sort"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/constants"
	"github.com/wcharczuk/blogctl/pkg/fileutil"
	"github.com/wcharczuk/blogctl/pkg/model"
	"github.com/wcharczuk/blogctl/pkg/resize"
	"github.com/wcharczuk/blogctl/pkg/stringutil"
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
	Log    logger.FullReceiver
}

//
// properties
//

// WithLogger sets the logger (optional).
func (e *Engine) WithLogger(log logger.FullReceiver) *Engine {
	e.Log = log
	return e
}

// Generate generates the blog to the given output directory.
func (e Engine) Generate() error {
	if err := e.InitializeOutputPath(); err != nil {
		return err
	}

	if err := e.InitializeThumbnailCache(); err != nil {
		return err
	}

	partials, err := e.DiscoverPartials()
	if err != nil {
		return err
	}

	data, err := e.DiscoverPosts(partials...)
	if err != nil {
		return err
	}

	if err := e.Render(data, partials...); err != nil {
		return err
	}

	return nil
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

// InitializeThumbnailCache creates the output path if it doesn't exist.
func (e Engine) InitializeThumbnailCache() error {
	if fileutil.Exists(e.Config.ThumbnailCachePathOrDefault()) {
		return nil
	}
	return MakeDir(e.Config.ThumbnailCachePathOrDefault())
}

// DiscoverPosts generates the blog data.
func (e Engine) DiscoverPosts(partials ...string) (*model.Data, error) {
	slugTemplate, err := e.ParseSlugTemplate()
	if err != nil {
		return nil, err
	}

	output := model.Data{
		Title:   e.Config.TitleOrDefault(),
		Author:  e.Config.AuthorOrDefault(),
		BaseURL: e.Config.BaseURLOrDefault(),
	}
	tags := make(map[string]*model.Tag)
	imagesPath := e.Config.PostsPathOrDefault()

	logger.MaybeSyncInfof(e.Log, "searching `%s` for images as posts", imagesPath)
	err = filepath.Walk(imagesPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == imagesPath {
			return nil
		}
		if info.IsDir() {
			logger.MaybeSyncInfof(e.Log, "reading post `%s`", currentPath)

			// check if we have an image
			post, err := e.GeneratePost(slugTemplate, currentPath, partials...)
			if err != nil {
				return err
			}
			output.Posts = append([]model.Post{*post}, output.Posts...)

			if !e.Config.SkipTags {
				for _, tag := range post.Meta.Tags {
					if tagPosts, ok := tags[tag]; ok {
						tagPosts.Posts = append(tagPosts.Posts, *post)
					} else {
						tags[tag] = &model.Tag{
							Tag:   tag,
							Posts: []model.Post{*post},
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// sort by metadata posted date
	// we don't really care about directory / filesystem order
	sort.Sort(model.Posts(output.Posts))

	// create previous and next links for each post.
	for index := range output.Posts {
		if index > 0 {
			output.Posts[index].Previous = &output.Posts[index-1]
		}
		if index < len(output.Posts)-1 {
			output.Posts[index].Next = &output.Posts[index+1]
		}
	}

	if !e.Config.SkipTags {
		// add tags, make sure they're sorted.
		for _, tag := range tags {
			sort.Sort(model.Posts(tag.Posts))
			output.Tags = append(output.Tags, *tag)
		}
		sort.Sort(model.Tags(output.Tags))
	}

	return &output, nil
}

// Render writes the templates out for each of the posts.
func (e Engine) Render(data *model.Data, partials ...string) error {
	outputPath := e.Config.OutputPathOrDefault()

	pagesPath := e.Config.PagesPathOrDefault()
	pages, err := ListDirectory(pagesPath)
	if err != nil {
		return err
	}
	for _, page := range pages {
		logger.MaybeSyncInfof(e.Log, "rendering page `%s`", page.Name())
		pageTemplate, err := e.CompileTemplate(filepath.Join(pagesPath, page.Name()), partials)
		if err != nil {
			return err
		}
		pageOutputPath := filepath.Join(outputPath, page.Name())
		if err := e.WriteTemplate(pageTemplate, pageOutputPath, &model.ViewModel{
			Config: e.Config,
			Post:   model.Posts(data.Posts).First(),
			Posts:  data.Posts,
			Tags:   data.Tags,
		}); err != nil {
			return err
		}
	}

	imagePostTemplatePath := e.Config.ImagePostTemplateOrDefault()
	var imagePostTemplate *template.Template
	if imagePostTemplatePath != "" {
		imagePostTemplate, err = e.CompileTemplate(imagePostTemplatePath, partials)
		if err != nil {
			return err
		}
	}

	textPostTemplatePath := e.Config.TextPostTemplateOrDefault()
	var textPostTemplate *template.Template
	if textPostTemplatePath != "" {
		textPostTemplate, err = e.CompileTemplate(textPostTemplatePath, partials)
		if err != nil {
			return err
		}
	}

	// foreach post, render the post with single to <slug>/index.html
	for _, post := range data.Posts {
		slugPath := filepath.Join(outputPath, post.Slug)
		logger.MaybeSyncInfof(e.Log, "rendering post `%s`", post.TitleOrDefault())

		if err := MakeDir(slugPath); err != nil {
			return exception.New(err)
		}

		postTemplate := imagePostTemplate
		if post.Image.IsZero() {
			postTemplate = textPostTemplate
		}

		if err := e.WriteTemplate(postTemplate, filepath.Join(slugPath, constants.FileIndex), &model.ViewModel{
			Config: e.Config,
			Posts:  data.Posts,
			Tags:   data.Tags,
			Post:   post,
		}); err != nil {
			return err
		}

		if post.ImagePath != "" {
			if err := e.ProcessThumbnails(post.ImagePath, slugPath); err != nil {
				return err
			}
		}
	}

	if !e.Config.SkipTags {
		tagTemplatePath := e.Config.TagTemplateOrDefault()
		if len(tagTemplatePath) > 0 && fileutil.Exists(tagTemplatePath) {
			tagTemplate, err := e.CompileTemplate(tagTemplatePath, partials)
			if err != nil {
				return err
			}
			for _, tag := range data.Tags {
				tagPath := filepath.Join(outputPath, "tags", stringutil.Slugify(tag.Tag))
				if err := MakeDir(tagPath); err != nil {
					return exception.New(err)
				}
				if err := e.WriteTemplate(tagTemplate, filepath.Join(tagPath, constants.FileIndex), &model.ViewModel{
					Config: e.Config,
					Posts:  data.Posts,
					Tags:   data.Tags,
					Tag:    tag,
				}); err != nil {
					return err
				}
			}
		}
	}

	staticPath := e.Config.StaticsPathOrDefault()
	if err := Copy(staticPath, outputPath); err != nil {
		return err
	}

	if !e.Config.SkipJSONData {
		if err := e.WriteDataJSON(data, filepath.Join(outputPath, constants.FileData)); err != nil {
			return err
		}
	}

	return nil
}

//
// utilities
//

// DiscoverPartials reads all the partials named in the config.
func (e Engine) DiscoverPartials() ([]string, error) {
	partialsPath := e.Config.PartialsPathOrDefault()

	partialFiles, err := ListDirectory(partialsPath)
	if err != nil {
		return nil, err
	}

	var partials []string
	for _, partial := range partialFiles {
		contents, err := ioutil.ReadFile(filepath.Join(partialsPath, partial.Name()))
		if err != nil {
			return nil, exception.New(err).WithMessagef("partial: %s", partial)
		}
		partials = append(partials, string(contents))
	}
	return partials, nil
}

// GeneratePost reads post contents and metadata from a folder.
func (e Engine) GeneratePost(slugTemplate *template.Template, path string, partials ...string) (*model.Post, error) {
	// sniff for image file
	// sniff for template file
	// sniff for metadata
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
			post.ImagePath = filepath.Join(path, name)
			if modTime.Before(fi.ModTime()) {
				modTime = fi.ModTime()
			}
			if post.Image, err = ReadImage(post.ImagePath); err != nil {
				return nil, err
			}
		} else if HasExtension(name, constants.TemplateExtensions...) && post.TemplatePath == "" {
			post.TemplatePath = filepath.Join(path, name)
			if modTime.Before(fi.ModTime()) {
				modTime = fi.ModTime()
			}
			if post.Template, err = e.CompileTemplate(post.TemplatePath, partials); err != nil {
				return nil, err
			}
		}
	}
	post.Slug = e.CreateSlug(slugTemplate, post)
	if post.Meta.Posted.IsZero() {
		post.Meta.Posted = modTime
	}
	if post.ImagePath == "" && post.TemplatePath == "" {
		return nil, exception.New("no image or template found").WithMessage(path)
	}
	return &post, nil
}

// ProcessThumbnails processes thumbnails.
func (e Engine) ProcessThumbnails(originalFilePath, destinationPath string) error {
	originalContents, err := ioutil.ReadFile(originalFilePath)
	if err != nil {
		return exception.New(err)
	}

	etag, err := fileutil.ETag(originalContents)
	if err != nil {
		return err
	}

	if e.ShouldGenerateThumbnails(etag) {
		if err := e.GenerateThumbnails(originalContents, etag); err != nil {
			return nil
		}
	}

	if err := e.CopyThumbnails(etag, destinationPath); err != nil {
		return err
	}

	return nil
}

// CleanThumbnailCache cleans the thumbnail cache by purging cached thumbnails for posts that may have been deleted.
func (e Engine) CleanThumbnailCache(dryRun bool) error {
	// for each post, generate the sha of the image ...
	postsPath := e.Config.PostsPathOrDefault()
	logger.MaybeSyncInfof(e.Log, "searching `%s` for posts", postsPath)
	postSums := map[string]bool{}
	err := filepath.Walk(postsPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == postsPath {
			return nil
		}
		if info.IsDir() {
			files, err := ListDirectory(currentPath)
			if err != nil {
				return err
			}
			for _, fi := range files {
				name := fi.Name()
				if HasExtension(name, constants.ImageExtensions...) {
					contents, err := ioutil.ReadFile(filepath.Join(currentPath, name))
					if err != nil {
						return err
					}
					etag, err := fileutil.ETag(contents)
					if err != nil {
						return err
					}

					postSums[etag] = true
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// for each thumbnail cache folder
	var orphanedCachedPosts []string
	thumbnailCachePath := e.Config.ThumbnailCachePathOrDefault()
	err = filepath.Walk(thumbnailCachePath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == thumbnailCachePath {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			// see if there is a matching sha'd image
			if _, ok := postSums[name]; !ok {
				orphanedCachedPosts = append(orphanedCachedPosts, name)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// purge folders
	for _, path := range orphanedCachedPosts {
		if !dryRun {
			if err := os.RemoveAll(filepath.Join(thumbnailCachePath, path)); err != nil {
				return err
			}
		}
		logger.MaybeSyncInfof(e.Log, "purging orphaned cached directory `%s`", path)
	}

	return nil
}

// CompileTemplate compiles a template.
func (e Engine) CompileTemplate(templatePath string, partials []string) (*template.Template, error) {
	contents, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, exception.New(err).WithMessagef("template path: %s", templatePath)
	}

	tmp := template.New("").Funcs(ViewFuncs())
	for _, partial := range partials {
		_, err := tmp.Parse(partial)
		if err != nil {
			return nil, exception.New(err).WithMessagef("template path: %s", templatePath)
		}
	}

	final, err := tmp.Parse(string(contents))
	if err != nil {
		return nil, exception.New(err).WithMessagef("template path: %s", templatePath)
	}
	return final, nil
}

// WriteDataJSON writes a data file to disk.
func (e Engine) WriteDataJSON(data *model.Data, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return exception.New(err)
	}
	defer f.Close()
	return exception.New(json.NewEncoder(f).Encode(data))
}

// GenerateThumbnails generates and copies our main thumbnails for the post image.
// - originalFilePath should be the path to the original image file
// - destinationPath should be the path to the output slug folder
func (e Engine) GenerateThumbnails(originalContents []byte, etag string) error {
	// decode jpeg into image.Image
	original, err := jpeg.Decode(bytes.NewBuffer(originalContents))
	if err != nil {
		exception.New(err)
	}

	for _, size := range e.Config.ImageSizesOrDefault() {
		if err := e.GenerateThumbnail(original, size, etag); err != nil {
			return err
		}
	}
	return nil
}

// GenerateThumbnail generates a thumbnail and stores it in the cache if it doesn't exist
// and copies the cached thumbail to the output directory.
func (e Engine) GenerateThumbnail(original image.Image, size int, etag string) error {
	// if the cached version doesnt exist, generate it
	// copy over the cached version
	thumbnailCachePath := e.Config.ThumbnailCachePathOrDefault()
	thumbnailPath := filepath.Join(thumbnailCachePath, etag, fmt.Sprintf("%d.jpg", size))
	if !fileutil.Exists(thumbnailPath) {
		logger.MaybeSyncInfof(e.Log, "generating cached thumbnail `%s` @ %dpx", etag, size)
		if err := MakeDir(filepath.Join(thumbnailCachePath, etag)); err != nil {
			return err
		}
		if err := e.Resize(original, thumbnailPath, uint(size)); err != nil {
			return err
		}
	}

	return nil
}

// ShouldGenerateThumbnails returns if we should process any thumbnais for a given etag.
func (e Engine) ShouldGenerateThumbnails(etag string) bool {
	thumbnailCachePath := e.Config.ThumbnailCachePathOrDefault()
	for _, size := range e.Config.ImageSizesOrDefault() {
		thumbnailPath := filepath.Join(thumbnailCachePath, etag, fmt.Sprintf(constants.ImageSizeFormat, size))
		if !fileutil.Exists(thumbnailPath) {
			return true
		}
	}
	return false
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

// CopyThumbnails copies all thumbnails to the destination path by etag from the thumbnail cache.
func (e Engine) CopyThumbnails(etag, destinationPath string) error {
	for _, size := range e.Config.ImageSizesOrDefault() {
		if err := e.CopyThumbnail(etag, destinationPath, size); err != nil {
			return err
		}
	}
	return nil
}

// CopyThumbnail copies a cached thumbnail to the output directory.
func (e Engine) CopyThumbnail(etag, destinationPath string, size int) error {
	thumbnailCachePath := e.Config.ThumbnailCachePathOrDefault()
	thumbnailPath := filepath.Join(thumbnailCachePath, etag, fmt.Sprintf("%d.jpg", size))
	outputPath := filepath.Join(destinationPath, fmt.Sprintf(constants.ImageSizeFormat, size))
	logger.MaybeSyncInfof(e.Log, "copying cached thumbnail `%s` @ %dpx", destinationPath, size)
	if err := Copy(thumbnailPath, outputPath); err != nil {
		return err
	}
	return nil
}

// WriteTemplate writes a template to a given path with a given data viewmodel.
func (e Engine) WriteTemplate(tpl *template.Template, outputPath string, data *model.ViewModel) error {
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

// ParseSlugTemplate ensures the slug template
func (e *Engine) ParseSlugTemplate() (*template.Template, error) {
	return ParseTemplate(e.Config.SlugTemplateOrDefault())
}

// CreateSlug creates a slug for a post.
func (e Engine) CreateSlug(slugTemplate *template.Template, p model.Post) string {
	output, _ := RenderString(slugTemplate, p)
	return output
}
