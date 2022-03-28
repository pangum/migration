package migration

import (
	"fmt"

	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
)

type sshLogger struct {
	logger *logging.Logger
}

func newSSHLogger(logger *logging.Logger) *sshLogger {
	return &sshLogger{
		logger: logger,
	}
}

func (sl *sshLogger) Printf(format string, v ...interface{}) {
	sl.logger.Info(`连接隧道`, field.String(`ssh`, fmt.Sprintf(format, v...)))
}
