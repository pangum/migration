package plugin

import (
	"github.com/pangum/migration/internal/command"
	"github.com/pangum/pangu"
)

type Creator struct {
	// 用于解决命名空间问题
}

func (c *Creator) New(app *pangu.App, migrate *command.Migrate) error {
	return app.Add(migrate)
}
