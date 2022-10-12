package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/rabbitq"
)

type LoggerConf struct {
	LogFile string
	Level   string
}

type DBConf struct {
	DBName     string
	DBUserName string
	DBPassword string
}

type Config struct {
	Logger LoggerConf
	DB     DBConf
	Rabbit rabbitq.RabbitCFG
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	return myConf
}
