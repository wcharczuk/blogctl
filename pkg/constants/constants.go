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
	TemplateIndex  = "index.html"
	TemplateHome   = "home.html"
	TemplateSingle = "single.html"
)

// FileIndex are known output file names.
const (
	FileIndex = "index.html"
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
	DiscoveryFileMetadata = "meta.yml"
)

// Extensions are file suffixes that indicate file type.
const (
	ExtensionJPG  = ".jpg"
	ExtensionJPEG = ".jpeg"
	ExtensionPNG  = ".png"
)

// ImageExtensions are known image extensions.
var (
	ImageExtensions = []string{
		ExtensionJPG,
		ExtensionJPEG,
		ExtensionPNG,
	}
)
