package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type LoggerConf struct {
	LogFile string
	Level   string
}

type HTTPConf struct {
	Host string
	Port string
}

type GRPCConf struct {
	Host string
	Port string
}

type StorageConf struct {
	StorageType string
}

type DBConf struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUserName string
	DBPassword string
}

type Config struct {
	Logger  LoggerConf
	HTTP    HTTPConf
	GRPC    GRPCConf
	Storage StorageConf
	DB      DBConf
}

func NewConfig(configFile string) (Config, error) {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	return myConf, err
}
