package engine

// TitleOrDefaultProvider is a type that provides a title.
type TitleOrDefaultProvider interface {
	TitleOrDefault() string
}
