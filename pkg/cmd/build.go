package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Build returns the build command.
func Build(flags config.PersistentFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "build",
		Short:   "Build the photoblog",
		Aliases: []string{"b", "build", "g", "generate"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPath, err := config.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}
			log := Logger(cfg, "build")
			fmt.Fprintf(log.Logger.Output, banner)
			if cfgPath != "" {
				log.Infof("using config path: %s", cfgPath)
			}

			if err := engine.MustNew(engine.OptConfig(cfg), engine.OptLog(log)).Generate(context.Background()); err != nil {
				logger.FatalExit(err)
			}
			log.Info("complete")
		},
	}
}
