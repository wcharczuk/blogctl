package config

import "github.com/wcharczuk/blogctl/pkg/constants"

// These are set by ldflags.
var (
	Version = ""
	GitRef  = ""
)

// Extra is just exta data you might want to pass into the renderer.
type Extra = map[string]interface{}

// Config is the blog config
type Config struct {
	// Title is the title for the blog.
	Title string `json:"title" yaml:"title"`
	// Author is your name.
	Author string `json:"author" yaml:"author"`
	// Description is a description for the blog, will be used in html head meta.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// BaseURL is the base url for the blog.
	BaseURL string `json:"baseURL,omitempty" yaml:"baseURL,omitempty"`
	// PostsPath is the path to the posts to compile.
	PostsPath string `json:"postsPath,omitempty" yaml:"postsPath,omitempty"`
	// PagesPath is the path to a folder with pages to compile.
	// They are rendered and copied to the root of the output.
	PagesPath string `json:"pagesPath,omitempty" yaml:"pagesPath,omitempty"`
	// OutputPath is the compiled site path.
	OutputPath string `json:"outputPath,omitempty" yaml:"outputPath,omitempty"`
	// PatialsPath is the path to a folder with partials to include
	// when rendering pages and the posts.
	PartialsPath string `json:"partialsPath,omitempty" yaml:"partialsPath,omitempty"`
	// StaticPath is the path to a folder with static files to copy over.
	StaticsPath string `json:"staticsPath,omitempty" yaml:"staticsPath,omitempty"`
	// ThumbnailCachePath is the path to the thumbnail cache.
	ThumbnailCachePath string `json:"thumbnailCachePath,omitempty" yaml:"thumbnailCachePath,omitempty"`
	// SlugTemplate is the template for post slugs.
	// It defaults to "/{{ .Meta.Posted.Year }}/{{ .Meta.Posted.Month }}/{{ .Meta.Posted.Day }}/{{ .Meta.Title | slugify }}/"
	SlugTemplate string `json:"slugTemplate,omitempty" yaml:"slugTemplate,omitempty"`
	// PostTemplate is the path to the post template file.
	// It is what is rendered when you go to /<POST_SLUG>/
	PostTemplatePath string `json:"postTemplatePath,omitempty" yaml:"postTemplatePath,omitempty"`
	// TagTemplate is the path to the tag template file.
	// It is what is rendered when you go to /tags/:tag_name
	TagTemplatePath string `json:"tagTemplatePath,omitempty" yaml:"tagTemplatePath,omitempty"`
	// ImageSizes lets you set what size thumbnails to create from post files.
	// This defaults to 2048px, 1024px, and 512px.
	ImageSizes []int `json:"imageSizes,omitempty" yaml:"imageSizes,omitempty"`
	// Extra is optional and allows you to provide variables for templates.
	Extra Extra `json:"extra,omitempty" yaml:"extra,omitempty"`
	// S3 governs how the blog is deployed.
	S3 S3 `json:"s3,omitempty" yaml:"s3,omitempty"`
	// Cloudfront governs options for how the s3 files are cached.
	Cloudfront Cloudfront `json:"cloudfront,omitempty" yaml:"cloudfront,omitempty"`

	// below are knobs you can turn to disable specific things.

	// SkipTags instructs the engine to not create tag summary pages.
	SkipTags bool `json:"skipTags,omitempty" yaml:"skipTags,omitempty"`
	// SkipJSONData instructs the engine not to create a data.json file.
	SkipJSONData bool `json:"skipJSONData,omitempty" yaml:"skipJSONData,omitempty"`
}

