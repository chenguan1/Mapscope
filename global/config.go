package global

import (
	"fmt"
	"github.com/spf13/viper"
)

const CONFIGFILE = "./config.toml"

func initConfig() error {
	viper.SetConfigType("toml")
	viper.SetConfigFile(CONFIGFILE)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil{
		return (fmt.Errorf("read config(%s) file failed. err: %v", CONFIGFILE, err))
	}

	return nil
}

