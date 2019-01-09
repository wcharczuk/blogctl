package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"

	"github.com/wcharczuk/blogctl/pkg/constants"
	"github.com/wcharczuk/blogctl/pkg/engine"
	"github.com/wcharczuk/blogctl/pkg/model"
	"github.com/wcharczuk/blogctl/pkg/stringutil"
)

// New returns a new post command.
func New(configPath *string, log *logger.Logger) *cobra.Command {
	var title, location, posted *string
	var tags *[]string
	cmd := &cobra.Command{
		Use:   "new [IMAGE_PATH]",
		Short: "Create a new blog post from a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imagePath := args[0]

			config, err := engine.ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}

			var postedDate time.Time
			if *posted != "" {
				postedDate, err = time.Parse("2006-01-02", *posted)
				if err != nil {
					log.SyncFatalExit(exception.New(err))
				}
			} else {
				postedDate, err = engine.ExtractCaptureDate(imagePath)
				if err != nil {
					log.SyncFatalExit(err)
				}
			}

			if *title == "" {
				*title = filepath.Base(imagePath)
			}

			path := fmt.Sprintf("%s/%s-%s", config.PostsPathOrDefault(), postedDate.Format("2006-01-02"), stringutil.Slugify(*title))
			log.SubContext("new").Infof("writing new post to %s", path)
			if err := engine.MakeDir(path); err != nil {
				log.SubContext("new").SyncFatal(err)
				os.Exit(1)
			}
			if err := engine.Copy(imagePath, filepath.Join(path, filepath.Base(imagePath))); err != nil {
				log.SubContext("new").SyncFatal(err)
				os.Exit(1)
			}

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
			if err := engine.WriteYAML(filepath.Join(path, constants.FileMeta), meta); err != nil {
				log.SubContext("new").SyncFatal(err)
				os.Exit(1)
			}
		},
	}

	title = cmd.Flags().String("title", "", "The title (optional, will default to the file name)")
	location = cmd.Flags().String("location", "", "The location (optional)")
	posted = cmd.Flags().String("posted", "", "The posted effective date (optional)")
	tags = cmd.Flags().StringArray("tag", nil, "Photo tags (optional)")
	return cmd
}
