package cmd

import (
	"github.com/blend/go-sdk/ansi/slant"
	"github.com/blend/go-sdk/logger"
	"github.com/wcharczuk/blogctl/pkg/config"
)

// Logger returns a new logger.
func Logger(cfg config.Config, name string) logger.Context {
	log := logger.MustNew(
		logger.OptConfig(cfg.Logger),
		logger.OptText(logger.OptTextHideTimestamp()),
		logger.OptSubContext("blogctl", name),
	)

	slant.Print(log.Output, "BLOGCTL")

	return log
}
