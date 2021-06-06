package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/fileutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/constants"
	"github.com/wcharczuk/blogctl/pkg/model"
	"github.com/wcharczuk/blogctl/pkg/resize"
)

// New returns a new engine..
func New(options ...Option) (*Engine, error) {
	var e Engine
	for _, opt := range options {
		if err := opt(&e); err != nil {
			return nil, err
		}
	}
	return &e, nil
}

// MustNew returns a new engine but panics on error.
func MustNew(options ...Option) *Engine {
	e, err := New(options...)
	if err != nil {
		panic(err)
	}
	return e
}

// Option is a mutator for an engine.
type Option func(*Engine) error

// OptConfig sets the engine config.
func OptConfig(cfg config.Config) Option {
	return func(e *Engine) error {
		e.Config = cfg
		return nil
	}
}

// OptLog sets the logger.
func OptLog(log logger.Log) Option {
	return func(e *Engine) error {
		e.Log = log
		return nil
	}
}

// OptParallelism sets the engine parallelism if relevant.
func OptParallelism(paralellism int) Option {
	return func(e *Engine) error {
		e.Parallelism = paralellism
		return nil
	}
}

// OptDryRun sets DryRun on the engine.
func OptDryRun(dryRun bool) Option {
	return func(e *Engine) error {
		e.DryRun = dryRun
		return nil
	}
}

// Engine returns a
type Engine struct {
	Config      config.Config
	Parallelism int
	DryRun      bool
	Log         logger.Log
}

// ParallelismOrDefault is the parallelism or a default.
func (e Engine) ParallelismOrDefault() int {
	if e.Parallelism > 0 {
		return e.Parallelism
	}
	return runtime.NumCPU()
}

// Build generates the blog to the given output directory.
func (e Engine) Build(ctx context.Context) error {
	renderContext, err := e.BuildRenderContext(ctx)
	if err != nil {
		return err
	}

	if err := e.InitializeOutputPath(); err != nil {
		return err
	}

	if err := e.InitializeThumbnailCache(); err != nil {
		return err
	}

	ctx = WithRenderContext(ctx, renderContext)
	if err := e.Render(ctx); err != nil {
		return err
	}

	columns, rows := renderContext.Stats.TableData()
	for index, column := range columns {
		logger.MaybeInfof(e.Log, "%s: %s", column, rows[0][index])
	}

	return nil
}

// InitializeOutputPath creates the output path if it doesn't exist.
func (e Engine) InitializeOutputPath() error {
	if Exists(e.Config.OutputPathOrDefault()) {
		if err := ex.New(os.RemoveAll(e.Config.OutputPathOrDefault())); err != nil {
			return err
		}
	}
	return MakeDir(e.Config.OutputPathOrDefault())
}

// InitializeThumbnailCache creates the output path if it doesn't exist.
func (e Engine) InitializeThumbnailCache() error {
	if Exists(e.Config.ThumbnailCachePathOrDefault()) {
		return nil
	}
	return MakeDir(e.Config.ThumbnailCachePathOrDefault())
}

