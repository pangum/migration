package migration

import (
	`github.com/storezhang/pangu`
	_ `github.com/storezhang/pangu-logging`
)

func init() {
	app := pangu.New()

	if err := app.Provides(newCommandMigrate); nil != err {
		panic(err)
	}
	if err := app.Invoke(func(command *commandMigrate, executor *migration) error {
		return app.Adds(command, executor)
	}); nil != err {
		panic(err)
	}
}
