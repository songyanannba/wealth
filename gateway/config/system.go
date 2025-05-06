package config

import (
	"fmt"
	"strings"
)

type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`                                  // 环境值
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`                               // 端口值
	DbType        string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                      // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType       string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                   // Oss类型
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	UseRedis      bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                // 使用redis
	LimitCountIP  int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP   int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	RouterPrefix  string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	ApiLog        bool   `mapstructure:"api-log" json:"api-log" yaml:"api-log"` // 是否记录api日志

	ListenIp       string `mapstructure:"listen-ip" json:"listen-ip" yaml:"listen-ip"`                   // 监听IP
	ConnectIp      string `mapstructure:"connect-ip" json:"connect-ip" yaml:"connect-ip"`                // 连接IP
	ApiAddr        string `mapstructure:"api-addr" json:"api-addr" yaml:"api-addr"`                      // api端口值
	MasterAddr     string `mapstructure:"master-addr" json:"master-addr" yaml:"master-addr"`             // master地址
	GameAddr       string `mapstructure:"game-addr" json:"game-addr" yaml:"game-addr"`                   // game地址
	BindAddr       string `mapstructure:"bind-addr" json:"bind-addr" yaml:"bind-addr"`                   // bind地址
	GateAddr       string `mapstructure:"gate-addr" json:"gate-addr" yaml:"gate-addr"`                   // gate地址
	BackendAddr    string `mapstructure:"backend-addr" json:"backend-addr" yaml:"backend-addr"`          // backend地址
	ApiDomain      string `mapstructure:"api-domain" json:"api-domain" yaml:"api-domain"`                // api域名
	GameDomain     string `mapstructure:"game-domain" json:"game-domain" yaml:"game-domain"`             // game域名
	GamePPDomain   string `mapstructure:"game-pp-domain" json:"game-pp-domain" yaml:"game-pp-domain"`    // game pp 域名
	GamePGDomain   string `mapstructure:"game-pg-domain" json:"game-pg-domain" yaml:"game-pg-domain"`    // game pg 域名
	StorageDomain  string `mapstructure:"storage-domain" json:"storage-domain" yaml:"storage-domain"`    // storage域名
	Migrate        bool   `mapstructure:"migrate" json:"migrate" yaml:"migrate"`                         // 是否自动迁移
	ConnectCluster bool   `mapstructure:"connect-cluster" json:"connect-cluster" yaml:"connect-cluster"` // 是否连接集群
	WsPath         string `mapstructure:"ws-path" json:"ws-path" yaml:"ws-path"`                         // websocket路径
	TestApiUrl     string `mapstructure:"test-api-url" json:"test-api-url" yaml:"test-api-url"`          // 测试的api地址

	Clusters []Cluster `mapstructure:"clusters" json:"clusters" yaml:"clusters"`
}

type Cluster struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Ip       string `mapstructure:"ip" json:"ip" yaml:"ip"`
	WsScheme string `mapstructure:"ws-scheme" json:"ws-scheme" yaml:"ws-scheme"`
}

func (c Cluster) GetUrl(port string) string {
	if strings.Contains(port, ":") {
		return port
	}
	return fmt.Sprintf("%s:%s", c.Ip, port)
}

func (c Cluster) GetIpByPortUrl(portUrl string) string {
	_, port, _ := strings.Cut(portUrl, ":")
	return fmt.Sprintf("%s:%s", c.Ip, port)
}
