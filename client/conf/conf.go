package conf

import "github.com/jinzhu/configor"

//const HOST = "127.0.0.1:8089"
//const PATH = "/zs_game"

const HOST = "127.0.0.1:8081"
//const HOST = "47.97.201.179:8081"
const PATH = "/gate_way"

var CliConf struct {
	Host string `default:"127.0.0.1:8189"`
	Port string `default:"8089"`
}

func CliConfInit() {

	configor.Load(&CliConf, "")
}
