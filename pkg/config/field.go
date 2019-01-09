package config

// Field is a field to prompt for on the config.
type Field struct {
	Prompt         string
	FieldReference *string
	Default        string
}
