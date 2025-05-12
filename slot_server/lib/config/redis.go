package config

type Redis struct {
	DB           int    `mapstructure:"db" json:"db" yaml:"db"`                                  // redis的哪个数据库
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`                            // 服务器地址:端口
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                // 密码
	PoolSize     int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`             // 密码
	MinIdleConns int    `mapstructure:"minIdle_conns" json:"minIdle_conns" yaml:"minIdle_conns"` // 密码
}
