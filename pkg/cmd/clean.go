package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Clean returns the clean command.
func Clean(flags *config.PersistentFlags) *cobra.Command {
	var dryRun *bool
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "Clean the thumbnail cache",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, configPath, err := engine.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}
			log := logger.MustNew(logger.OptConfig(cfg.Logger)).SubContext("blogctl").SubContext("clean")
			if configPath != "" {
				log.Infof("using config path: %s", configPath)
			}

			if err := engine.New(cfg).WithLogger(log).CleanThumbnailCache(*dryRun); err != nil {
				logger.FatalExit(err)
			}
			log.Info("complete")
		},
	}

	dryRun = cmd.Flags().Bool("dry-run", false, "If we should only print which paths we would delete")
	return cmd
}
