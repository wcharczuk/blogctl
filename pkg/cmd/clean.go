package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Clean returns the clean command.
func Clean(flags config.Flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "Clean the thumbnail cache",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, configPath, err := config.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}

			log := Logger(flags, "clean")
			slant.Print(log.Output, "BLOGCTL")

			if configPath != "" {
				log.Infof("using config path: %s", configPath)
			}

			if err := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptLog(log),
				engine.OptParallelism(*flags.Parallelism),
				engine.OptDryRun(*flags.DryRun),
			).CleanThumbnailCache(context.Background()); err != nil {
				logger.FatalExit(err)
			}
			log.Info("complete")
		},
	}
	return cmd
}
