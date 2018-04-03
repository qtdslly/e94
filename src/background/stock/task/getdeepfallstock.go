package task

import (
	"fmt"
	"time"

	"background/stock/model"
	"github.com/jinzhu/gorm"
	"background/common/logger"
	"background/stock/strategy"
	"background/stock/tools/util"
	"background/stock/service"
	"sync"
)

func GetLargeFallStockInfo(db *gorm.DB){
	var stockList []*model.StockList
	if err := db.Find(&stockList).Error ; err != nil{
		logger.Error(err)
		return
	}
	var process *sync.Mutex
	process = new(sync.Mutex)
	var Count int = 0
	var flag int = 0
	for _,stock := range stockList{
		for{
			if Count > 20{
				time.Sleep(time.Second)
			}else{
				break
			}
		}

		go func(){
			process.Lock()
			Count++
			flag++
			logger.Debug("开始第",flag,"只股票，共",len(stockList),"只")
			process.Unlock()
			GetLargeFallStockInfoByCode(stock.Code,db)
			process.Lock()
			Count--
			process.Unlock()
		}()
		time.Sleep(time.Millisecond * 100)
	}

	for{
		time.Sleep(time.Second * 60)
		logger.Debug("已进行到第",flag,"只股票，共",len(stockList),"只")
		if flag == len(stockList) && Count == 0{
			logger.Debug("执行完毕")
			break
		}
	}
}


func GetLargeFallStockInfoByCode(stockCode string,db *gorm.DB){
	var err error
	nowDate := fmt.Sprintf("%04d-%02d-%02d",time.Now().Year(),time.Now().Month(),time.Now().Day())

	stockNDaysPreInfo := strategy.GetStockInfoByDatePre(stockCode,nowDate,200,db)
	jysCode := util.GetJysCodeByStockCode(stockCode)
	var stockNowInfo *model.RealTimeStock
	if err,stockNowInfo = service.GetRealTimeStockInfoByStockCode(jysCode,stockCode) ; err != nil{
		logger.Error(err)
		return
	}
	if (stockNDaysPreInfo.HighPrice - stockNowInfo.NowPrice) / stockNowInfo.NowPrice > 1 && (stockNowInfo.NowPrice - stockNDaysPreInfo.LowPrice) / stockNDaysPreInfo.LowPrice < 0.3{
		var deepFallStock model.DeepFallStock
		if err = db.Where("code = ?", stockCode).First(&deepFallStock).Error ; err != nil{
			if err != gorm.ErrRecordNotFound{
				logger.Error(err)
				return
			}
		}

		deepFallStock.Code = stockCode
		deepFallStock.NowPrice = stockNowInfo.NowPrice
		deepFallStock.Date = stockNowInfo.DealDate
		deepFallStock.HighPrice = stockNDaysPreInfo.HighPrice
		deepFallStock.HighPriceDate = stockNDaysPreInfo.HighPriceDate
		deepFallStock.LowPrice = stockNDaysPreInfo.LowPrice
		deepFallStock.LowPriceDate = stockNDaysPreInfo.LowPriceDate

		if err = db.Save(&deepFallStock).Error; err != nil{
			logger.Error(err)
			return
		}
		//util.SendEmail("股票交易提示",
		//	"<div><h2>股票代码:" + stock.Code + "跌幅较大</h2></br>" +
		//		"<h4>较最高价下跌" + fmt.Sprint((stockNDaysPreInfo.HighPrice - stockNowInfo.NowPrice) / stockNDaysPreInfo.HighPrice) + "</h4></br>" +
		//		"<h4>较最低价上涨:" + fmt.Sprint((stockNowInfo.NowPrice - stockNDaysPreInfo.LowPrice) / stockNDaysPreInfo.LowPrice) + "</h4></br>" +
		//		"</div>")
	}else{
		var deepFallStock model.DeepFallStock
		if err = db.Where("code = ?",stockCode).First(&deepFallStock).Error ; err != nil{
			if err == gorm.ErrRecordNotFound{
				return
			}
			logger.Error(err)
			return
		}
		if err = db.Delete(&deepFallStock).Error; err != nil{
			logger.Error(err)
			return
		}
	}
}