package constants

const (
	// DefaultConfigPath is the default config file name.
	DefaultConfigPath = "./config.yml"
	// DefaultOutputPath is the default output path.
	DefaultOutputPath = "./dist"
	// DefaultPostsPath is the default source path for post images.
	DefaultPostsPath = "./posts"
	// DefaultStaticsPath is the default source path static files.
	DefaultStaticsPath = "./static"
	// DefaultLayoutPath is the default layout path.
	DefaultLayoutPath = "./layout"
	// DefaultThumbnailCachePath is the default thumbnail cache path.
	DefaultThumbnailCachePath = "./thumbnails"
	// DefaultPagesPath is the default partials path.
	DefaultPagesPath = "./layout/pages"
	// DefaultPartialsPath is the default partials path.
	DefaultPartialsPath = "./layout/partials"
	// DefaultPostTemplatePath is the default post template path.
	DefaultPostTemplatePath = "./layout/post.html"
	// DefaultTagTemplatePath is the default tag template path.
	DefaultTagTemplatePath = "./layout/tag.html"
)

// DefaultSlugTemplate is the default slug format.
const (
	DefaultSlugTemplate = `{{ .Meta.Posted | time_format "2006/01/02" }}/{{ .Meta.Title | slugify }}`
)

// OutputFiles are known output file names.
const (
	FileIndex         = "index.html"
	FileMeta          = "meta.yml"
	FileData          = "data.json"
	FileImageOriginal = "original.jpg"
)

// Sizes are the default sizes for the resized images.
const (
	SizeLarge  = 2048
	SizeMedium = 1024
	SizeSmall  = 512
)

// DefaultImageSizes are the default resize dimensions.
// The size given will be the largest dimension of the final image.
// If the image is landscape, it will be the width.
// If the image is portrait, it will be the height.
var (
	DefaultImageSizes = []int{
		SizeLarge,
		SizeMedium,
		SizeSmall,
	}
)

// ImageSizeFormat is the format for the thumbnail images..
const (
	ImageSizeFormat = "%d.jpg"
)

// Extensions are file suffixes that indicate file type.
const (
	ExtensionJPG  = ".jpg"
	ExtensionJPEG = ".jpeg"
	ExtensionHTML = ".html"
)

// Known Extensions
var (
	ImageExtensions = []string{
		ExtensionJPG,
		ExtensionJPEG,
	}
	TemplateExtensions = []string{
		ExtensionHTML,
	}
)
