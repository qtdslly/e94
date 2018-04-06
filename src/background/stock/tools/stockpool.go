package main

import (
	"background/stock/config"
	"background/stock/model"
	"background/stock/task"
	"background/stock/tools/util"
	"background/common/logger"
	"flag"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
	"fmt"
	"background/stock/service"
)

/*
将策略获取的股票与股票基本信息结合放入股票池中
*/

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

	var deepFallStocks []model.DeepFallStock
	if err = db.Find(&deepFallStocks).Error ; err != nil{
		logger.Error(err)
		return
	}

	p := time.Now()

	for _ , deepFallStock := range deepFallStocks{
		var stock model.StockList
		if err = db.Where("code = ?",deepFallStock.Code).Error ; err != nil{
			logger.Error(err)
			return
		}

		jysCode := util.GetJysCodeByStockCode(deepFallStock.Code)
		var realStock *model.RealTimeStock
		if err,realStock = service.GetRealTimeStockInfoByStockCode(jysCode,deepFallStock.Code) ; err != nil{
			logger.Error(err)
			return
		}

		/*上市日期需大于2年*/
		if stock.TimeToMarket < fmt.Sprintf("%04d%02d%02d",p.Year() - 2 ,p.Month(),p.Day()){
			continue
		}

		/*市值需小于等于1000亿*/
		if stock.Totals * realStock.NowPrice > 1000{
			continue
		}
		if stock.Rev < 10 { /*收入同比需大于等于10%*/
			continue
		}
		if stock.Npr < 10{ /*净利润率需大于等于10%*/
			continue
		}
		if stock.Profit < 10{ /*利润同比需大于等于10%*/
			continue
		}
		if stock.Gpr < 5{ /*毛利率需大于等于5%*/
			continue
		}

	}

}
