package main

import (
	"background/stock/model"
	"background/stock/config"
	"background/common/logger"
	apiths "background/stock/api/tonghuashun"
	"log"
	"flag"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)
func main(){

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
		logger.Print(config.GetDBName())
		logger.Print(config.GetDBSource())

		logger.Fatal("Open db Failed!!!!", err)


		return
	}

	db.LogMode(true)

	model.InitModel(db)

	var stocks []model.StockList
	if err = db.Find(&stocks).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,stock := range stocks{
		//apiths.GetComprehensive(stock.Code,db)
		if apiths.GetControlInfo(stock.Code,db) != nil{
			continue
		}
	}
}
