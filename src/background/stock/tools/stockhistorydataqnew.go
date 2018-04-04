package main

import (
	//"background/stock/model"
	"background/common/logger"
	"background/stock/config"

	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Luxurioust/excelize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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


	files, err := ioutil.ReadDir(config.GetStorageRoot() + "TransData/HistoryDataNew")
	if err != nil {
		logger.Error(err)
		return
	}
	for _, f := range files {
		GetHistoryDataQNewFromExcel(f)
		return 
	}
}


func GetHistoryDataQNewFromExcel(fileName string){
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		logger.Error(err)
		return
	}
	rows := xlsx.GetRows("Sheet1")
	for i, row := range rows {
		if i == 0{
			continue
		}
		for _, value := range row {
			fmt.Print(value, "\t")
		}
		fmt.Println()
	}
}



