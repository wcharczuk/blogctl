package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/yaml"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
	"github.com/wcharczuk/blogctl/pkg/model"
)

// Show returns the show tree of commands.
func Show(flags config.Flags) *cobra.Command {
	var outputFormat *string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details about posts, tags, or cache metdata",
	}
	outputFormat = cmd.PersistentFlags().StringP("output", "o", "name", "The output format; one of `name`, `table`, `json`, `yaml`")

	var postsOrderBy *string
	var postsOrderDesc *bool
	posts := &cobra.Command{
		Use:   "posts",
		Short: "Show posts",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _, err := config.ReadConfig(flags)
			sh.Fatal(err)
			e := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptParallelism(*flags.Parallelism),
				engine.OptDryRun(*flags.DryRun),
			)

			posts, err := e.DiscoverPosts(context.Background())
			sh.Fatal(err)

			switch strings.ToLower(*postsOrderBy) {
			case "location":
				if *postsOrderDesc {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].Meta.Location > posts.Posts[j].Meta.Location })
				} else {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].Meta.Location < posts.Posts[j].Meta.Location })
				}
			case "posted":
				if *postsOrderDesc {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].Meta.Posted.Before(posts.Posts[j].Meta.Posted) })
				} else {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].Meta.Posted.After(posts.Posts[j].Meta.Posted) })
				}
			case "slug":
				if *postsOrderDesc {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].Slug > posts.Posts[j].Slug })
				} else {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].Slug < posts.Posts[j].Slug })
				}
			case "title":
				if *postsOrderDesc {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].TitleOrDefault() > posts.Posts[j].TitleOrDefault() })
				} else {
					sort.Slice(posts.Posts, func(i, j int) bool { return posts.Posts[i].TitleOrDefault() < posts.Posts[j].TitleOrDefault() })
				}
			default:
				sh.Fatal(fmt.Errorf("invalid post order by: %s", *postsOrderBy))
			}

			switch strings.ToLower(*outputFormat) {
			case "name":
				for _, post := range posts.Posts {
					fmt.Fprintf(os.Stdout, "%s\n", post.TitleOrDefault())
				}
			case "json":
				sh.Fatal(json.NewEncoder(os.Stdout).Encode(posts.Posts))
			case "yaml":
				sh.Fatal(yaml.NewEncoder(os.Stdout).Encode(posts.Posts))
			case "table":
				sh.Fatal(ansi.TableForSlice(os.Stdout, model.Posts(posts.Posts).TableRows()))
			default:
				sh.Fatal(fmt.Errorf("invalid output format: %s", *outputFormat))
			}
		},
	}
	postsOrderBy = posts.Flags().String("order-by", "title", "Which field to order the posts by; one of `location`, `posted`, `slug`, or `title`")
	postsOrderDesc = posts.Flags().Bool("desc", false, "The posts sort order (true will sort descending)")

	var tagsOrderBy *string
	var tagsOrderDesc *bool
	tags := &cobra.Command{
		Use:   "tags",
		Short: "Show tags",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _, err := config.ReadConfig(flags)
			sh.Fatal(err)
			e := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptParallelism(*flags.Parallelism),
				engine.OptDryRun(*flags.DryRun),
			)

			posts, err := e.DiscoverPosts(context.Background())
			sh.Fatal(err)
			switch strings.ToLower(*tagsOrderBy) {
			case "tag":
				if *tagsOrderDesc {
					sort.Slice(posts.Tags, func(i, j int) bool { return posts.Tags[i].Tag > posts.Tags[j].Tag })
				} else {
					sort.Slice(posts.Tags, func(i, j int) bool { return posts.Tags[i].Tag < posts.Tags[j].Tag })
				}
			case "posts":
				if *tagsOrderDesc {
					sort.Slice(posts.Tags, func(i, j int) bool { return len(posts.Tags[i].Posts) < len(posts.Tags[j].Posts) })
				} else {
					sort.Slice(posts.Tags, func(i, j int) bool { return len(posts.Tags[i].Posts) > len(posts.Tags[j].Posts) })
				}
			}

			switch strings.ToLower(*outputFormat) {
			case "name":
				for _, tag := range posts.Tags {
					fmt.Fprintf(os.Stdout, "%s\n", tag.Tag)
				}
			case "json":
				sh.Fatal(json.NewEncoder(os.Stdout).Encode(posts.Tags))
			case "yaml":
				sh.Fatal(yaml.NewEncoder(os.Stdout).Encode(posts.Tags))
			case "table":
				sh.Fatal(ansi.TableForSlice(os.Stdout, model.Tags(posts.Tags).TableRows()))
			default:
				sh.Fatal(fmt.Errorf("invalid output format: %s", *outputFormat))
			}
		},
	}

	tagsOrderBy = tags.Flags().String("order-by", "tag", "Which field to order the tags by; one of `tag`, or `posts`")
	tagsOrderDesc = tags.Flags().Bool("desc", false, "The tags sort order (true will sort descending)")

	cmd.AddCommand(posts)
	cmd.AddCommand(tags)
	return cmd
}
