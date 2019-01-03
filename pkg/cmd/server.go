package cmd

import (
	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// Server returns the server command.
func Server(configPath *string, log *logger.Logger) *cobra.Command {
	var bindAddr *string
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a static fileserver",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}

			files := config.OutputPathOrDefault()
			app := web.New().WithBindAddr(*bindAddr)
			app.WithLogger(log)
			app.ServeStatic("/", files)
			app.SetStaticRewriteRule("/", "/$", func(filePath string, matchedPieces ...string) string {
				if len(matchedPieces) > 0 {
					return matchedPieces[0] + "index.html"
				}
				return filePath
			})
			if err := graceful.Shutdown(app); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
	bindAddr = cmd.Flags().String("bind-addr", ":9000", "The bind address for the static webserver.")
	return cmd
}
