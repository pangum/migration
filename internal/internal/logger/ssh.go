package logger

import (
	"fmt"

	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
)

type Ssh struct {
	logger logging.Logger
}

func NewSsh(logger logging.Logger) *Ssh {
	return &Ssh{
		logger: logger,
	}
}

func (s *Ssh) Printf(format string, v ...any) {
	s.logger.Info("连接隧道", field.New("ssh", fmt.Sprintf(format, v...)))
}
