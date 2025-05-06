package core

import (
	"fmt"
	"gateway/global"
	"gateway/utils/env"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
func Viper(path ...string) *viper.Viper {
	global.SvName = "gate_way"

	config := env.GetConfigFileName()
	fmt.Printf("您正在使用%s环境,config的路径为%s\n", env.Mode, config)

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)

		logMode := global.GVA_CONFIG.Mysql.LogMode

		if err = v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}

		// 修改日志级别
		if global.GVA_CONFIG.Mysql.LogMode != logMode {
			fmt.Println("mysql log mode changed: " + logMode + " => " + global.GVA_CONFIG.Mysql.LogMode)
			level := global.GetLogLevel()
			global.GVA_USER_DB.Logger = global.GVA_USER_DB.Logger.LogMode(level)
			if global.GVA_USER_DB != global.GVA_USER_DB {
				global.GVA_USER_DB.Logger = global.GVA_USER_DB.Logger.LogMode(level)
			}
		}
	})

	if err = v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
