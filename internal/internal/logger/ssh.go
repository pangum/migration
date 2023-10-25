package logger

import (
	"fmt"

	"github.com/goexl/gox/field"
	"github.com/goexl/log"
)

type Ssh struct {
	logger log.Logger
}

func NewSsh(logger log.Logger) *Ssh {
	return &Ssh{
		logger: logger,
	}
}

func (s *Ssh) Printf(format string, v ...any) {
	s.logger.Info("连接隧道", field.New("ssh", fmt.Sprintf(format, v...)))
}
