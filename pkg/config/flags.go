package config

import (
	"runtime"

	"github.com/spf13/cobra"
)

// NewFlags returns a new persistent flags collection.
func NewFlags(cmd *cobra.Command) PersistentFlags {
	return PersistentFlags{
		ConfigPath:  cmd.PersistentFlags().StringP("config", "c", "config.yml", "The config path to use for cli settings"),
		LoggerFlags: cmd.PersistentFlags().StringArray("log-flag", []string{"all", "-debug"}, "The logger flags to use with the cli"),
		Parallelism: cmd.PersistentFlags().IntP("parallelism", "p", runtime.NumCPU(), "The parallelism settings"),
	}
}

// PersistentFlags returns the persistent flag values.
type PersistentFlags struct {
	ConfigPath  *string
	LoggerFlags *[]string
	Parallelism *int
}
