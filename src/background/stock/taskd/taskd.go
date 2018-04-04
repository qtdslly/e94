package main

import (
	"background/stock/config"
	"background/stock/model"
	"background/stock/task"
	"background/common/logger"
	"flag"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
	"fmt"
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
			var p = time.Now()
			if fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) > "1500"{
				task.SyncAllRealTimeStockInfo(db)
				//task.SyncHoldStockRealTimeInfo()
			}
			time.Sleep(time.Hour * 24)
		}
	}()
	
	go task.TransPromptAll(db)
	//task.GetLargeFallStockInfo(db)
	//go func() {
	//	for {
	//		var p = time.Now()
	//		if fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) > "1200"{
	//			task.SyncAllRealTimeStockInfo(db)
	//			//task.SyncHoldStockRealTimeInfo()
	//			break
	//		}
	//
	//	}
	//}()



	for {
		time.Sleep(time.Minute * 5)
	}
}
