package global

import (
	"fmt"
	"gateway/config"
	"gateway/utils/queue"
	"gateway/utils/timer"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"sync"
)

var (
	//用户
	GVA_USER_DB *gorm.DB // 读写库

	//meme
	GVA_MEME_DB *gorm.DB // 读写库

	//GVA_REDIS *redis.ClusterClient
	GVA_REDIS *redis.Client

	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	GVA_LOG    *ZapLogger

	GVA_Timer timer.Timer = timer.NewTimerTask()

	lock   sync.RWMutex
	SvName string

	ChanQueue *queue.Queue

	QueueDataKeyMap *queue.QueueDataKeyMap

	//RatGRPCCli *grpcclient.GRPCClient
)

func NoLog(tx *gorm.DB) *gorm.DB {
	return tx.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
}

func GetListenUrl(port string) string {
	if !strings.Contains(port, ":") {
		return fmt.Sprintf("%s:%s", GVA_CONFIG.System.ListenIp, port)
	}
	return port
}

func GetConnectUrl(port string) string {
	if !strings.Contains(port, ":") {
		return fmt.Sprintf("%s:%s", GVA_CONFIG.System.ConnectIp, port)
	}
	return port
}

func GetLogLevel() logger.LogLevel {
	switch GVA_CONFIG.Mysql.GetLogMode() {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
