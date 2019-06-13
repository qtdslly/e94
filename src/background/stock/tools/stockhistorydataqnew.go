package main

import (
	"background/common/logger"
	"background/stock/config"

	"io"
	"os/exec"
	"os"
	"sync"
	"time"
	"flag"
	"io/ioutil"
	"encoding/csv"

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

	times := 0
	for{
		times++
		files, err := ioutil.ReadDir(config.GetStorageRoot() + "TransData/HistoryDataNew")
		if err != nil {
			logger.Error(err)
			return
		}

		var process *sync.Mutex
		process = new(sync.Mutex)
		var Count int = 0
		k := 0
		for _, f := range files {
			k++

			if IsHaveDone(f.Name(),db) {
				MoveFile(config.GetStorageRoot() + "TransData/HistoryDataNew/" + f.Name(), config.GetStorageRoot() + "TransData/havedone/" + f.Name())
				continue
			}
			for{
				if Count > 20{
					time.Sleep(time.Millisecond * 100)
				}else{
					break
				}
			}
			go func(){

				process.Lock()
				Count++
				process.Unlock()
				if err = GetHistoryDataQNewFromExcel(f.Name(),db) ; err != nil{
					logger.Error(err)
					return
				}
				//MoveFile(config.GetStorageRoot() + "TransData/HistoryDataNew/" + f.Name(), config.GetStorageRoot() + "TransData/havedone/" + f.Name())
				process.Lock()
				Count--
				process.Unlock()
			}()

			time.Sleep(time.Millisecond * 100)
		}

		for{
			if Count == 0 && k == len(files){
				break
			}else{
				time.Sleep(time.Second * 10)
			}
		}

		if times == 100{
			break
		}
	}
}


func MoveFile(srcFileName, destFileName string) {
	cmdStr := "mv " + srcFileName + " " + destFileName
	logger.Debug(cmdStr)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		logger.Error("Movfile Failed:srcFileName:"+srcFileName, " destFileName:", destFileName, err)
	} else {
		logger.Debug("Movfile Success:" + srcFileName)
	}
}


func IsHaveDone(name string,db *gorm.DB) (bool){
	code := name[0:6]
	year := name[7:11] + "%"

	var Count int
	if err := db.Where("code = ? and date like ?",code,year).Table("stock_history_data_q").Count(&Count).Error ; err !=nil{
		logger.Error(err)
		return false
	}

	if Count == 0{
		logger.Debug("Count:",Count)
		return false
	}
	logger.Debug("Count:",Count)
	return true
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

		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("记录集错误:", err)
			return err
		}
		k++
		if k == 1{
			continue
		}
		sql := "insert into stock_history_data_q(`code`,`date`,`open`,`high`,`close`,`low`,`volume`,`amount`) select '" + stockCode + "','"
		for i := 0; i < len(record); i++ {
			sql = sql + record[i] + "','"
		}
		sql = sql[0:len(sql) - 2]
		sql = sql + " from dual where not exists (select 1 from stock_history_data_q where `code` = '" + stockCode + "' and `date` = '" + record[0] + "');"

		if err = db.Exec(sql).Error ; err != nil{
			logger.Error(err)
			return err
		}
	}
	return nil
}



