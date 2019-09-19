package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Tags returns the tags command.
func Tags(flags config.Flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Display tags for posts",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, configPath, err := config.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}
			log := Logger(flags, "tags")
			slant.Print(log.Output, "BLOGCTL")
			if configPath != "" {
				log.Infof("using config path: %s", configPath)
			}
			e := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptLog(log),
				engine.OptParallelism(*flags.Parallelism),
				engine.OptDryRun(*flags.DryRun),
			)

			posts, err := e.DiscoverPosts(context.Background())
			if err != nil {
				logger.FatalExit(err)
			}

			for _, tag := range posts.Tags {
				fmt.Fprintf(log.Output, "%s (%d)\n", tag.Tag, len(tag.Posts))
			}
		},
	}
	return cmd
}
