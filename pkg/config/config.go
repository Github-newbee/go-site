package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func NewConfig(p string) *viper.Viper {
	// 获取当前操作系统环境变量
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}
