package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/elliotchance/sshtunnel"
	"github.com/go-sql-driver/mysql"
	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
	"github.com/pangum/pangu/app"
	"github.com/rubenv/sql-migrate"
	"golang.org/x/crypto/ssh"
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

func (m *migration) migration() (err error) {
	if 0 == len(m.resources) {
		return
	}
	err = pangu.New().Invoke(m.migrate)

	return
}

func (m *migration) migrate(config *pangu.Config, logger *logging.Logger) (err error) {
	_panguConfig := new(panguConfig)
	if err = config.Load(_panguConfig); nil != err {
		return
	}
	database := _panguConfig.Database

	if !database.Migration.Enable() {
		return
	}

	var migrations migrate.MigrationSource
	logger.Info(`数据迁移开始`, field.Int(`count`, len(m.resources)))
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
	if dsn, err = m.dsn(database, logger); nil != err {
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
	logger.Info(`数据迁移成功`, field.Int(`count`, len(m.resources)))

	return
}

func (m *migration) dsn(conf config, logger *logging.Logger) (dsn string, err error) {
	if nil == conf.SSH || !conf.SSH.Enable() {
		return
	}

	var auth ssh.AuthMethod
	if `` != conf.Password {
		auth = ssh.Password(conf.Password)
	} else {
		auth = sshtunnel.PrivateKeyFile(conf.SSH.Keyfile)
	}
	host := fmt.Sprintf(`%s@%s`, conf.Username, conf.Addr)
	tunnel := sshtunnel.NewSSHTunnel(host, auth, conf.Addr, `65513`)
	tunnel.Log = newSSHLogger(logger)
	go func() {
		err = tunnel.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	conf.Addr = fmt.Sprintf(`127.0.0.1:%d`, tunnel.Local.Port)
	dsn, err = conf.dsn()

	return
}

func (m *migration) clear(db *sql.DB, table string, ms migrate.MigrationSource) (err error) {
	var migrations []*migrate.Migration
	if migrations, err = ms.FindMigrations(); nil != err {
		return
	}

	migrateIds := make([]string, 0, len(migrations))
	for _, __migration := range migrations {
		migrateIds = append(migrateIds, fmt.Sprintf(`'%s'`, __migration.Id))
	}

	ids:=strings.Join(migrateIds, `, `)
	if _, err = db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id NOT IN(%s)`, table, ids)); nil != err {
		// 表不存在不需要清理
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if noSuchTable == mysqlErr.Number {
				err = nil
			}
		}
	}

	return
}
