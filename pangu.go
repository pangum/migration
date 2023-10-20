package migration

import (
	"github.com/pangum/migration/internal/command"
	"github.com/pangum/migration/internal/core"
	"github.com/pangum/pangu"
)

func init() {
	pangu.New().Get().Dependencies().Build().Provide(
		core.New,
		command.New,
		new(command.Creator).New,
	)
}
