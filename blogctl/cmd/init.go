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
			if err := engine.MakeDir(filepath.Join(name, constants.DefaultPostsPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.DefaultPagesPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.DefaultPartialsPath)); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.MakeDir(filepath.Join(name, constants.DefaultStaticPath)); err != nil {
				log.SyncFatalExit(err)
			}
			//create the config
			config, err := RenderTemplate(configTemplate, Vars{
				"title":        name,
				"postsPath":    constants.DefaultPostsPath,
				"postTemplate": constants.DefaultPostTemplate,
				"pagesPath":    constants.DefaultPagesPath,
				"partialsPath": constants.DefaultPartialsPath,
				"staticPath":   constants.DefaultStaticPath,
			})
			if err != nil {
				log.SyncFatalExit(err)
			}
			if err := WriteFile(filepath.Join(name, constants.DefaultConfigPath), []byte(config)); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
}

const (
	configTemplate = `title: {{ .Var "title" }}
postsPath: {{ .Var "postsPath" }}
outputPath: {{ .Var "outputPath" }}
postTemplate: {{ .Var "postTemplate" }}
pagesPath: {{ .Var "pagesPath" }}
partialsPath: {{ .Var "partialsPath" }}
staticPath: {{ .Var "sttaicPath" }}
`
)
