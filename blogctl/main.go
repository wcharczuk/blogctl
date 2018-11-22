package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/photoblog/blogctl/cmd"
)

/*
general usage

blogctl
	init : touch an empy instance of the photo blog
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
	configPath := blogctl.PersistentFlags().String("config", "config.yml", "The config file path")

	log := logger.All()

	blogctl.AddCommand(cmd.Init(configPath, log))
	blogctl.AddCommand(cmd.New(configPath, log))
	blogctl.AddCommand(cmd.Build(configPath, log))
	blogctl.AddCommand(cmd.Deploy(configPath, log))
	blogctl.AddCommand(cmd.Server(configPath, log))
	blogctl.Execute()
}
