package main

import (
	//"background/stock/model"
	"background/common/logger"
	"background/stock/config"

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
		if err = GetHistoryDataQNewFromExcel(f.Name(),db) ; err != nil{
			logger.Error(err)
			return
		}
	}
}


func GetHistoryDataQNewFromExcel(fileName string,db *gorm.DB)(error){
	file, err := os.Open(config.GetStorageRoot() + "TransData/HistoryDataNew/" + fileName)
	if err != nil {
		logger.Error("File Name : ",fileName, " error : ", err)
		return err
	}
	defer file.Close()

	stockCode := fileName[0:6]
	reader := csv.NewReader(file)
	k := 0
	for {
		k++
		if k == 1{
			continue
		}
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("记录集错误:", err)
			return err
		}
		sql := "insert into stock_history_data_q_new(`code`,`date`,`open`,`high`,`close`,`low`,`volume`,`amount`) select '" + stockCode + "','"
		for i := 0; i < len(record); i++ {
			sql = sql + record[i] + "','"
		}
		sql = sql[0:len(sql) - 2]
		sql = sql + " from dual where not exists (select 1 from stock_history_data_q_new where `code` = '" + stockCode + "' and `date` = '" + record[0] + "');"

		if err = db.Exec(sql).Error ; err != nil{
			logger.Error(err)
			return err
		}
	}
	return nil
}



