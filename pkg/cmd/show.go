package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/selector"
	"github.com/blend/go-sdk/sh"

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

	var postsOrderBy, postsSelector *string
	var postsOrderDesc *bool
	posts := &cobra.Command{
		Use:   "posts",
		Short: "Show posts",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _, err := config.ReadConfig(flags)
			Fatal(err)
			e := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptParallelism(*flags.Parallelism),
				engine.OptDryRun(*flags.DryRun),
			)

			posts, err := e.DiscoverPosts(context.Background())
			Fatal(err)

			if *postsSelector != "" {
				sel, err := selector.Parse(*postsSelector)
				Fatal(err)
				posts.Posts = model.Posts(posts.Posts).FilterBySelector(sel)
			}

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
	postsSelector = posts.Flags().StringP("labels", "l", "", "Filter posts with a given label selector (ex. `tree` for tagged with `tree)")
	postsOrderBy = posts.Flags().String("order-by", "title", "Which field to order the posts by; one of `location`, `posted`, `slug`, or `title`")
	postsOrderDesc = posts.Flags().Bool("desc", false, "The posts sort direction (`true` will sort descending, `false` ascending)")

	var tagsOrderBy *string
	var tagsOrderDesc *bool
	var tagsSimilar *bool
	tags := &cobra.Command{
		Use:   "tags",
		Short: "Show tags",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _, err := config.ReadConfig(flags)
			Fatal(err)
			e := engine.MustNew(
				engine.OptConfig(cfg),
				engine.OptParallelism(*flags.Parallelism),
				engine.OptDryRun(*flags.DryRun),
			)

			posts, err := e.DiscoverPosts(context.Background())
			Fatal(err)
			switch strings.ToLower(*tagsOrderBy) {
			case "tag":
				if *tagsOrderDesc {
					sort.Slice(posts.Tags, func(i, j int) bool { return posts.Tags[i].Tag > posts.Tags[j].Tag })
				} else {
					sort.Slice(posts.Tags, func(i, j int) bool { return posts.Tags[i].Tag < posts.Tags[j].Tag })
				}
			case "posts":
				if *tagsOrderDesc {
					sort.Slice(posts.Tags, func(i, j int) bool { return len(posts.Tags[i].Posts) > len(posts.Tags[j].Posts) })
				} else {
					sort.Slice(posts.Tags, func(i, j int) bool { return len(posts.Tags[i].Posts) < len(posts.Tags[j].Posts) })
				}
			}

			tags := posts.Tags

			if *tagsSimilar {
				tags = filterSimilar(tags, 1)
			}

			switch strings.ToLower(*outputFormat) {
			case "name":
				for _, tag := range tags {
					fmt.Fprintf(os.Stdout, "%s\n", tag.Tag)
				}
			case "json":
				sh.Fatal(json.NewEncoder(os.Stdout).Encode(tags))
			case "yaml":
				sh.Fatal(yaml.NewEncoder(os.Stdout).Encode(tags))
			case "table":
				sh.Fatal(ansi.TableForSlice(os.Stdout, model.Tags(tags).TableRows()))
			default:
				sh.Fatal(fmt.Errorf("invalid output format: %s", *outputFormat))
			}
		},
	}

	tagsSimilar = tags.Flags().Bool("similar", false, "Show only tags that have an small edit distance to each other")
	tagsOrderBy = tags.Flags().String("order-by", "tag", "Which field to order the tags by; one of `tag`, or `posts`")
	tagsOrderDesc = tags.Flags().Bool("desc", false, "The tags sort order (true will sort descending)")

	cmd.AddCommand(posts)
	cmd.AddCommand(tags)
	return cmd
}

func filterSimilar(tags []model.Tag, editDistance int) []model.Tag {
	var output []model.Tag

	for ai, a := range tags {
		var didAddA bool
		for bi, b := range tags {
			if ai == bi {
				continue
			}
			if distance := computeDistance(a.Tag, b.Tag); distance <= editDistance {
				if !didAddA {
					output = append(output, a)
					didAddA = true
				}
				output = append(output, b)
			}
		}
	}

	return output
}

// ComputeDistance computes the levenshtein distance between the two
// strings passed as an argument. The return value is the levenshtein distance
//
// Works on runes (Unicode code points) but does not normalize
// the input strings. See https://blog.golang.org/normalization
// and the golang.org/x/text/unicode/norm pacage.
func computeDistance(a, b string) int {
	if len(a) == 0 {
		return utf8.RuneCountInString(b)
	}

	if len(b) == 0 {
		return utf8.RuneCountInString(a)
	}

	if a == b {
		return 0
	}

	// We need to convert to []rune if the strings are non-ASCII.
	// This could be avoided by using utf8.RuneCountInString
	// and then doing some juggling with rune indices,
	// but leads to far more bounds checks. It is a reasonable trade-off.
	s1 := []rune(a)
	s2 := []rune(b)

	// swap to save some memory O(min(a,b)) instead of O(a)
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	lenS1 := len(s1)
	lenS2 := len(s2)

	// init the row
	x := make([]uint16, lenS1+1)
	// we start from 1 because index 0 is already 0.
	for i := 1; i < len(x); i++ {
		x[i] = uint16(i)
	}

	// make a dummy bounds check to prevent the 2 bounds check down below.
	// The one inside the loop is particularly costly.
	_ = x[lenS1]
	// fill in the rest
	for i := 1; i <= lenS2; i++ {
		prev := uint16(i)
		var current uint16
		for j := 1; j <= lenS1; j++ {
			if s2[i-1] == s1[j-1] {
				current = x[j-1] // match
			} else {
				current = min(min(x[j-1]+1, prev+1), x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[lenS1] = prev
	}
	return int(x[lenS1])
}

func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
