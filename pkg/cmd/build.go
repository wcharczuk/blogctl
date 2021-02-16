package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi/slant"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Build returns the build command.
func Build(flags config.Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPaths, err := config.ReadConfig(flags)
			Fatal(err)

			log := Logger(flags, "build")
			slant.Print(log.Output, "BLOGCTL")

			log.Infof("using logger flags: %v", log.Flags.String())
			if len(cfgPaths) > 0 {
				log.Infof("using config path(s): %s", strings.Join(cfgPaths, ", "))
			}
			log.Infof("using parallelism: %d", *flags.Parallelism)

			if err := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptLog(log),
				engine.OptParallelism(*flags.Parallelism),
			).Build(context.Background()); err != nil {
				Fatal(err)
			}
		},
	}
}
