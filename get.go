package migration

import (
	"github.com/pangum/migration/internal/core"
)

var _ = Get

// Get 取得合并
func Get() *core.Migration {
	return core.New()
}
