package main

import (
	//"background/stock/model"
	"background/common/logger"
	"background/stock/config"

	"fmt"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
)

func main(){
	var err error
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err = config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Config Failed!!!!", err)
		return
	}

	logger.SetLevel(config.GetLoggerLevel())

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)


	files, _ := ioutil.ReadDir(config.GetStorageRoot() + "TransData/HistoryDataNew")
	for _, f := range files {
		fmt.Println(f.Name())
	}
}



