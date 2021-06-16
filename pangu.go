package migration

import (
	`github.com/storezhang/pangu`
	_ `github.com/storezhang/pangu-logging`
)

func init() {
	if err := pangu.New().Provides(newMigration, newMigrate); nil != err {
		panic(err)
	}
}
