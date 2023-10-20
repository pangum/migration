package config

type Migration struct {
	// 是否启用数据迁移
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 升级记录表
	Table string `default:"migration" json:"table" yaml:"table" xml:"table" toml:"table" validate:"required"`
}

func (mc Migration) Enable() bool {
	return nil == mc.Enabled || *mc.Enabled
}
