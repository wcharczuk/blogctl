package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/web"

	"github.com/wcharczuk/blogctl/pkg/config"
)

// Server returns the server command.
func Server(flags config.Flags) *cobra.Command {
	var bindAddr *string
	var cached *bool
	var statics *[]string
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a static fileserver",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPath, err := config.ReadConfig(flags)
			Fatal(err)

			log := Logger(flags, "server")
			slant.Print(log.Output, "BLOGCTL")

			if cfgPath != "" {
				log.Infof("using config path: %s", cfgPath)
			}
			if *cached {
				log.Infof("using cached static file server")
			}

			files := cfg.OutputPathOrDefault()
			app, err := web.New(web.OptConfig(cfg.Web), web.OptBindAddr(*bindAddr), web.OptLog(log))
			Fatal(err)

			filePaths := append(*statics, files)
			log.Infof("using static search paths: %s", strings.Join(filePaths, ", "))

			if *cached {
				log.Infof("using cached static file server")
				app.ServeStaticCached("/", filePaths)
			} else {
				log.Infof("using live static file server")
				app.ServeStatic("/", filePaths)
			}

			app.SetStaticRewriteRule("/", "/$", func(filePath string, matchedPieces ...string) string {
				if len(matchedPieces) > 0 {
					return matchedPieces[0] + "index.html"
				}
				return filePath
			})
			Fatal(graceful.Shutdown(app))
		},
	}
	bindAddr = cmd.Flags().String("bind-addr", ":9000", "The bind address for the static webserver.")
	statics = cmd.Flags().StringArray("static", nil, "Alternate static directories to serve from.")
	cached = cmd.Flags().Bool("cached", false, "If we should cache static files in memory.")
	return cmd
}
