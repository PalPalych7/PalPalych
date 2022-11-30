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
	config, err := NewConfig(configFile)
	fmt.Println("config=", config, err)
	if err != nil {
		return
	}
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	storage := sqlstorage.New(config.DB.DBName, config.DB.DBUserName, config.DB.DBPassword, config.DB.DBHost, config.DB.DBPort) //nolint
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	err = storage.Connect(ctx)
	if err != nil {
		logg.Fatal(err.Error())
	}
	err = storage.DBConnect.Ping()
	if err != nil {
		logg.Fatal(err.Error())
	}
	logg.Info("Connected to storage")
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
			logg.Debug("I not sleep:)", currDateStr)

			// отпрaвка оповещений (ещё не отправленных)
			myEventList, err2 := storage.GetNotSendEventByDate(currDateStr)

			fmt.Println(myEventList, err2)
			for _, myEvent := range myEventList {
				fmt.Println(myEvent)
				logg.Debug("find message for sending - ", myEvent)
				myMess, errMarsh := json.Marshal(myEvent)
				if errMarsh != nil {
					logg.Error("error json.Marshal", errMarsh)
				}
				if erSemdMess := myRQ.SendMess(myMess); erSemdMess != nil {
					logg.Error("send message error-", errMarsh)
				} else {
					logg.Info("message was successful sended")
					// внесём сообщение в список "отправленных"
					if err3 := storage.SetSendMessID(myEvent.ID); err3 != nil {
						logg.Error("set information about sendung error -", err3)
					}
				}
			}

			//  Удаление прошлогодних
			myEventList, err2 = storage.GetEventByDate(lastYearStr)
			fmt.Println(myEventList, err2)
			for _, myEvent := range myEventList {
				fmt.Println(myEvent)
				logg.Debug("found event for delete - ", myEvent)
				if errDel := storage.DeleteEvent(myEvent.ID); errDel != nil {
					logg.Error("delete error -", errDel)
				} else {
					logg.Info("event was successful deleted")
				}
			}
			time.Sleep(time.Second * time.Duration(config.Rabbit.SleepSecond))
		}
	}()
	<-ctx.Done()
	storage.DBConnect.Close()
	myRQ.Shutdown()
	logg.Info("Finish")
}
