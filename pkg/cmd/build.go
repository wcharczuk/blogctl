package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Build returns the build command.
func Build(flags *config.PersistentFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "build",
		Short:   "Build the photoblog",
		Aliases: []string{"b", "build", "g", "generate"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPath, err := engine.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}
			log := logger.MustNew(logger.OptConfig(cfg.Logger)).SubContext("blogctl").SubContext("build")

			if cfgPath != "" {
				log.Infof("using config path: %s", cfgPath)
			}

			if err := engine.New(cfg).WithLogger(log).Generate(); err != nil {
				logger.FatalExit(err)
			}
			log.Info("complete")
		},
	}
}
