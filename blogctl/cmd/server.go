package cmd

import (
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	"github.com/spf13/cobra"
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
			if err := graceful.Shutdown(app); err != nil {
				log.SyncFatalExit(err)
			}
		},
	}
	bindAddr = cmd.Flags().String("bind-addr", ":9000", "The bind address for the static webserver.")
	return cmd
}
