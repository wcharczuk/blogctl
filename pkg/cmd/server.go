package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"

	"github.com/wcharczuk/blogctl/pkg/config"
)

// Server returns the server command.
func Server(flags *config.PersistentFlags) *cobra.Command {
	var bindAddr *string
	var statics *[]string
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"s", "server"},
		Short:   "Start a static fileserver",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPath, err := config.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}

			log := Logger(cfg, "server")

			if cfgPath != "" {
				log.Infof("using config path: %s", cfgPath)
			}

			files := cfg.OutputPathOrDefault()
			app := web.New(web.OptConfig(cfg.Web), web.OptBindAddr(*bindAddr), web.OptLog(log))
			if len(*statics) > 0 {
				filePaths := append(*statics, files)
				log.Infof("using static search paths: %s", strings.Join(filePaths, ", "))
				app.ServeStatic("/", filePaths)
			} else {
				log.Infof("using static search path: %s", files)
				app.ServeStatic("/", []string{files})
			}

			app.SetStaticRewriteRule("/", "/$", func(filePath string, matchedPieces ...string) string {
				if len(matchedPieces) > 0 {
					return matchedPieces[0] + "index.html"
				}
				return filePath
			})
			if err := graceful.Shutdown(app); err != nil {
				logger.FatalExit(err)
			}
		},
	}
	bindAddr = cmd.Flags().String("bind-addr", ":9000", "The bind address for the static webserver.")
	statics = cmd.Flags().StringArray("static", nil, "Alternate static directories to serve from.")
	return cmd
}
