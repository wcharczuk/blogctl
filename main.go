package main

import (
	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/sh"

	"github.com/wcharczuk/blogctl/pkg/cmd"
	"github.com/wcharczuk/blogctl/pkg/config"
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

	flags := config.NewFlags(blogctl)

	// add commands
	blogctl.AddCommand(cmd.Init(flags))
	blogctl.AddCommand(cmd.Build(flags))
	blogctl.AddCommand(cmd.Clean(flags))
	blogctl.AddCommand(cmd.Deploy(flags))
	blogctl.AddCommand(cmd.Fix(flags))
	blogctl.AddCommand(cmd.New(flags))
	blogctl.AddCommand(cmd.Server(flags))
	blogctl.AddCommand(cmd.Show(flags))

	sh.Fatal(blogctl.Execute())
}
