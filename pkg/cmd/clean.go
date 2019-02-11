package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Clean returns the clean command.
func Clean(configPath *string, log *logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "clean",
		Short:   "Clean the thumbnail cache",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := engine.ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.New(config).WithLogger(log.SubContext("clean")).CleanThumbnailCache(); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
}
