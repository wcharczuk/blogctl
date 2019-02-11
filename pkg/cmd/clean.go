package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Clean returns the clean command.
func Clean(configPath *string, log *logger.Logger) *cobra.Command {
	var dryRun *bool
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "Clean the thumbnail cache",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := engine.ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.New(config).WithLogger(log.SubContext("clean")).CleanThumbnailCache(*dryRun); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}

	dryRun = cmd.Flags().Bool("dry-run", false, "If we should only print which paths we would delete")
	return cmd
}
