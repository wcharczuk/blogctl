package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/photoblog/pkg/constants"
	"github.com/wcharczuk/photoblog/pkg/engine"
	"github.com/wcharczuk/photoblog/pkg/model"
	"github.com/wcharczuk/photoblog/pkg/stringutil"
)

// New returns a new post command.
func New(configPath *string, log *logger.Logger) *cobra.Command {
	var title, location *string
	cmd := &cobra.Command{
		Use:   "new [IMAGE_PATH]",
		Short: "Create a new blog post from a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imagePath := args[0]

			config, err := ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}

			captureDate, err := engine.ExtractCaptureDate(imagePath)
			if err != nil {
				log.SyncFatalExit(err)
			}

			if *title == "" {
				*title = filepath.Base(imagePath)
			}

			path := fmt.Sprintf("%s/%s-%s", config.ImagesOrDefault(), captureDate.Format("2006-01-02"), stringutil.Slugify(*title))
			log.Infof("writing new post to %s", path)
			if err := engine.MakeDir(path); err != nil {
				log.SyncFatalExit(err)
			}
			if err := engine.Copy(imagePath, filepath.Join(path, filepath.Base(imagePath))); err != nil {
				log.SyncFatalExit(err)
			}
			meta := model.Meta{
				Title:    *title,
				Location: *location,
				Posted:   captureDate,
			}
			if err := WriteYAML(filepath.Join(path, constants.DiscoveryFileMeta), meta); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}

	title = cmd.Flags().String("title", "", "An optional title (default will be the file name)")
	location = cmd.Flags().String("location", "", "An optional location")
	return cmd
}
