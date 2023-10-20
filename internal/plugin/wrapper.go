package plugin

type Wrapper struct {
	Db *Config `json:"db,omitempty" yaml:"db" xml:"db" toml:"db" validate:"required"`
}
