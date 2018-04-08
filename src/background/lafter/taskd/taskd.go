package main

import (
	"background/lafter/config"
	"background/lafter/model"
	"background/lafter/task"
	"background/common/logger"

	"flag"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"time"
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

	var providers []model.Provider
	if err = db.Where("status = 1").Find(&providers).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,provider := range providers{
		go task.GetContentByProvider(provider,db)
	}

	for{
		time.Sleep(time.Minute * 30)
	}
}
