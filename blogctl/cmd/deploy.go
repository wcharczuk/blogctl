package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
)

// Deploy returns the deploy command.
func Deploy(configPath *string, log *logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}
		},
	}

	return cmd
}
