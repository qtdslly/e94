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

	//go func(){
	//	for{
	//		var p = time.Now()
	//		if (fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) >= "0930" && fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) <= "1130") ||
	//			(fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) >= "1300" && fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) <= "1500"){
	//			go task.TransPromptAll(db)
	//		}
	//		time.Sleep(time.Minute)
	//	}
	//}()

	//go func(){
	//	for{
	//		var p = time.Now()
	//		if fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) > "1500"{
	//			go task.SyncAllRealTimeStockInfo(db)
	//		}
	//		time.Sleep(time.Hour * 3)
	//	}
	//}()
	//
	go func(){
		for{
			var p = time.Now()
			if fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) < "1500"{
				go task.GetTonghuashun(db)
			}
			time.Sleep(time.Hour)
		}
	}()

	//task.GetLargeFallStockInfo(db)

	for {
		time.Sleep(time.Minute * 5)
	}
}
