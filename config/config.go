package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func New[ConfigType any](args ...string) (ConfigType, error) {
	configFile := "config"
	if len(args) > 0 && args[0] != "" {
		configFile = args[0]
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	dir, err := GetLaunchDir()
	if err == nil {
		viper.AddConfigPath(dir)
	}

	var conf ConfigType
	err = viper.ReadInConfig()
	if err != nil {
		return conf, err
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
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
