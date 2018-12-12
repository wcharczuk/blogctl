package config

import "github.com/wcharczuk/photoblog/pkg/constants"

// These are set by ldflags.
var (
	Version = ""
	GitRef  = ""
)

// Config is the blog config
type Config struct {
	// Title is the title for the blog.
	Title string `json:"title" yaml:"title"`
	// Author is your name.
	Author string `json:"author" yaml:"author"`
	// Description is a description for the blog, will be used in html head meta.
	Description string `json:"description" yaml:"description"`
	// BaseURL is the base url for the blog.
	BaseURL string `json:"baseURL" yaml:"baseURL"`
	// PostsPath is the path to the posts to compile.
	PostsPath string `json:"postsPath" yaml:"postsPath"`
	// PagesPath is the path to a folder with pages to compile.
	// They are rendered and copied to the root of the output.
	PagesPath string `json:"pagesPath" yaml:"pagesPath"`
	// OutputPath is the compiled site path.
	OutputPath string `json:"outputPath" yaml:"outputPath"`
	// PatialsPath is the path to a folder with partials to include
	// when rendering pages and the posts.
	PartialsPath string `json:"partialsPath" yaml:"partialsPath"`
	// StaticPath is the path to a folder with static files to copy over.
	StaticPath string `json:"staticPath" yaml:"staticPath"`
	// ThumbnailCachePath is the path to the thumbnail cache.
	ThumbnailCachePath string `json:"thumbnailCachePath" yaml:"thumbnailCachePath"`
	// SlugTemplate is the template for post slugs.
	// It defaults to "/{{ .Meta.Posted.Year }}/{{ .Meta.Posted.Month }}/{{ .Meta.Posted.Day }}/{{ .Meta.Title | slugify }}/"
	SlugTemplate string `json:"slugTemplate" yaml:"slugTemplate"`
	// PostTemplate is the path to the post template file.
	// It is what is rendered when you go to /<POST_SLUG>/
	PostTemplate string `json:"postTemplate" yaml:"postTemplate"`
	// TagTemplate is the path to the tag template file.
	// It is what is rendered when you go to /tags/:tag_name
	TagTemplate string `json:"tagTemplate" yaml:"tagTemplate"`
	// ImageSizes lets you set what size thumbnails to create from post files.
	// This defaults to 2048px, 1024px, and 512px.
	ImageSizes []int `json:"imageSizes" yaml:"imageSizes"`
	// Extra is optional and allows you to provide variables for templates.
	Extra Extra `json:"extra" yaml:"extra"`
	// S3 governs how the blog is deployed.
	S3 S3 `json:"s3" yaml:"s3"`
	// S3 governs how the blog is deployed.
	Cloudfront Cloudfront `json:"cloudfront" yaml:"cloudfront"`
}

// TitleOrDefault returns the title or a default.
func (c Config) TitleOrDefault() string {
	if c.Title != "" {
		return c.Title
	}
	return "Unset"
}

// AuthorOrDefault returns the author or a default.
func (c Config) AuthorOrDefault() string {
	return c.Author
}

// BaseURLOrDefault returns the base url or a default.
func (c Config) BaseURLOrDefault() string {
	return c.BaseURL
}

// PostsPathOrDefault returns the images path.
func (c Config) PostsPathOrDefault() string {
	if c.PostsPath != "" {
		return c.PostsPath
	}
	return constants.DefaultPostsPath
}

// OutputPathOrDefault returns the output path.
func (c Config) OutputPathOrDefault() string {
	if c.OutputPath != "" {
		return c.OutputPath
	}
	return constants.DefaultOutputPath
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
	if c.PostTemplate != "" {
		return c.PostTemplate
	}
	return constants.DefaultPostTemplate
}

// TagTemplateOrDefault returns the single tag template or a default.
func (c Config) TagTemplateOrDefault() string {
	if c.TagTemplate != "" {
		return c.TagTemplate
	}
	return constants.DefaultTagTemplate
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

// StaticPathOrDefault returns static file paths or defaults.
func (c Config) StaticPathOrDefault() string {
	if c.StaticPath != "" {
		return c.StaticPath
	}
	return constants.DefaultStaticPath
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

// Extra is just exta data you might want to pass into the renderer.
type Extra map[string]interface{}
