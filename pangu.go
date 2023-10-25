package migration

import (
	"github.com/pangum/migration/internal/command"
	"github.com/pangum/migration/internal/core"
	"github.com/pangum/pangu"
)

func init() {
	creator := new(command.Creator)
	pangu.New().Get().Dependency().Put(
		core.New,
		creator.New,
	).Build().Get(
		creator.Add,
	).Build().Build().Apply()
}
