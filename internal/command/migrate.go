package command

import (
	"github.com/pangum/logging"
	"github.com/pangum/migration/internal/core"
	"github.com/pangum/pangu/runtime"
)

type Migrate struct {
	*runtime.Command

	migration core.Migration
	logger    logging.Logger
}

func New(migration core.Migration, logger logging.Logger) *Migrate {
	return &Migrate{
		Command: runtime.NewCommand("migration").Aliases("m").Usage("数据迁移").Build(),

		migration: migration,
		logger:    logger,
	}
}

func (m *Migrate) Run(_ *runtime.Context) error {
	return m.migration.Migrate()
}
