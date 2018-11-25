package constants

const (
	// DefaultOutputPath is the default output path.
	DefaultOutputPath = "dist"
	// DefaultPostsPath is the default source path for post images.
	DefaultPostsPath = "posts"
	// DefaultStaticPath is the default source path static files.
	DefaultStaticPath = "static"
	// DefaultPagesPath is the default partials path.
	DefaultPagesPath = "layout/pages"
	// DefaultPartialsPath is the default partials path.
	DefaultPartialsPath = "layout/partials"
	// DefaultConfigPath is the default config file name.
	DefaultConfigPath = "config.yml"
	// DefaultPostTemplate is the default post template path.
	DefaultPostTemplate = "layout/post.html"
)

// OutputFiles are known output file names.
const (
	FileIndex = "index.html"
	FileMeta  = "meta.yml"
)

// Sizes are the default sizes for the resized images.
const (
	SizeLarge  = 2048
	SizeMedium = 1024
	SizeSmall  = 512
)

// DefaultImageSizes are the default resize dimensions.
var (
	DefaultImageSizes = []int{
		SizeLarge,
		SizeMedium,
		SizeSmall,
	}
)

// ImageOriginal is the original image.
const (
	ImageOriginal   = "original.jpg"
	ImageSizeFormat = "%d.jpg"
)

// Extensions are file suffixes that indicate file type.
const (
	ExtensionJPG  = ".jpg"
	ExtensionJPEG = ".jpeg"
)

// ImageExtensions are known image extensions.
var (
	ImageExtensions = []string{
		ExtensionJPG,
		ExtensionJPEG,
	}
)
