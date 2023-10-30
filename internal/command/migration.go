package command

import (
	"context"

	"github.com/pangum/migration/internal/core"
	"github.com/pangum/pangu/runtime"
)

type Migration struct {
	*runtime.Command

	migration *core.Migration
}

func New(migration *core.Migration) *Migration {
	return &Migration{
		Command: runtime.NewCommand("migration").Aliases("m").Usage("数据迁移").Build(),

		migration: migration,
	}
}

func (m *Migration) Run(_ *runtime.Context) error {
	return m.migration.Migrate()
}

func (m *Migration) Before(_ context.Context) error {
	return m.migration.Migrate()
}
