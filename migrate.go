package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
	"github.com/pangum/pangu/app"
	"github.com/rubenv/sql-migrate"
)

const noSuchTable = 1146

var (
	_migration *migration
	once       sync.Once
	_          app.Executor = (*migration)(nil)
)

type migration struct {
	resources []fs.FS
}

// New 创建新的数据迁移
func New() *migration {
	once.Do(func() {
		_migration = &migration{
			resources: make([]fs.FS, 0),
		}
	})

	return _migration
}

func (m *migration) migration() (err error) {
	if 0 == len(m.resources) {
		return
	}

	err = pangu.New().Invoke(func(config *pangu.Config, logger *logging.Logger) (err error) {
		_panguConfig := new(panguConfig)
		if err = config.Load(_panguConfig); nil != err {
			return
		}
		database := _panguConfig.Database

		var migrations migrate.MigrationSource
		logger.Info("数据迁移开始", field.Int("count", len(m.resources)))
		// 设置升级记录的表名，默认值是grop_migrations
		migrate.SetTable(database.Migration.Table)
		migrate.SetIgnoreUnknown(true)

		// 开始升级数据库
		// 如果升级有错误，应退出程序
		// 不应该完成启动，导致数据库错误越来越离谱
		migrations = &migrate.HttpFileSystemMigrationSource{
			FileSystem: http.FS(m.resources[0]),
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
		logger.Info("数据迁移成功", field.Int("count", len(m.resources)))

		return
	})

	return
}

func (m *migration) AddSource(migration fs.FS) {
	m.resources = append(m.resources, migration)
}

func (m *migration) Run() (err error) {
	return m.migration()
}

func (m *migration) Name() string {
	return `数据迁移`
}

func (m *migration) Type() app.ExecutorType {
	return app.ExecutorTypeBeforeServe
}

func (m *migration) ExecuteType() app.ExecuteType {
	return app.ExecuteTypeReturn
}

func (m *migration) clear(db *sql.DB, table string, ms migrate.MigrationSource) (err error) {
	var migrations []*migrate.Migration
	if migrations, err = ms.FindMigrations(); nil != err {
		return
	}

	migrateIds := make([]string, 0, len(migrations))
	for _, migration := range migrations {
		migrateIds = append(migrateIds, fmt.Sprintf("'%s'", migration.Id))
	}

	if _, err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id NOT IN(%s)", table, strings.Join(migrateIds, ", "))); nil != err {
		// 表不存在不需要清理
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if noSuchTable == mysqlErr.Number {
				err = nil
			}
		}
	}

	return
}
