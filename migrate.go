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
	"github.com/goexl/gox"
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

	_config := _panguConfig.Database
	if !_config.Migration.Enable() {
		return
	}

	var migrations migrate.MigrationSource
	logger.Info("数据迁移开始", field.New("count", len(m.resources)))
	migrate.SetTable(_config.Migration.Table)
	migrate.SetIgnoreUnknown(true)

	// 开始升级数据库
	// 如果升级有错误，应退出程序
	// 不应该完成启动，导致数据库错误越来越离谱
	migrations = &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(m.resources[0]),
	}

	if err = m.setupSSH(_config, logger); nil != err {
		return
	}

	if dsn, de := _config.dsn(); nil != de {
		err = de
	} else if db, oe := sql.Open(_config.Type, dsn); nil != oe {
		err = oe
	} else if ce := m.clear(db, _config.Migration.Table, migrations); nil != ce {
		err = ce
	} else if _, ee := migrate.Exec(db, _config.Type, migrations, migrate.Up); nil != ee {
		err = ee
	} else {
		logger.Info("数据迁移成功", field.New("count", len(m.resources)))
	}

	return
}

func (m *migration) setupSSH(conf *config, logger *logging.Logger) (err error) {
	if !conf.sshEnabled() {
		return
	}

	password := conf.Password
	keyfile := conf.SSH.Keyfile
	auth := gox.If[ssh.AuthMethod]("" != password, ssh.Password(password), sshtunnel.PrivateKeyFile(keyfile))
	host := fmt.Sprintf("%s@%s", conf.Username, conf.Addr)
	tunnel := sshtunnel.NewSSHTunnel(host, auth, conf.Addr, "65513")
	tunnel.Log = newSSHLogger(logger)
	go func() {
		err = tunnel.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	conf.Addr = fmt.Sprintf("127.0.0.1:%d", tunnel.Local.Port)

	return
}

func (m *migration) clear(db *sql.DB, table string, ms migrate.MigrationSource) (err error) {
	var migrations []*migrate.Migration
	if migrations, err = ms.FindMigrations(); nil != err {
		return
	}

	migrateIds := make([]string, 0, len(migrations))
	for _, __migration := range migrations {
		migrateIds = append(migrateIds, fmt.Sprintf("'%s'", __migration.Id))
	}

	ids := strings.Join(migrateIds, ",")
	if _, err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id NOT IN(%s)", table, ids)); nil != err {
		// 表不存在不需要清理
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if noSuchTable == mysqlErr.Number {
				err = nil
			}
		}
	}

	return
}
