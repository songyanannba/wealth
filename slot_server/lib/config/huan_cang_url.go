package config

type HuanCangUrl struct {
	GetUserScoreUrl    string `mapstructure:"getUserScoreUrl" json:"getUserScoreUrl" yaml:"getUserScoreUrl"`
	UpdateUserScoreUrl string `mapstructure:"updateUserScoreUrl" json:"updateUserScoreUrl" yaml:"updateUserScoreUrl"`
	AddUserScoreUrl    string `mapstructure:"addUserScoreUrl" json:"addUserScoreUrl" yaml:"addUserScoreUrl"`
	GetIntegralProduct string `mapstructure:"getIntegralProduct" json:"getIntegralProduct" yaml:"getIntegralProduct"`
	GetUserInfo        string `mapstructure:"getUserInfo" json:"getUserInfo" yaml:"getUserInfo"`
}
