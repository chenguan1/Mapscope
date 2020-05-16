package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

const CONFIGFILE = "./config.toml"

func Initialize() error {
	viper.SetConfigType("toml")
	viper.SetConfigFile(CONFIGFILE)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return (fmt.Errorf("read config(%s) file failed. err: %v", CONFIGFILE, err))
	}

	return nil
}

func PathDatasets(user string) string {
	data := viper.GetString("paths.data")
	dir := filepath.Join(data, "datasets", user)
	dir, _ = filepath.Abs(dir)
	return dir
}

func PathDatasources(user string) string {
	data := viper.GetString("paths.data")
	dir := filepath.Join(data, "datasources", user)
	dir, _ = filepath.Abs(dir)
	return dir
}

func PathUploads(user string) string {
	data := viper.GetString("paths.data")
	dir := filepath.Join(data, "uploads", user)
	dir, _ = filepath.Abs(dir)
	return dir
}

func PathCaches() string {
	data := viper.GetString("paths.data")
	dir := filepath.Join(data, "caches")
	dir, _ = filepath.Abs(dir)
	return dir
}
