package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Build returns the build command.
func Build(configPath *string, log *logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "build",
		Short:   "Build the photoblog",
		Aliases: []string{"b", "build", "g", "generate"},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := engine.ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.New(config).WithLogger(log.SubContext("build")).Generate(); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
}
