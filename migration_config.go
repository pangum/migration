package migration

type migrationConfig struct {
	// 是否启用数据迁移
	Enable bool `default:"true" json:"enable" yaml:"enable"`
	// 升级记录表
	Table string `default:"migration" json:"table" yaml:"table" validate:"required"`
}
