package global

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"slot_server/lib/config"
	"slot_server/lib/utils/queue"
	"slot_server/lib/utils/timer"
	"strings"
	"sync"
	"time"
)

var (
	//用户
	GVA_SLOT_SERVER_DB *gorm.DB //

	//GVA_REDIS *redis.ClusterClient
	GVA_REDIS *redis.Client

	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	GVA_LOG    *ZapLogger

	GVA_Timer timer.Timer = timer.NewTimerTask()

	lock     sync.RWMutex
	SvName   string
	Location *time.Location

	QueueDataKeyMap *queue.QueueDataKeyMap
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
