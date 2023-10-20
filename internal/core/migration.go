package core

import (
	"database/sql"
	"errors"
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
	"github.com/goexl/simaqian"
	"github.com/pangum/migration/internal/internal/constant"
	"github.com/pangum/migration/internal/internal/logger"
	"github.com/pangum/migration/internal/plugin"
	"github.com/pangum/pangu"
	"github.com/rubenv/sql-migrate"
	"golang.org/x/crypto/ssh"
)

var (
	once     sync.Once
	instance *Migration
)

type Migration struct {
	resources []fs.FS
}

func New() *Migration {
	once.Do(func() {
		instance = &Migration{
			resources: make([]fs.FS, 0),
		}
	})

	return instance
}

func (m *Migration) Add(migration fs.FS) {
	m.resources = append(m.resources, migration)
}

func (m *Migration) Migrate() error {
	return gox.If(0 != len(m.resources), pangu.New().Get().Dependency().Build().Get(m.New))
}

func (m *Migration) New(config *pangu.Config, logger simaqian.Logger) (err error) {
	wrapper := new(plugin.Wrapper)
	if ge := config.Build().Get(wrapper); nil != ge {
		err = ge
	} else {
		err = m.new(wrapper.Db, logger)
	}

	return
}

func (m *Migration) new(config *plugin.Config, logger simaqian.Logger) (err error) {
	if !config.Migration.Enable() {
		return
	}

	migrations := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(m.resources[0]),
	}
	logger.Info("数据迁移开始", field.New("count", len(m.resources)))
	migrate.SetTable(config.Migration.Table)
	migrate.SetIgnoreUnknown(true)

	if ese := m.enableSSH(config, logger); nil != ese {
		err = ese
	} else if dsn, de := config.Dsn(); nil != de {
		err = de
	} else if db, oe := sql.Open(config.Type, dsn); nil != oe {
		err = oe
	} else if ce := m.clear(db, config.Migration.Table, migrations); nil != ce {
		err = ce
	} else if _, ee := migrate.Exec(db, config.Type, migrations, migrate.Up); nil != ee {
		err = ee
	} else {
		logger.Info("数据迁移成功", field.New("count", len(m.resources)))
	}

	return
}

func (m *Migration) enableSSH(conf *plugin.Config, external simaqian.Logger) (err error) {
	if !conf.SshEnabled() {
		return
	}

	password := conf.Password
	keyfile := conf.SSH.Keyfile
	auth := gox.Ift("" != password, ssh.Password(password), sshtunnel.PrivateKeyFile(keyfile))
	host := fmt.Sprintf("%s@%s", conf.Username, conf.Addr)
	if tunnel, ne := sshtunnel.NewSSHTunnel(host, auth, conf.Addr, "65513"); nil != ne {
		err = ne
	} else {
		tunnel.Log = logger.NewSsh(external)
		go m.startTunnel(tunnel)
		time.Sleep(100 * time.Millisecond)
		conf.Addr = fmt.Sprintf("127.0.0.1:%d", tunnel.Local.Port)
	}

	return
}

func (m *Migration) clear(db *sql.DB, table string, source migrate.MigrationSource) (err error) {
	if migrations, fme := source.FindMigrations(); nil != fme {
		err = fme
	} else {
		err = m.delete(db, table, migrations)
	}

	return
}

func (m *Migration) delete(db *sql.DB, table string, migrations []*migrate.Migration) (err error) {
	ids := make([]string, 0, len(migrations))
	for _, migration := range migrations {
		ids = append(ids, fmt.Sprintf("'%s'", migration.Id))
	}

	exec := fmt.Sprintf("DELETE FROM %s WHERE id NOT IN(%s)", table, strings.Join(ids, ","))
	if _, err = db.Exec(exec); nil != err {
		// 表不存在不需要清理
		mysqlError := new(mysql.MySQLError)
		if errors.As(err, &mysqlError) && constant.NoSuchTableCode == mysqlError.Number {
			err = nil
		}
	}

	return
}

func (m *Migration) startTunnel(tunnel *sshtunnel.SSHTunnel) {
	_ = tunnel.Start()
}
