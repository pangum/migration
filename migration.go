package migration

import (
	`database/sql`
	`io/fs`
	`net/http`

	`github.com/go-sql-driver/mysql`
	`github.com/rubenv/sql-migrate`
	`github.com/storezhang/glog`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/pangu`
	`xorm.io/builder`
	`xorm.io/xorm`
)

const noSuchTable = 1146

type Migration struct {
	// 文件名称
	Id string `xorm:"varchar(64) notnull default('')"`
	// 升级时间
	AppliedAt gox.Timestamp `xorm:"created default('2020-02-04 09:55:52')"`

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

	if err = m.cleanDeletedMigrations(migrations, engine); nil != err {
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

func (m *Migration) cleanDeletedMigrations(ms migrate.MigrationSource, engine *xorm.Engine) (err error) {
	var (
		migrates     []*migrate.Migration
		migrateFiles = make([]string, 0)
	)

	if migrates, err = ms.FindMigrations(); nil != err {
		return
	}
	for _, m := range migrates {
		migrateFiles = append(migrateFiles, m.Id)
	}

	cond := builder.NotIn("id", migrateFiles)
	if _, err = engine.Where(cond).Delete(&Migration{}); nil != err {
		// 表不存在不需要清理
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if noSuchTable == mysqlErr.Number {
				err = nil
			}
		}
	}

	return
}
