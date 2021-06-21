package migration

import (
	`github.com/storezhang/pangu`
	_ `github.com/storezhang/pangu-logging`
)

func init() {
	app := pangu.New()
	migrate := New()

	if err := app.Adds(migrate); nil != err {
		panic(err)
	}
	if err := app.Provides(newCommandMigrate); nil != err {
		panic(err)
	}
	if err := app.Invoke(func(command *commandMigrate) error {
		return app.Adds(command)
	}); nil != err {
		panic(err)
	}
}
