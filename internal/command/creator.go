package command

import (
	"github.com/pangum/pangu"
)

type Creator struct {
	// 用于解决命名空间问题
}

func (c *Creator) New(app *pangu.App, migrate *Migrate) error {
	return app.Add(migrate)
}
