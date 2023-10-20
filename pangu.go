package migration

import (
	"github.com/pangum/migration/internal/command"
	"github.com/pangum/migration/internal/core"
	"github.com/pangum/migration/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	pangu.New().Get().Dependencies().Build().Provide(
		core.New,
		command.New,
		new(plugin.Creator).New,
	)
}
