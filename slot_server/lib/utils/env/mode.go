package env

import "os"

type ModeType string

var (
	Dev   ModeType = "dev"
	Test  ModeType = "test"
	Debug ModeType = "debug"
	Pre   ModeType = "pre"
	Prod  ModeType = "prod"
)

var ModeLevel = map[ModeType]int{
	Dev:   1,
	Test:  2,
	Debug: 3,
	Pre:   4,
	Prod:  5,
}

var Mode ModeType

func init() {
	Mode = ModeType(os.Getenv("RUN_ENV"))
	if Mode == "" {
		Mode = Dev
	}
}

// LtMode 当前模式级别是否小于等于指定模式
func LtMode(mode ModeType) bool {
	return ModeLevel[Mode] <= ModeLevel[mode]
}

// GtMode 当前模式级别是否大于等于指定模式
func GtMode(mode ModeType) bool {
	return ModeLevel[Mode] >= ModeLevel[mode]
}

func GetConfigFileName() string {
	//switch Mode {
	//case Dev:
	//	return "config.yaml"
	//case Test, Debug:
	//	return "config.dev.yaml"
	//case Pre:
	//	return "config.pre.yaml"
	//case Prod:
	//	return "config.prod.yaml"
	//}
	return "app.yaml"
}
