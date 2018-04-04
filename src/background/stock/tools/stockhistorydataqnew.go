package main

import (
	//"background/stock/model"
	"background/common/logger"
	"background/stock/config"

	"fmt"
	"flag"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"encoding/csv"
	"io"
	"os"
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
		GetHistoryDataQNewFromExcel(config.GetStorageRoot() + "TransData/HistoryDataNew/" + f.Name())
		return
	}
}


func GetHistoryDataQNewFromExcel(fileName string){
	file, err := os.Open(fileName)
	if err != nil {
		logger.Error("File Name : ",fileName, " error : ", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("记录集错误:", err)
			return
		}
		for i := 0; i < len(record); i++ {
			fmt.Print(record[i] + " ")
		}
		fmt.Print("\n")
	}
}



