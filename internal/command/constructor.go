package command

import (
	"github.com/pangum/migration/internal/core"
	"github.com/pangum/pangu"
)

type Constructor struct {
	// 构造方法
}

func (c *Constructor) New(migration *core.Migration) *Migration {
	return New(migration)
}

func (c *Constructor) Add(app *pangu.App, migration *Migration) error {
	return app.Add(migration)
}
