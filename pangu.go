package migration

import (
	`github.com/storezhang/pangu`
	_ `github.com/storezhang/pangu-logging`
)

func init() {
	app := pangu.New()

	if err := app.Provides(newMigration, newMigrate); nil != err {
		panic(err)
	}
	if err := app.Invoke(func(migrate *Migrate) error {
		return app.Adds(migrate)
	}); nil != err {
		panic(err)
	}
}
