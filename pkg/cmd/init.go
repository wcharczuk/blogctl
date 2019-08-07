package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/sh"

	"github.com/spf13/cobra"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/constants"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Init returns the init command.
func Init(flags config.Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "init [NAME]",
		Short: "Initialize a new blog",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if name == "" {
				sh.Fatalf("must provide a folder name")
			}

			slant.Print(os.Stdout, "blogctl")
			fmt.Println("Initializing a new blog")
			fmt.Println("Please provide the fields for the `config.yml`")
			fmt.Println("They will be prompted as `Field (explanation) [default value]`")

			// create a new config
			config := new(config.Config)
			fields := config.Fields()
			var value string
			for _, field := range fields {
				if field.Default != "" {
					value = sh.Promptf("%s [default %s]: ", field.Prompt, field.Default)
				} else {
					value = sh.Promptf("%s: ", field.Prompt)
				}
				if value != "" {
					*field.FieldReference = value
				} else {
					*field.FieldReference = field.Default
				}
			}

			if err := engine.MakeDir(name); err != nil {
				sh.Fatal(err)
			}
			if err := engine.MakeDir(filepath.Join(name, config.PostsPathOrDefault())); err != nil {
				sh.Fatal(err)
			}
			if err := engine.MakeDir(filepath.Join(name, config.PagesPathOrDefault())); err != nil {
				sh.Fatal(err)
			}
			if err := engine.MakeDir(filepath.Join(name, config.PartialsPathOrDefault())); err != nil {
				sh.Fatal(err)
			}
			if err := engine.MakeDir(filepath.Join(name, config.StaticsPathOrDefault())); err != nil {
				sh.Fatal(err)
			}
			if err := engine.MakeDir(filepath.Join(name, config.StaticsPathOrDefault(), "css")); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteYAML(filepath.Join(name, *flags.ConfigPath), config); err != nil {
				sh.Fatal(err)
			}

			/* write individual files */
			if err := engine.WriteFile(filepath.Join(name, config.PartialsPathOrDefault(), "header.html"), []byte(headerHTML)); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteFile(filepath.Join(name, config.PartialsPathOrDefault(), "footer.html"), []byte(footerHTML)); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteFile(filepath.Join(name, config.PagesPathOrDefault(), constants.FileIndex), []byte(indexHTML)); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteFile(filepath.Join(name, config.ImagePostTemplateOrDefault()), []byte(imageHTML)); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteFile(filepath.Join(name, config.TextPostTemplateOrDefault()), []byte(textHTML)); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteFile(filepath.Join(name, config.TagTemplateOrDefault()), []byte(tagHTML)); err != nil {
				sh.Fatal(err)
			}
			if err := engine.WriteFile(filepath.Join(name, config.StaticsPathOrDefault(), "css/site.css"), []byte(siteCSS)); err != nil {
				sh.Fatal(err)
			}
		},
	}
}

const (
	headerHTML = `{{ define "header" }}
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>{{ .TitleOrDefault }}</title>
	<meta name="author" content="{{ .Config.Author }}">
	<meta name="description" content="{{ .Config.Description }}">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="/css/site.css">
</head>
<body>
	<div class="content">
{{ end }}`

	footerHTML = `{{ define "footer" }}
	</div>
</body>
</html>
{{ end }}`

	indexHTML = `{{ template "header" . }}
{{ range $index, $post := .Posts }}
	<div class="post">
		<img src="{{ $post.ImageSourceSmall }}" />
	</div>
{{ else }}
	<h2>No Posts.</h2>
{{ end }}
{{ template "footer" . }}`

	imageHTML = `{{ template "header" . }}
<div class="image post">
	<img src="{{ .Post.ImageSourceLarge}}" />
</div>
{{ template "footer" . }}`

	textHTML = `{{ template "header" . }}
<div class="text post">
	{{ render .Post }}
</div>
{{ template "footer" . }}`

	tagHTML = `{{ template "header" . }}
<div class="tag">
	{{ range $index, $post := .Tag.Posts }}
	<div class="post">
		<img src="{{$post.ImageSourceSmall}}" />
	</div>
	{{ else }}
	<h2>No Posts For Tag.</h2>
	{{ end }}
</div>
{{ template "footer" . }}`

	siteCSS = `body { font-family: 'sans-serif'; margin: 0; padding: 0; }

.post { display: inline-block; }

.post img {
	width: auto;
	max-width: calc(100vw - 20px);
	height: auto;
	max-height: calc(100vh - 20px);
}`
)
