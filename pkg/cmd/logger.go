package cmd

import (
	"github.com/blend/go-sdk/logger"

	"github.com/wcharczuk/blogctl/pkg/config"
)

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
