package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/logger"
	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/rabbitq"
	sqlstorage "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	fmt.Println(flag.Args(), configFile)
	config := NewConfig(configFile)
	fmt.Println("config=", config)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	storage := sqlstorage.New(config.DB.DBName, config.DB.DBUserName, config.DB.DBPassword)
	logg.Info("Connected to storage:", storage)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	err := storage.Connect(ctx)
	if err != nil {
		logg.Fatal(err.Error())
	}
	err = storage.DBConnect.Ping()
	if err != nil {
		logg.Fatal(err.Error())
	}
	myRQ, err := rabbitq.CreateQueue(config.Rabbit, ctx)
	if err != nil {
		logg.Fatal(err.Error())
	}
	logg.Info("Connected to Rabit! - ", myRQ)
	go func() {
		for {
			currDate := time.Now()
			currDateStr := currDate.Format("02.01.2006")
			lastYearStr := currDate.AddDate(-1, 0, 0).Format("02.01.2006")
			logg.Info("Проснулись. Сегодня ", currDateStr)

			// отпарвка оповещений
			myEventList, err2 := storage.GetEventByDate(currDateStr)
			fmt.Println(myEventList, err2)
			for _, myEvent := range myEventList {
				fmt.Println(myEvent)
				logg.Info("найдено сообщение для отправки - ", myEvent)
				myMess, errMarsh := json.Marshal(myEvent)
				if errMarsh != nil {
					logg.Error("ошибка json.Marshal", errMarsh)
				}
				if erSemdMess := myRQ.SendMess(myMess); erSemdMess != nil {
					logg.Error("ошибка отправки сообщения-", errMarsh)
				} else {
					logg.Info("сообщение успешно отпралвено")
				}
			}

			//  Удаление прошлогодних
			myEventList, err2 = storage.GetEventByDate(lastYearStr)
			fmt.Println(myEventList, err2)
			for _, myEvent := range myEventList {
				fmt.Println(myEvent)
				logg.Info("найдено событие для удаления - ", myEvent)
				if errDel := storage.DeleteEvent(myEvent.ID); errDel != nil {
					logg.Error("ошибка удаления сообщения -", errDel)
				} else {
					logg.Info("событие успешно удалено")
				}
			}
			time.Sleep(time.Hour * 24)
		}
	}()
	<-ctx.Done()
	storage.DBConnect.Close()
	myRQ.Shutdown()
	logg.Info("Finish")
}
