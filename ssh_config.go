package migration

type sshConfig struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 地址
	Addr string `json:"addr" yaml:"addr" xml:"addr" toml:"addr" validate:"required,hostname_port|hostname"`
	// 用户名
	Username string `json:"username" yaml:"username" xml:"username" toml:"username"`
	// 密码
	Password string `json:"password" yaml:"password" xml:"password" toml:"password" validate:"required_without=Keyfile"`
	// 私钥文件地址
	Keyfile string `json:"keyfile" yaml:"keyfile" xml:"keyfile" toml:"keyfile" validate:"required_without=Password"`
}

func (sc *sshConfig) Enable() bool {
	return nil == sc.Enabled || *sc.Enabled
}
