package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// TODO(Mitch): Add contraint for ConfigType such that only structs are valid generic types
func New[ConfigType any](args ...any) (*ConfigType, error) {
	viper.SetEnvPrefix("habits")

	pflag.String("config", "config", "Path to config file")
	viper.BindEnv("config_file")
	viper.SetDefault("config_file", "config")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	for i, arg := range args {
		switch argument := arg.(type) {
		case ConfigType:
			SetDefaultConfig(&argument)
		case *ConfigType:
			SetDefaultConfig(argument)
		default:
			return nil, fmt.Errorf("Invalid argument type: argument %d is type '%s'", i, reflect.TypeOf(argument))
		}
	}

	viper.SetConfigName(viper.GetString("config_file"))
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// dir, err := GetLaunchDir()
	// if err == nil {
	// 	viper.AddConfigPath(dir)
	// }

	return FetchConfig[ConfigType]()
}

func SetDefaultConfig[ConfigType any](defaultConf *ConfigType) {
	SetDefaultConfigStruct("", *defaultConf)
}

func SetDefaultConfigStruct(prefix string, confStruct any) {
	val := reflect.ValueOf(confStruct)
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := val.Type().Field(i)

		if !fieldVal.IsValid() {
			continue
		}

		fieldName := fieldTyp.Name
		if mapkey, ok := fieldTyp.Tag.Lookup("mapstructure"); ok {
			fieldName = mapkey
		}

		viper.SetDefault(prefix+fieldName, fieldVal.Interface())

		if fieldTyp.Type.Kind() == reflect.Struct {
			SetDefaultConfigStruct(prefix+fieldName+".", fieldVal.Interface())
		}

	}
}

func FetchConfig[ConfigType any]() (*ConfigType, error) {
	var conf ConfigType
	var err error

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
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
