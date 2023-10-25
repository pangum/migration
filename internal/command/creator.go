package command

import (
	"github.com/pangum/migration/internal/core"
	"github.com/pangum/pangu"
)

type Creator struct {
	// 用于解决命名空间问题
}

func (c *Creator) New(migration *core.Migration) *Migration {
	return New(migration)
}

func (c *Creator) Add(app *pangu.App, migration *Migration) error {
	return app.Add(migration)
}
