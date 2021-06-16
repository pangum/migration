package migration

import (
	`database/sql`
	`io/fs`
	`net/http`
	`strings`

	`github.com/go-sql-driver/mysql`
	`github.com/rubenv/sql-migrate`
	`github.com/storezhang/glog`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/pangu`
)

const noSuchTable = 1146

type Migration struct {
	migrations []fs.FS       `xorm:"-"`
	config     *pangu.Config `xorm:"-"`
	logger     glog.Logger   `xorm:"-"`
}

func newMigration(config *pangu.Config, logger glog.Logger) *Migration {
	return &Migration{
		migrations: make([]fs.FS, 0, 0),
		config:     config,
		logger:     logger,
	}
}

func (m *Migration) Migrate() (err error) {
	if 0 == len(m.migrations) {
		return
	}

	panguConfig := new(panguConfig)
	if err = m.config.Load(panguConfig); nil != err {
		return
	}
	database := panguConfig.Database

	var migrations migrate.MigrationSource
	m.logger.Info("数据迁移开始", field.Int("count", len(m.migrations)))
	// 设置升级记录的表名，默认值是grop_migrations
	migrate.SetTable(database.Migration.Table)
	migrate.SetIgnoreUnknown(true)

	// 开始升级数据库
	// 如果升级有错误，应退出程序
	// 不应该完成启动，导致数据库错误越来越离谱
	migrations = &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(m.migrations[0]),
	}

	var dsn string
	if dsn, err = database.dsn(); nil != err {
		return
	}
	var db *sql.DB
	if db, err = sql.Open(database.Type, dsn); nil != err {
		return
	}
	defer func() {
		if closeErr := db.Close(); nil != closeErr {
			err = closeErr
		}
	}()

	if err = m.clear(db, database.Migration.Table, migrations); nil != err {
		return
	}
	_, err = migrate.Exec(db, database.Type, migrations, migrate.Up)
	m.logger.Info("数据迁移成功", field.Int("count", len(m.migrations)))

	return
}

func (m *Migration) addSource(migration fs.FS) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migration) shouldMigration() bool {
	return 0 != len(m.migrations)
}

func (m *Migration) clear(db *sql.DB, table string, ms migrate.MigrationSource) (err error) {
	var (
		migrations   []*migrate.Migration
		migrateFiles = make([]string, 0)
	)

	if migrations, err = ms.FindMigrations(); nil != err {
		return
	}
	for _, migration := range migrations {
		migrateFiles = append(migrateFiles, migration.Id)
	}

	var stmt *sql.Stmt
	if stmt, err = db.Prepare("DELETE FROM ? WHERE id NOT IN(?)"); nil != err {
		return
	}
	if _, err = stmt.Exec(table, strings.Join(migrateFiles, ",")); nil != err {
		// 表不存在不需要清理
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if noSuchTable == mysqlErr.Number {
				err = nil
			}
		}
	}

	return
}
