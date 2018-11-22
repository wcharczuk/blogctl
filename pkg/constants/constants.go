package constants

const (
	// OutputPath is the default output path.
	OutputPath = "dist"
	// ImagesPath is the default source path for images.
	ImagesPath = "images"
	// LayoutPath is the default layouts path.
	LayoutPath = "layout"
	// PartialsPath is the default partials path.
	PartialsPath = "layout/partials"
	// ConfigPath is the default config file name.
	ConfigPath = "config.yml"
)

// Files are known filenames for layout templates.
const (
	TemplatePost = "post.html"
)

// OutputFiles are known output file names.
const (
	OutputFileIndex = "index.html"
)

// ImageOriginal is the original image.
const (
	ImageOriginal = "original.jpg"
	Image2048     = "2048.jpg"
	Image1024     = "1024.jpg"
	Image512      = "512.jpg"
)

// Partials are important / reused paritals or controls.
const (
	PartialImage = "image.html"
	PartialList  = "list.html"
	PartialAbout = "about.html"
)

// Discovery Paths
const (
	DiscoveryPathImages   = "images"
	DiscoveryPathStatic   = "static"
	DiscoveryPathPartials = "partials"
	DiscoveryFileMeta     = "meta.yml"
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
