package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wcharczuk/blogctl/pkg/engine"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// Server returns the server command.
func Server(configPath *string, log logger.Log) *cobra.Command {
	var bindAddr, static *string
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"s", "server"},
		Short:   "Start a static fileserver",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := engine.ReadConfig(*configPath)
			if err != nil {
				logger.FatalExit(err)
			}

			files := config.OutputPathOrDefault()
			app := web.New(web.OptBindAddr(*bindAddr), web.OptLog(log))
			if *static != "" {
				log.Infof("using static search path: %s", *static)
				log.Infof("using static search path: %s", files)
				app.ServeStatic("/", []string{*static, files})
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
	static = cmd.Flags().String("static", "", "An alternate static directory to serve from.")
	return cmd
}
