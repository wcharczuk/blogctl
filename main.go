package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wcharczuk/blogctl/pkg/cmd"
)

/*
general usage

blogctl
	init : touch an empy instance of the photo blog
	new : create a new post
	build : compile the posts into static pages
	deploy : push it to aws/gcp/*
	server : start a local server against the output folder

flags:
--config
*/

func main() {
	blogctl := &cobra.Command{
		Use: "blogctl",
	}
	configPath := blogctl.PersistentFlags().String("config", "./config.yml", "The config file path")

	// add commands
	blogctl.AddCommand(cmd.Init(configPath))
	blogctl.AddCommand(cmd.New(configPath))
	blogctl.AddCommand(cmd.Build(configPath))
	blogctl.AddCommand(cmd.Clean(configPath))
	blogctl.AddCommand(cmd.Deploy(configPath))
	blogctl.AddCommand(cmd.Server(configPath))

	if err := blogctl.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
