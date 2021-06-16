package migration

import (
	`github.com/storezhang/glog`
	`github.com/storezhang/pangu/app`
	`github.com/storezhang/pangu/command`
)

var _ app.Command = (*Migrate)(nil)

// Migrate 描述一个提供服务的命令
type Migrate struct {
	command.Base

	migration Migration
	logger    glog.Logger
}

// 创建服务命令
func newMigrate(logger glog.Logger) *Migrate {
	return &Migrate{
		Base: command.NewBase("migrate", "数据迁移", "m"),

		logger: logger,
	}
}

func (m *Migrate) SetMigration(migration Migration) {
	m.migration = migration
}

func (m *Migrate) Run(_ *app.Context) error {
	return m.migration.Migrate()
}
