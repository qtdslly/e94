package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"flag"
	"background/stock/config"
	"log"
	"background/common/systemcall"
	"github.com/jinzhu/gorm"
	"background/common/logger"
	"background/stock/model"
	"background/stock/service"

	_ "github.com/go-sql-driver/mysql"
)

func main(){

	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	// set file descriptor limit
	systemcall.SetFileLimit()

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.DB().SetMaxIdleConns(10)

	model.InitModel(db)

	url := "http://hq.sinajs.cn/list=sh600570"
	resp, _ := http.Get(url)

	data, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(data))

	err ,realTimeStock := service.GetRealTimeStockObject("sh","600570",string(data))

	if err != nil{
		logger.Error("================")
		return
	}
	logger.Printf(fmt.Sprint(realTimeStock.YestdayClosePrice))
	db.Save(&realTimeStock)
}
