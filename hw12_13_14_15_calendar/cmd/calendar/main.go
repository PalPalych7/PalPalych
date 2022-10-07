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
	internalhttpGRPC "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/server/HTTP_GRPC"
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
		storage := sqlstorage.New(config.DB.DBName, config.DB.DBUserName, config.DB.DBPassword)
		logg.Info("Get new storage:", storage)
		err := storage.Connect(context.Background())
		if err != nil {
			logg.Fatal(err.Error())
		}
		err = storage.DBConnect.Ping()
		if err != nil {
			logg.Fatal(err.Error())
		}
		calendar = app.New(logg, storage)
	}
	calendar.Logg.Info("Get new calendar:", calendar)

	server := internalhttp.NewServer(calendar, config.HTTP.Host+":"+config.HTTP.Port)
	serverGRPC := internalhttpGRPC.NewServer(calendar, config.GRPC.Host+":"+config.GRPC.Port)

	calendar.Logg.Info("server:", server)
	calendar.Logg.Info("serverGRPC:", serverGRPC)

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
		if err := serverGRPC.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	go func() {
		if err := serverGRPC.Start(ctx); err != nil {
			logg.Error("failed to start GRPC server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	logg.Info("Finish")
}
