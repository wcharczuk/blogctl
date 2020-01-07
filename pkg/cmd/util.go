package cmd

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/wcharczuk/blogctl/pkg/config"
)

// Fatal prints an error and exits the process.
func Fatal(err error) {
	if err != nil {
		var opts []logger.TextOutputFormatterOption
		if env.Env().Bool("NO_COLOR") {
			opts = append(opts, logger.OptTextNoColor())
		}
		tf := logger.NewTextOutputFormatter(opts...)
		fmt.Fprintf(os.Stderr, "%s %+v", tf.FormatFlag("error", ansi.ColorRed), err)
		os.Exit(1)
	}
}

// Logger returns a new logger.
func Logger(cfg config.Flags, name string) *logger.Logger {
	log := logger.MustNew(
		logger.OptFlags(logger.NewFlags(*cfg.LoggerFlags...)),
		logger.OptText(logger.OptTextHideTimestamp()),
		logger.OptPath("blogctl", name),
	)
	if *cfg.Debug {
		log.Flags.Enable(logger.Debug)
	}
	return log
}
