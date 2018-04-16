package main

import (
	"background/stock/service"
	"background/stock/model"
	"background/common/logger"
	"background/stock/config"
	"background/stock/tools/util"

	"fmt"
	"flag"

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

	var holdStocks []model.HoldStockInfo
	if err = db.Where("simulation_id = 0").Find(&holdStocks).Error ; err != nil{
		logger.Error(err)
		return
	}

	var tishi string
	var allProfit float64
	for _ , holdStock := range holdStocks{
		var everyProfit float64
		jysCode := util.GetJysCodeByStockCode(holdStock.Code)
		var realTimeStock *model.RealTimeStock
		if err, realTimeStock = service.GetRealTimeStockInfoByStockCode(jysCode,holdStock.Code) ; err != nil{
			logger.Error(err)
			return
		}

		everyProfit = (realTimeStock.NowPrice - realTimeStock.YestdayClosePrice) * float64(holdStock.AllCount)
		allProfit += everyProfit
		stockName := util.GetNameByCode(holdStock.Code,db)
		tishi += stockName + " " + fmt.Sprint(everyProfit) + "\t"
	}
	tishi += "all:" + fmt.Sprint(allProfit)

	fmt.Print(tishi)
	return
}
