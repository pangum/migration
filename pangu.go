package migration

import (
	`github.com/pangum/pangu`
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

	var cmd *commandMigrate
	if err := app.Invoke(func(command *commandMigrate) {
		cmd = command
	}); nil != err {
		panic(err)
	}
	if err := app.Adds(cmd); nil != err {
		panic(err)
	}
}
