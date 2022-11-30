package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/logger"
	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/rabbitq"
	sqlstorage "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/sender_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	fmt.Println(flag.Args(), configFile)
	config, err := NewConfig(configFile)
	fmt.Println("config=", config, err)
	if err != nil {
		return
	}
	fmt.Println("config=", config)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	logg.Info("Start!")
	storage := sqlstorage.New(config.DB.DBName, config.DB.DBUserName, config.DB.DBPassword, config.DB.DBHost, config.DB.DBPort) //nolint
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	err = storage.Connect(ctx)
	if err != nil {
		fmt.Println(1, err)
		logg.Fatal(err.Error())
	}
	err = storage.DBConnect.Ping()
	if err != nil {
		fmt.Println(2, err)
		logg.Fatal(err.Error())
	}
	logg.Info("Connected to storage")

	myRQ, err := rabbitq.CreateQueue(config.Rabbit, ctx)
	if err != nil {
		logg.Fatal(err.Error())
	}
	logg.Info("Connected to Rabit! - ", myRQ)

	msgs := myRQ.Consume()
	logg.Println("start consuming...")

	for m := range msgs {
		logg.Info("receive new message: ", string(m))
		// запишем информацию о получении сообщения в БД
		if err3 := storage.SendMessStat(string(m)); err3 != nil {
			logg.Error("SendMessStat error -", err3)
		} else {
			logg.Info("successful SendMessStat")
		}

	}
	myRQ.Shutdown()
	logg.Info("Finish")
}
