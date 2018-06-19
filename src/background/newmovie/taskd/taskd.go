package main

import (
	"background/newmovie/config"
	"background/newmovie/model"
	"background/newmovie/task"
	"background/common/logger"
	"flag"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {


	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	logger.SetLevel(config.GetLoggerLevel())

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)

	model.InitModel(db)

	go func(){
		for{
			task.CheckSystemPlayUrl(db)
			task.CheckOtherPlayUrl(db)

			time.Sleep(time.Minute * 10)
		}
	}()

	for {
		time.Sleep(time.Minute * 5)
	}
}