// Fields returns fields to prompt for when creating a new config.
func (c *Config) Fields() []Field {
	return []Field{
		{Prompt: "Title (the title of the blog)", FieldReference: &c.Title},
		{Prompt: "Author (your name)", FieldReference: &c.Author},
		{Prompt: "Base URL (fully qualified, i.e https://foo.com)", FieldReference: &c.BaseURL},
		{Prompt: "Output Path (where the compiled blog goes)", FieldReference: &c.OutputPath, Default: constants.DefaultOutputPath},
		{Prompt: "Posts Path (where your posts live)", FieldReference: &c.PostsPath, Default: constants.DefaultPostsPath},
		{Prompt: "Pages Path (pages to render)", FieldReference: &c.PagesPath, Default: constants.DefaultPagesPath},
		{Prompt: "Partials Path (partials to include)", FieldReference: &c.PartialsPath, Default: constants.DefaultPartialsPath},
		{Prompt: "Statics Path (files to copy to output)", FieldReference: &c.StaticsPath, Default: constants.DefaultStaticsPath},
		{Prompt: "Slug Template (template literal for slugs)", FieldReference: &c.SlugTemplate, Default: constants.DefaultSlugTemplate},
		{Prompt: "Post Template Path (template file to use for each post)", FieldReference: &c.PostTemplatePath, Default: constants.DefaultPostTemplatePath},
		{Prompt: "Tag Template Path (template file to use for each tag)", FieldReference: &c.TagTemplatePath, Default: constants.DefaultTagTemplatePath},
		{Prompt: "Thumbnail Cache Path (resized image cache path)", FieldReference: &c.ThumbnailCachePath, Default: constants.DefaultThumbnailCachePath},
	}
}

// TitleOrDefault returns the title or a default.
func (c Config) TitleOrDefault() string {
	return c.Title
}

// AuthorOrDefault returns the author or a default.
func (c Config) AuthorOrDefault() string {
	return c.Author
}

// BaseURLOrDefault returns the base url or a default.
func (c Config) BaseURLOrDefault() string {
	return c.BaseURL
}

// OutputPathOrDefault returns the output path.
func (c Config) OutputPathOrDefault() string {
	if c.OutputPath != "" {
		return c.OutputPath
	}
	return constants.DefaultOutputPath
}

// PostsPathOrDefault returns the images path.
func (c Config) PostsPathOrDefault() string {
	if c.PostsPath != "" {
		return c.PostsPath
	}
	return constants.DefaultPostsPath
}

// SlugTemplateOrDefault returns the slug template or default.
func (c Config) SlugTemplateOrDefault() string {
	if c.SlugTemplate != "" {
		return c.SlugTemplate
	}
	return constants.DefaultSlugTemplate
}

// PostTemplateOrDefault returns the single post template or a default.
func (c Config) PostTemplateOrDefault() string {
	if c.PostTemplatePath != "" {
		return c.PostTemplatePath
	}
	return constants.DefaultPostTemplatePath
}

// TagTemplateOrDefault returns the single tag template or a default.
func (c Config) TagTemplateOrDefault() string {
	if c.TagTemplatePath != "" {
		return c.TagTemplatePath
	}
	return constants.DefaultTagTemplatePath
}

// PagesPathOrDefault returns page file paths or defaults.
func (c Config) PagesPathOrDefault() string {
	if c.PagesPath != "" {
		return c.PagesPath
	}
	return constants.DefaultPagesPath
}

// PartialsPathOrDefault returns partial file paths or defaults.
func (c Config) PartialsPathOrDefault() string {
	if c.PartialsPath != "" {
		return c.PartialsPath
	}
	return constants.DefaultPartialsPath
}

// StaticsPathOrDefault returns static file paths or defaults.
func (c Config) StaticsPathOrDefault() string {
	if c.StaticsPath != "" {
		return c.StaticsPath
	}
	return constants.DefaultStaticsPath
}

// ThumbnailCachePathOrDefault returns static file paths or defaults.
func (c Config) ThumbnailCachePathOrDefault() string {
	if c.ThumbnailCachePath != "" {
		return c.ThumbnailCachePath
	}
	return constants.DefaultThumbnailCachePath
}

// ImageSizesOrDefault returns the image sizes or a default set.
func (c Config) ImageSizesOrDefault() []int {
	if c.ImageSizes != nil {
		return c.ImageSizes
	}
	return constants.DefaultImageSizes
}
