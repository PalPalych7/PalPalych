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
type HttpConf struct {
	Host string
	Port string
}
type StorageConf struct {
	StorageType string
}

type DBConf struct {
	DbName     string
	DbUserName string
	DbPassword string
}

type Config struct {
	Logger  LoggerConf
	Http    HttpConf
	Storage StorageConf
	DB      DBConf
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	//	fmt.Println("q=", q, "mc=", myConf)
	return myConf
}
