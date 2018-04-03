package task

import (
	"background/stock/model"
	"background/common/logger"
	"fmt"
	"background/stock/service"
)

func SyncHoldStockRealTimeInfo(){
	shHoldStockList := []string{"600050","600570","600839","601258","603368"}
	szHoldStockList := []string{"000897","002175","002269","002339","002395","002522","300033","300168","300315"}

	for _,stockCode := range shHoldStockList{
		var realTimeStock *model.RealTimeStock
		var err error
		if err,realTimeStock = service.GetRealTimeStockInfoByStockCode("sh", stockCode) ; err != nil{
			if err.Error() == "股票代码不存在"{
				continue
			}else{
				return
			}
		}
		logger.Printf("当前价:" + fmt.Sprint(realTimeStock.NowPrice))
	}
	for _,stockCode := range szHoldStockList{
		var realTimeStock *model.RealTimeStock
		var err error
		if err,realTimeStock = service.GetRealTimeStockInfoByStockCode("sz", stockCode) ; err != nil{
			if err.Error() == "股票代码不存在"{
				continue
			}else{
				return
			}
		}
		logger.Printf("当前价:" + fmt.Sprint(realTimeStock.NowPrice))
	}
}
