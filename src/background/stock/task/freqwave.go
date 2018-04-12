package task

import (
	"github.com/jinzhu/gorm"
	"background/stock/model"
	"background/common/logger"
	"time"
	"fmt"
	"background/stock/strategy"
)

/*
寻找波动频繁的股票
*/

func GetFreqWaveStock(db *gorm.DB){
	var err error
	var stocks []model.StockList
	if err = db.Find(&stocks).Error ; err != nil{
		logger.Error(err)
		return
	}
	p := time.Now()
	for _ , stock := range stocks{
		if stock.TimeToMarket < fmt.Sprintf("%04d%02d%02d",p.Year() - 2 ,p.Month(),p.Day()){
			continue
		}

		start := fmt.Sprintf("%04d-%02d-%02d",p.Year() - 2 ,p.Month(),p.Day())
		end := fmt.Sprintf("%04d-%02d-%02d",p.Year() ,p.Month(),p.Day())
		var stockHistoryDataQNews []*model.StockHistoryDataQ
		if err , stockHistoryDataQNews = strategy.GetStockHistoryDataQByDate(stock.Code,start,end,db) ; err != nil{
			logger.Error(err)
			return
		}

		low := stockHistoryDataQNews[0].Close
		high := low
		count := 0
		flag := false
		for _,stockHistoryDataQNew := range stockHistoryDataQNews{
			if (stockHistoryDataQNew.Close - low) / low > 0.15 && !flag{
				count++
				flag = true
			}
			if (high - stockHistoryDataQNew.Close) / high > 0.15 && flag{
				count++
				flag = false
			}

			if stockHistoryDataQNew.Close > high {
				high = stockHistoryDataQNew.Close

				if flag{
					low = high
				}
			}
			if stockHistoryDataQNew.Close < low{
				low = stockHistoryDataQNew.Close
				if !flag{
					high = low
				}
			}

		}

	}
}