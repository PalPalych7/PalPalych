package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/logger"
	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/rabbitq"
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
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	myRQ, err := rabbitq.CreateQueue(config.Rabbit, ctx)
	if err != nil {
		logg.Fatal(err.Error())
	}
	logg.Info("Connected to Rabit! - ", myRQ)

	msgs := myRQ.Consume()
	logg.Println("start consuming...")

	for m := range msgs {
		logg.Println("receive new message: ", string(m))
	}
	myRQ.Shutdown()
	logg.Info("Finish")
}