// DiscoverPosts generates the blog data.
func (e Engine) DiscoverPosts(ctx context.Context) (*model.Data, error) {
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
	postsPath := e.Config.PostsPathOrDefault()

	var postIndex int
	logger.MaybeInfof(e.Log, "searching `%s` for posts", postsPath)
	err = filepath.Walk(postsPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if currentPath == postsPath {
			return nil
		}
		if info.IsDir() {
			defer func() {
				postIndex++
			}()

			logger.MaybeDebugf(e.Log, "%s: reading post (%d)", currentPath, postIndex)

			// check if we have an image
			post, err := e.GeneratePost(ctx, slugTemplate, currentPath, postIndex)
			if err != nil {
				return err
			}
			output.Posts = append([]*model.Post{post}, output.Posts...)

			if !e.Config.SkipGenerateTags {
				for _, tag := range post.Meta.Tags {
					if tagPosts, ok := tags[tag]; ok {
						tagPosts.Posts = append(tagPosts.Posts, post)
					} else {
						tags[tag] = &model.Tag{
							Tag:   tag,
							Posts: []*model.Post{post},
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
	sort.Sort(model.Posts(output.Posts).Sort(e.Config.PostSortKeyOrDefault(), e.Config.PostSortAscendingOrDefault()))

	// create previous and next links for each post.
	for index := range output.Posts {
		if index > 0 {
			output.Posts[index].Previous = output.Posts[index-1]
		}
		if index < len(output.Posts)-1 {
			output.Posts[index].Next = output.Posts[index+1]
		}
	}

	if !e.Config.SkipGenerateTags {
		// add tags, make sure they're sorted.
		for _, tag := range tags {
			sort.Sort(model.Posts(tag.Posts).Sort(e.Config.PostSortKeyOrDefault(), e.Config.PostSortAscendingOrDefault()))
			output.Tags = append(output.Tags, *tag)
		}
		sort.Sort(model.Tags(output.Tags))
	}

	return &output, nil
}

// BuildRenderContext builds the render context used by the render function.
func (e Engine) BuildRenderContext(ctx context.Context) (*model.RenderContext, error) {
	partials, err := e.DiscoverPartials(ctx)
	if err != nil {
		return nil, err
	}
	data, err := e.DiscoverPosts(ctx)
	if err != nil {
		return nil, err
	}
	return &model.RenderContext{
		Data:     data,
		Partials: partials,
		Stats: model.Stats{
			NumPosts:      data.NumPosts(),
			NumTags:       data.NumTags(),
			NumImagePosts: data.NumImagePosts(),
			NumTextPosts:  data.NumTextPosts(),
			Earliest:      data.EarliestPost(),
			Latest:        data.LatestPost(),
		},
	}, nil
}

// Render writes the templates out for each of the posts.
func (e Engine) Render(ctx context.Context) error {
	renderContext := GetRenderContext(ctx)

	logger.MaybeInfof(e.Log, "rendering site with parallelism %d", e.ParallelismOrDefault())
	var err error

	outputPath := e.Config.OutputPathOrDefault()

	var defaultImagePostTemplate *template.Template
	imagePostTemplatePath := e.Config.ImagePostTemplateOrDefault()
	if imagePostTemplatePath != "" {
		_, defaultImagePostTemplate, err = e.CompileTemplate(imagePostTemplatePath, renderContext.Partials)
		if err != nil {
			return err
		}
	}

	var defaultTextPostTemplate *template.Template
	textPostTemplatePath := e.Config.TextPostTemplateOrDefault()
	if textPostTemplatePath != "" {
		_, defaultTextPostTemplate, err = e.CompileTemplate(textPostTemplatePath, renderContext.Partials)
		if err != nil {
			return err
		}
	}

	posts := make(chan interface{}, len(renderContext.Data.Posts))
	batchErrors := make(chan error, len(renderContext.Data.Posts))
	for _, post := range renderContext.Data.Posts {
		posts <- post
	}
	async.NewBatch(posts, func(ctx context.Context, workItem interface{}) error {
		post := workItem.(*model.Post)

		var postTemplate *template.Template
		if post.Text.SourcePath != "" {
			if post.Text.Template, post.Template, err = e.CompileTemplate(post.Text.SourcePath, renderContext.Partials); err != nil {
				return ex.New(err)
			}
		}

		if post.IsText() {
			postTemplate = defaultTextPostTemplate
		} else {
			postTemplate = defaultImagePostTemplate
		}

		slugPath := filepath.Join(outputPath, post.Slug)
		if err := MakeDir(slugPath); err != nil {
			return ex.New(err)
		}

		outputIndexPath := filepath.Join(slugPath, constants.FileIndex)
		logger.MaybeDebugf(e.Log, "%s: processing page", outputIndexPath)
		var postTextOutput string
		if postTextOutput, err = e.RenderTemplateToFile(postTemplate, outputIndexPath, &model.ViewModel{
			Config: e.Config,
			Posts:  renderContext.Data.Posts,
			Tags:   renderContext.Data.Tags,
			Post:   *post,
		}); err != nil {
			return err
		}
		if post.IsText() {
			post.Text.Output = postTextOutput
		}

		if post.Image.SourcePath != "" {
			if !e.Config.SkipCopyOriginalImage {
				if err := e.CopyImageOriginal(ctx, post.Image.SourcePath, slugPath); err != nil {
					return err
				}
			}
			if err := e.ProcessThumbnails(ctx, post.Image.SourcePath, slugPath); err != nil {
				return err
			}
		}
		return nil
	}, async.OptBatchParallelism(e.ParallelismOrDefault()), async.OptBatchErrors(batchErrors)).Process(ctx)

	if len(batchErrors) > 0 {
		return <-batchErrors
	}

	pagesPath := e.Config.PagesPathOrDefault()
	pages, err := ListDirectory(pagesPath)
	if err != nil {
		return err
	}
	for _, page := range pages {
		pageSourcePath := filepath.Join(pagesPath, page.Name())
		pageOutputPath := filepath.Join(outputPath, page.Name())

		logger.MaybeDebugf(e.Log, "%s: rendering page", pageOutputPath)
		_, pageTemplate, err := e.CompileTemplate(pageSourcePath, renderContext.Partials)
		if err != nil {
			return err
		}
		if _, err := e.RenderTemplateToFile(pageTemplate, pageOutputPath, &model.ViewModel{
			Config: e.Config,
			Post:   model.Posts(renderContext.Data.Posts).First(),
			Posts:  renderContext.Data.Posts,
			Tags:   renderContext.Data.Tags,
		}); err != nil {
			return err
		}
	}

	if !e.Config.SkipGenerateTags {
		tagTemplatePath := e.Config.TagTemplateOrDefault()
		if len(tagTemplatePath) > 0 && Exists(tagTemplatePath) {
			_, tagTemplate, err := e.CompileTemplate(tagTemplatePath, renderContext.Partials)
			if err != nil {
				return err
			}
			for _, tag := range renderContext.Data.Tags {
				tagPath := filepath.Join(outputPath, "tags", stringutil.Slugify(tag.Tag))
				if err := MakeDir(tagPath); err != nil {
					return ex.New(err)
				}
				if _, err := e.RenderTemplateToFile(tagTemplate, filepath.Join(tagPath, constants.FileIndex), &model.ViewModel{
					Config: e.Config,
					Posts:  renderContext.Data.Posts,
					Tags:   renderContext.Data.Tags,
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

	if !e.Config.SkipGenerateJSONData {
		dataOutputPath := filepath.Join(outputPath, constants.FileData)
		logger.MaybeDebugf(e.Log, "%s: rendering page", dataOutputPath)
		if err := e.WriteDataJSON(renderContext.Data, dataOutputPath); err != nil {
			return err
		}
	}

	return nil
}

// CleanThumbnailCache cleans the thumbnail cache by purging cached thumbnails for posts that may have been deleted.
func (e Engine) CleanThumbnailCache(ctx context.Context) error {
	postsPath := e.Config.PostsPathOrDefault()
	logger.MaybeInfof(e.Log, "%s: searching for posts", postsPath)
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
						return ex.New(err)
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
		return ex.New(err)
	}

	thumbnailCachePath := e.Config.ThumbnailCachePathOrDefault()
	logger.MaybeInfof(e.Log, "%s: comparing as thumbnail cache", thumbnailCachePath)

	// for each thumbnail cache folder
	var orphanedCachedPosts []string
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
		return ex.New(err)
	}

	// purge folders
	if len(orphanedCachedPosts) > 0 {
		for _, path := range orphanedCachedPosts {
			if !e.DryRun {
				if err := os.RemoveAll(filepath.Join(thumbnailCachePath, path)); err != nil {
					return ex.New(err)
				}
				logger.MaybeInfof(e.Log, "%s: purging orphaned cached directory", path)
			} else {
				logger.MaybeInfof(e.Log, "%s: (dry-run) would purge orphaned cached directory", path)
			}
		}
	} else {
		logger.MaybeInfof(e.Log, "no orphaned cached posts found")
	}

	return nil
}

//
// utilities
//

// DiscoverPartials reads all the partials named in the config.
// These are then injected into any subsequent renders as potential helper views.
func (e Engine) DiscoverPartials(ctx context.Context) ([]string, error) {
	partialsPath := e.Config.PartialsPathOrDefault()

	partialFiles, err := ListDirectory(partialsPath)
	if err != nil {
		return nil, err
	}

	var partials []string
	for _, partial := range partialFiles {
		contents, err := ioutil.ReadFile(filepath.Join(partialsPath, partial.Name()))
		if err != nil {
			return nil, ex.New(err).WithMessagef("partial: %s", partial)
		}
		partials = append(partials, string(contents))
	}
	return partials, nil
}

// GeneratePost reads post contents and metadata from a folder.
func (e Engine) GeneratePost(ctx context.Context, slugTemplate *template.Template, path string, postIndex int) (*model.Post, error) {
	files, err := ListDirectory(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ex.New("no child files found in post directory").WithMessage(path)
	}

	post := model.Post{
		Index: postIndex,
	}

	var postModTime time.Time
	for _, fi := range files {
		name := fi.Name()
		if strings.ToLower(name) == constants.FileMeta {
			if err := ReadYAML(filepath.Join(path, name), &post.Meta); err != nil {
				return nil, err
			}
		} else if HasExtension(name, constants.ImageExtensions...) {
			if post.Image.SourcePath != "" {
				return nil, ex.New("multiple image files found in post directory", ex.OptMessage(path))
			}
			post.Image.SourcePath = filepath.Join(path, name)
			if postModTime.Before(fi.ModTime()) {
				postModTime = fi.ModTime()
			}
			if post.Image, err = ReadImage(post.Image.SourcePath); err != nil {
				return nil, err
			}
		} else if HasExtension(name, constants.TemplateExtensions...) {
			if post.Text.SourcePath != "" {
				return nil, ex.New("multiple text files found in post directory", ex.OptMessage(path))
			}
			post.Text.SourcePath = filepath.Join(path, name)
			if postModTime.Before(fi.ModTime()) {
				postModTime = fi.ModTime()
			}
		} else {
			logger.MaybeDebugf(e.Log, "%s: ignoring erroneous post file", filepath.Join(path, name))
		}
	}

	post.Slug = e.CreateSlug(slugTemplate, post)
	post.ModTime = postModTime
	if post.Meta.Posted.IsZero() {
		post.Meta.Posted = postModTime
	}
	if post.IsImage() {
		post.Image.Sizes = e.GetImageSizePaths(post)
	}
	if post.Image.SourcePath == "" && post.Text.SourcePath == "" {
		return nil, ex.New("no image or text post data found", ex.OptMessage(path))
	}
	return &post, nil
}

// ProcessThumbnails processes thumbnails.
func (e Engine) ProcessThumbnails(ctx context.Context, originalFilePath, destinationPath string) error {
	originalContents, err := ioutil.ReadFile(originalFilePath)
	if err != nil {
		return ex.New(err)
	}

	etag, err := fileutil.ETag(originalContents)
	if err != nil {
		return err
	}

	if e.ShouldGenerateThumbnails(etag) {
		logger.MaybeInfof(e.Log, "%s: generating thumbnails", originalFilePath)
		if err := e.GenerateThumbnails(originalContents, originalFilePath, etag); err != nil {
			return nil
		}
	}

	if err := e.CopyThumbnails(etag, destinationPath); err != nil {
		return err
	}

	return nil
}

// CompileTemplate compiles a template.
func (e Engine) CompileTemplate(templatePath string, partials []string) (contents string, final *template.Template, err error) {
	fileContents, fileErr := ioutil.ReadFile(templatePath)
	if fileErr != nil {
		err = ex.New(fileErr).WithMessagef("template path: %s", templatePath)
		return
	}
	contents = string(fileContents)
	tmp := template.New(templatePath).Funcs(ViewFuncs())
	for _, partial := range partials {
		_, err = tmp.Parse(partial)
		if err != nil {
			err = ex.New(err).WithMessagef("template path: %s", templatePath)
			return
		}
	}

	final, err = tmp.Parse(string(contents))
	if err != nil {
		err = ex.New(err).WithMessagef("template path: %s", templatePath)
		return
	}
	return
}

// WriteDataJSON writes a data file to disk.
func (e Engine) WriteDataJSON(data *model.Data, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return ex.New(err)
	}
	defer f.Close()
	return ex.New(json.NewEncoder(f).Encode(data))
}

// GenerateThumbnails generates and copies our main thumbnails for the post image.
// - originalContents should be the bytes of the original image file
// - etag should be the sha sum as an etag, it is used as a path component in the file cache
func (e Engine) GenerateThumbnails(originalContents []byte, originalImagePath, etag string) error {
	// decode jpeg into image.Image
	original, err := jpeg.Decode(bytes.NewBuffer(originalContents))
	if err != nil {
		ex.New(err)
	}

	for _, size := range e.Config.ImageSizesOrDefault() {
		if err := e.GenerateThumbnail(original, size, originalImagePath, etag); err != nil {
			return err
		}
	}
	return nil
}

// GenerateThumbnail generates a thumbnail and stores it in the cache if it doesn't exist
// and copies the cached thumbail to the output directory.
func (e Engine) GenerateThumbnail(original image.Image, size int, imagePath, etag string) error {
	// if the cached version doesnt exist, generate it
	// copy over the cached version
	thumbnailCachePath := e.Config.ThumbnailCachePathOrDefault()
	thumbnailPath := filepath.Join(thumbnailCachePath, etag, fmt.Sprintf("%d.jpg", size))
	if !Exists(thumbnailPath) {
		logger.MaybeDebugf(e.Log, "%s: generating cached thumbnail @ %dpx", imagePath, size)
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
		if !Exists(thumbnailPath) {
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
		return ex.New(err)
	}
	defer out.Close()
	// write new image to file
	if err := jpeg.Encode(out, resized, nil); err != nil {
		return ex.New(err)
	}
	return nil
}

// CopyImageOriginal copies the original image to the destination.
func (e Engine) CopyImageOriginal(ctx context.Context, originalPath, destinationPath string) error {
	logger.MaybeDebugf(e.Log, "%s: copying original image", destinationPath)
	return Copy(originalPath, filepath.Join(destinationPath, constants.FileImageOriginal))
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
	logger.MaybeDebugf(e.Log, "%s: copying cached thumbnail @ %dpx", destinationPath, size)
	if err := Copy(thumbnailPath, outputPath); err != nil {
		return err
	}
	return nil
}

// RenderTemplateToFile writes a template to a given path with a given data viewmodel.
func (e Engine) RenderTemplateToFile(tpl *template.Template, outputPath string, data *model.ViewModel) (string, error) {
	buffer := new(bytes.Buffer)
	if err := tpl.Execute(buffer, data); err != nil {
		return "", ex.New(err)
	}
	if err := ioutil.WriteFile(outputPath, buffer.Bytes(), 0666); err != nil {
		return "", ex.New(err)
	}

	return buffer.String(), nil
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

// GetImageSizePaths gets the map that corresponds to the image sizes and the image path.
func (e Engine) GetImageSizePaths(post model.Post) map[string]string {
	output := make(map[string]string)
	if !e.Config.SkipCopyOriginalImage {
		output["original"] = filepath.Join(post.Slug, constants.FileImageOriginal)
	}
	for _, size := range e.Config.ImageSizesOrDefault() {
		output[strconv.Itoa(size)] = filepath.Join(post.Slug, fmt.Sprintf("%d.jpg", size))
	}
	return output
}
