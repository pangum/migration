package migration

import (
	"fmt"
	"strings"

	"github.com/goexl/exc"
	"github.com/goexl/gox/field"
)

type config struct {
	// 数据库类型
	// nolint: lll
	Type string `default:"mysql" json:"type" yaml:"type" xml:"type" toml:"type" validate:"required,oneof=mysql sqlite3 mssql oracle psql"`

	// 地址，填写服务器地址
	Addr string `default:"127.0.0.1:3306" json:"addr" yaml:"addr" xml:"addr" toml:"addr" validate:"required,hostname_port"`
	// 授权，用户名
	Username string `json:"username,omitempty" yaml:"username" xml:"username" toml:"username"`
	// 授权，密码
	Password string `json:"password,omitempty" yaml:"password" xml:"password" toml:"password"`
	// 连接协议
	// nolint: lll
	Protocol string `default:"tcp" json:"protocol" yaml:"protocol" xml:"protocol" toml:"protocol" validate:"required,oneof=tcp udp"`

	// 连接的数据库名
	Schema string `json:"schema" yaml:"schema" xml:"schema" toml:"schema" validate:"required"`

	// 额外参数
	Parameters string `json:"parameters,omitempty" yaml:"parameters" xml:"parameters" toml:"parameters"`
	// SQLite填写数据库文件的路径
	Path string `default:"data.db" json:"path,omitempty" yaml:"path" xml:"path" toml:"path"`

	// SSH代理连接
	SSH *sshConfig `json:"ssh" yaml:"ssh" xml:"ssh" toml:"ssh"`

	// 数据迁移配置
	Migration migrationConfig `json:"migrate" yaml:"migrate" xml:"migration" toml:"migration" validate:"required"`
}

func (c *config) dsn() (dsn string, err error) {
	switch strings.ToLower(c.Type) {
	case `mysql`:
		dsn = fmt.Sprintf("%s:%s@%s(%s)", c.Username, c.Password, c.Protocol, c.Addr)
		if `` != strings.TrimSpace(c.Schema) {
			dsn = fmt.Sprintf("%s/%s", dsn, strings.TrimSpace(c.Schema))
		}
	case "sqlite3":
		dsn = c.Path
	default:
		err = exc.NewField("不支持的数据库类型", field.New("type", c.Type))
	}
	if nil != err {
		return
	}

	// 增加参数
	if "" != strings.TrimSpace(c.Parameters) {
		dsn = fmt.Sprintf("%s?%s", dsn, strings.TrimSpace(c.Parameters))
	}

	return
}
