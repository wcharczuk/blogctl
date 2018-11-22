package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/photoblog/pkg/engine"
)

// Build returns the build command.
func Build(configPath *string, log *logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}

			engine := engine.Engine{
				Config: config,
				Log:    log,
			}
			if err := engine.Generate(); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
}
