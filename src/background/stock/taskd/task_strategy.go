package main

import (
	"flag"
	"log"

	"background/common/logger"
	"background/stock/model"
	"background/stock/config"

	"background/stock/strategy"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
	"fmt"
)

func main(){
	//fmt.Println(fmt.Sprintf("%02d%02d%02d",time.Now().Hour(),time.Now().Minute(),time.Now().Second()))

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

	db.LogMode(false)

	model.InitModel(db)

	sTwoYearsBeforeDate := fmt.Sprintf("%04d%02d%02d",time.Now().Year() - 2,time.Now().Month(),time.Now().Day())
	sNowDate := fmt.Sprintf("%04d-%02d-%02d",time.Now().Year(),time.Now().Month(),time.Now().Day())
	//只取上市时间大于2年的股票
	var stockList []*model.StockList
	if err := db.Where("timetomarket < ?" , sTwoYearsBeforeDate).Find(&stockList).Error ; err != nil{
		logger.Error(err)
		return
	}

	var process *sync.Mutex
	process = new(sync.Mutex)
	var Count int = 0
	for index,stock := range stockList{
		for{
			if Count > 30{
				time.Sleep(time.Millisecond * 100)
			}else{
				break
			}
		}
		logger.Debug("股票总数:",len(stockList)," 开始模拟第",index + 1,"只股票,进度:", float64(index + 1) / float64(len(stockList)) * 100.00,"% 代码:", stock.Code," 名称:",stock.Name)

		go func(){
			process.Lock()
			Count++
			process.Unlock()
			sDate := stock.TimeToMarket[0:4] + "-" + stock.TimeToMarket[4:6] + "-" + stock.TimeToMarket[6:8]
			strategy.LowBuyHighSell(stock.Code,sDate,sNowDate,db)
			process.Lock()
			Count--
			process.Unlock()
		}()
		time.Sleep(time.Millisecond * 100)
	}

	for{
		time.Sleep(time.Second * 60)
		
	}
}










