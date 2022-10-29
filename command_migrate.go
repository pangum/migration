package migration

import (
	"github.com/pangum/logging"
	"github.com/pangum/pangu/app"
	"github.com/pangum/pangu/cmd"
)

var _ app.Command = (*commandMigrate)(nil)

// 数据迁移执行命令
type commandMigrate struct {
	cmd.Command

	migrate migration
	logger  *logging.Logger
}

// 创建数据迁移命令
func newCommandMigrate(logger *logging.Logger) *commandMigrate {
	return &commandMigrate{
		Command: *cmd.New(`migrate`, cmd.Aliases(`m`), cmd.Description(`数据迁移`)),
		logger:  logger,
	}
}

func (m *commandMigrate) SetMigration(migration migration) {
	m.migrate = migration
}

func (m *commandMigrate) Run(_ *app.Context) error {
	return m.migrate.migration()
}
