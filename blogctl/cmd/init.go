package cmd

import (
	"path/filepath"

	"github.com/wcharczuk/photoblog/pkg/constants"
	"github.com/wcharczuk/photoblog/pkg/engine"

	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
)

// Init returns the init command.
func Init(configPath *string, log *logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "init [NAME]",
		Short: "Initialize a new blog",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if err := engine.MakeDir(name); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.ImagesPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.LayoutPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.DiscoveryPathPages)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.DiscoveryPathPartials)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.DiscoveryPathStatic)); err != nil {
				log.SyncFatalExit(err)
			}
			//create the config
			config, err := RenderTemplate(configTemplate, Vars{"title": name})
			if err != nil {
				log.SyncFatalExit(err)
			}
			if err := WriteFile(filepath.Join(name, constants.ConfigPath), []byte(config)); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
}

const (
	configTemplate = `title: {{ .Var "name" }}
images: images
output: dist
layout:
  post: layout/post.html
  pages: layout/pages
  partials: layout/partials
  static: static
`
)
