package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/stringutil"

	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/constants"
	"github.com/wcharczuk/blogctl/pkg/engine"
	"github.com/wcharczuk/blogctl/pkg/model"
)

// New returns a new post command.
func New(flags config.Flags) *cobra.Command {
	var title, location, posted *string
	var tags *[]string
	cmd := &cobra.Command{
		Use:   "new [IMAGE_PATH]",
		Short: "Create a new blog post from a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imagePath := args[0]

			cfg, cfgPaths, err := config.ReadConfig(flags)
			Fatal(err)

			log := Logger(flags, "new")
			slant.Print(log.Output, "BLOGCTL")

			if len(cfgPaths) > 0 {
				log.Infof("using config path(s): %s", strings.Join(cfgPaths, ", "))
			}

			var postedDate time.Time
			if *posted != "" {
				postedDate, err = time.Parse("2006-01-02", *posted)
				Fatal(err)
			} else {
				postedDate, err = engine.ExtractCaptureDate(imagePath)
				Fatal(err)
			}

			if *title == "" {
				*title = filepath.Base(imagePath)
			}

			path := fmt.Sprintf("%s/%s-%s", cfg.PostsPathOrDefault(), postedDate.Format("2006-01-02"), stringutil.Slugify(*title))
			log.Infof("writing new post to %s", path)
			if _, err := os.Stat(path); err == nil {
				Fatal(fmt.Errorf("post directory already exists, aborting"))
			}
			fullPath := filepath.Join(path, filepath.Base(imagePath))
			Fatal(engine.Copy(imagePath, fullPath))

			var metaTags []string
			if tags != nil {
				metaTags = *tags
			}
			meta := model.Meta{
				Title:    *title,
				Location: *location,
				Posted:   postedDate,
				Tags:     metaTags,
			}
			Fatal(engine.WriteYAML(filepath.Join(path, constants.FileMeta), meta))
		},
	}

	title = cmd.Flags().String("title", "", "The title (optional, will default to the file name)")
	location = cmd.Flags().String("location", "", "The location (optional)")
	posted = cmd.Flags().String("posted", "", "The posted effective date (optional)")
	tags = cmd.Flags().StringArray("tag", nil, "Photo tags (optional)")
	return cmd
}
