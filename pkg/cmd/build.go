package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/sh"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Build returns the build command.
func Build(flags config.Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPath, err := config.ReadConfig(flags)
			sh.Fatal(err)

			log := Logger(flags, "build")
			slant.Print(log.Output, "BLOGCTL")

			log.Infof("using logger flags: %v", log.Flags.String())
			if cfgPath != "" {
				log.Infof("using config path: %s", cfgPath)
			}
			log.Infof("using parallelism: %d", *flags.Parallelism)

			if err := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptLog(log),
				engine.OptParallelism(*flags.Parallelism),
			).Build(context.Background()); err != nil {
				sh.Fatal(err)
			}
		},
	}
}
