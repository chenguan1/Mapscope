package config

import "path/filepath"

type Config struct {
	DataFolder string
}

var config Config

func init()  {
	df, _ := filepath.Abs("./data")
	config = Config{
		DataFolder: df,
	}
}

func Load(path string) *Config {
	return &config
}

func Get() *Config {
	return &config
}