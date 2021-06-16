package migration

type panguConfig struct {
	// 关系型数据库配置
	Database config `json:"database" yaml:"database" validate:"required"`
}
