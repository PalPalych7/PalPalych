package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/app"
	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/home/palpalych/calend/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	fmt.Println(flag.Args(), configFile)

	config := NewConfig(configFile)
	fmt.Println("config=", config)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	var calendar *app.App
	if config.Storage.StorageType == "memory" {
		logg.Info("Work with memory")
		storage := memorystorage.New()
		logg.Info("Get new storage:", storage)
		calendar = app.New(logg, storage)
	} else {
		logg.Info("Work with sql")
		storage := sqlstorage.New(config.DB.DbName, config.DB.DbUserName, config.DB.DbPassword)
		logg.Info("Get new storage:", storage)
		calendar = app.New(logg, storage)
	}
	calendar.Logg.Info("Get new calendar:", calendar)

	server := internalhttp.NewServer( /*logg,*/ calendar, config.Http.Host+":"+config.Http.Port)
	fmt.Println("server=", server)

	calendar.Logg.Info("server:", server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	logg.Info("Finish")
}
