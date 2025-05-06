package config

type Server struct {
	JWT    JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	System System `mapstructure:"system" json:"system" yaml:"system"`

	// gorm
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`

	//meme battle
	MemeBattleMysql Mysql `mapstructure:"mysqlMemeBattle" json:"mysqlMemeBattle" yaml:"mysqlMemeBattle"`

	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	//Local Local `mapstructure:"local" json:"local" yaml:"local"`

	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	// 密钥配置
	Keys Keys `mapstructure:"keys" json:"keys" yaml:"keys"`

	HuanCangUrl HuanCangUrl `mapstructure:"huanCangUrl" json:"huanCangUrl" yaml:"huanCangUrl"`
}
