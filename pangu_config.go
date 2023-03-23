package migration

type wrapper struct {
	// 关系型数据库配置
	Db *config `json:"db" yaml:"db" xml:"db" toml:"db" validate:"required"`
}
