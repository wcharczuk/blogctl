package cmd

import (
	"github.com/blend/go-sdk/logger"
	"github.com/wcharczuk/blogctl/pkg/config"
)

// banner is the banner displayed at the beginning of most commands.
// it is generated here: http://patorjk.com/software/taag/#p=display&h=3&f=Slant&t=BLOGCTL
const banner = `    ____  __   ____  ____________________
   / __ )/ /  / __ \/ ____/ ____/_  __/ /
  / __  / /  / / / / / __/ /     / / / /
 / /_/ / /__/ /_/ / /_/ / /___  / / / /___
/_____/_____\____/\____/\____/ /_/ /_____/

`

// Logger returns a new logger.
func Logger(cfg config.Config, name string) logger.Context {
	return logger.MustNew(logger.OptConfig(cfg.Logger), logger.OptText(logger.OptTextHideTimestamp())).SubContext("blogctl").SubContext(name)
}
