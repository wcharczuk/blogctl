package config

import (
	"runtime"

	"github.com/spf13/cobra"
)

// NewFlags returns a new persistent flags collection.
func NewFlags(cmd *cobra.Command) Flags {
	return Flags{
		ConfigPath:  cmd.PersistentFlags().StringP("config", "c", "config.yml", "The config path to use for cli settings"),
		DryRun:      cmd.PersistentFlags().Bool("dry-run", false, "If we should only print the plan, and not realize changes"),
		LoggerFlags: cmd.PersistentFlags().StringArray("log-flag", []string{"all", "-debug"}, "The logger flags to use with the cli"),
		Parallelism: cmd.PersistentFlags().IntP("parallelism", "p", runtime.NumCPU(), "The parallelism settings"),
	}
}

// Flags returns the commandline flag config settings.
type Flags struct {
	DryRun      *bool
	ConfigPath  *string
	LoggerFlags *[]string
	Parallelism *int
}
