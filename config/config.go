package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseConnectionString string `mapstructure:"database_connection_string"`
	NatsConnectionString     string `mapstructure:"nats_connection_string"`
}

func New() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	dir, err := GetLaunchDir()
	if err == nil {
		viper.AddConfigPath(dir)
	}

	var conf Config
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return conf
}

func GetLaunchDir() (string, error) {
	var dir string
	excutable, err := os.Executable()
	if err != nil {
		return dir, err
	}

	dir = filepath.Dir(excutable)
	return dir, nil
}
