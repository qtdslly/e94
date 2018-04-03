package strategy

import (
	"background/stock/model"
	"github.com/jinzhu/gorm"
	"background/common/logger"
	"strconv"
	"strings"
	"fmt"
)

/**
获取指定时间范围的股票前复权历史数据
 */
func GetStockHistoryDataQByDate(stockCode ,begin , end string ,db *gorm.DB)(error,[]*model.StockHistoryDataQ){
	var stockHistoryDataQs []*model.StockHistoryDataQ
	if err := db.Order("date asc").Where("code = ? and date >= ? and date <= ?" , stockCode,begin,end).Find(&stockHistoryDataQs).Error;err != nil{
		logger.Error(err)
		return err,stockHistoryDataQs
	}
	return nil,stockHistoryDataQs
}

/**
获取指定股票所有的前复权历史数据
 */
func GetStockHistoryDataQByCode(stockCode string,db *gorm.DB)(error,[]*model.StockHistoryDataQ){
	var stockHistoryDataQs []*model.StockHistoryDataQ
	if err := db.Order("date asc").Where("code = ?" , stockCode).Find(&stockHistoryDataQs).Error;err != nil{
		logger.Error(err)
		return err,stockHistoryDataQs
	}
	return nil,stockHistoryDataQs
}

func GetStockHistoryDataQByDatePre(stockCode,date string,count int ,db *gorm.DB)(error,[]*model.StockHistoryDataQ){
	var stockHistoryDataQs []*model.StockHistoryDataQ
	if err := db.Order("date desc").Where("code = ? and date < ?" , stockCode,date).Limit(count).Find(&stockHistoryDataQs).Error;err != nil{
		logger.Error(err)
		return err,stockHistoryDataQs
	}

	return nil,stockHistoryDataQs
}

func GetStockHistoryDataQByDateAfter(stockCode,date string,count int ,db *gorm.DB)(error,[]*model.StockHistoryDataQ){
	var stockHistoryDataQs []*model.StockHistoryDataQ
	if err := db.Order("date asc").Where("code = ? and date > ?" , stockCode,date).Limit(count).Find(&stockHistoryDataQs).Error;err != nil{
		logger.Error(err)
		return err,stockHistoryDataQs
	}
	return nil,stockHistoryDataQs
}

func GetStockHistoryTransDateByStockCode(stockCode string,db *gorm.DB)(error,[]string){
	var transDates []string
	if err := db.Order("date asc").Table("stock_history_data_q").Where("code = ?",stockCode).Pluck("date",&transDates).Error ; err != nil{
		logger.Error(err)
		return err,transDates
	}
	return nil,transDates
}

func GetHighAndLowHistoryStockInfoByStockInfos(stockHistoryDataQs []*model.StockHistoryDataQ)(highStockInfo,lowStockInfo model.StockHistoryDataQ){
	if len(stockHistoryDataQs) == 0{
		return highStockInfo,lowStockInfo
	}
	highStockInfo.Close = stockHistoryDataQs[0].Close
	lowStockInfo.Close  = stockHistoryDataQs[0].Close
	for _,stockHistoryDataQ := range stockHistoryDataQs{
		if(stockHistoryDataQ.Close > highStockInfo.Close){
			highStockInfo = *stockHistoryDataQ
		}
		if(stockHistoryDataQ.Close < lowStockInfo.Close){
			lowStockInfo = *stockHistoryDataQ
		}
	}
	return highStockInfo,lowStockInfo
}


/*
获取指定日期前N个交易区间内的股票最高价，最低价信息
*/
func GetStockInfoByDatePre(stockCode ,date string,days int,db *gorm.DB)(model.StockNDaysPreInfo){
	var err error
	var stockNDaysPreInfo model.StockNDaysPreInfo
	var beforeStockHistoryDataQs []*model.StockHistoryDataQ
	if err ,beforeStockHistoryDataQs = GetStockHistoryDataQByDatePre(stockCode,date,days,db) ; err != nil{
		logger.Error(err)
		return stockNDaysPreInfo
	}

	//获取最高点和最低点股票数据
	highStockHistoryDataQ,lowStockHistoryDataQ := GetHighAndLowHistoryStockInfoByStockInfos(beforeStockHistoryDataQs)
	stockNDaysPreInfo.HighPrice = highStockHistoryDataQ.Close
	stockNDaysPreInfo.HighPriceDate = highStockHistoryDataQ.Date
	stockNDaysPreInfo.LowPrice = lowStockHistoryDataQ.Close
	stockNDaysPreInfo.LowPriceDate = lowStockHistoryDataQ.Date
	return stockNDaysPreInfo

}

/*
获取涨幅
*/
func GetRose(first,second float64)(result float64){
	if second == 0.00{
		return 0.00;
	}
	return first  / second
}

func LowBuyHighSell(stockCode ,begin,end string,db *gorm.DB){
	var err error

	var allMoney,surplusMoney float64;
	allMoney = 100000
	surplusMoney = allMoney
	var isFirstBuy  = true
	var holdStockInfo model.HoldStockInfo
	if err = db.Where("code = ?",stockCode).First(&holdStockInfo).Error ; err != nil{
		if err == gorm.ErrRecordNotFound{
			holdStockInfo.HoldDays = 0
			holdStockInfo.HoldMoney = 0.00
			holdStockInfo.Cost = 0.00
			holdStockInfo.Fee = 0.00
			holdStockInfo.AllCount = 0
			holdStockInfo.AvgPrice = 0.00
			holdStockInfo.Date = ""
			holdStockInfo.Price = 0.00
			holdStockInfo.Code = stockCode
			holdStockInfo.Profit = 0.00
			isFirstBuy  = false
		}
	}


	var transDates []string
	if err ,transDates = GetStockHistoryTransDateByStockCode(stockCode,db) ; err != nil{
		logger.Error(err)
		return
	}

	if len(transDates) < 60{
		logger.Debug("上市时间短，采集数据不够")
		return
	}

	var beforeStockHistoryDataQs []*model.StockHistoryDataQ
	if err ,beforeStockHistoryDataQs = GetStockHistoryDataQByDatePre(stockCode,begin,180,db) ; err != nil{
		logger.Error(err)
		return
	}

	//获取最高点和最低点股票数据
	highStockHistoryDataQ,lowStockHistoryDataQ := GetHighAndLowHistoryStockInfoByStockInfos(beforeStockHistoryDataQs)

	if(len(beforeStockHistoryDataQs) < 60){
		var tmpStockHistoryDataQs []*model.StockHistoryDataQ
		if err ,tmpStockHistoryDataQs = GetStockHistoryDataQByDateAfter(stockCode,transDates[0],60,db) ; err != nil{
			logger.Error(err)
			return
		}

		highStockHistoryDataQ,lowStockHistoryDataQ = GetHighAndLowHistoryStockInfoByStockInfos(tmpStockHistoryDataQs)
	}

	if highStockHistoryDataQ.Close == lowStockHistoryDataQ.Close || highStockHistoryDataQ.Close == 0.00 || lowStockHistoryDataQ.Close == 0.00 {
		return
	}

	//如果最高价与最低价相差幅度小于50%,直接退出
	var rose float64
	rose = GetRose(highStockHistoryDataQ.Close - lowStockHistoryDataQ.Close,highStockHistoryDataQ.Close)
	if rose < 0.2{
		logger.Debug("最高价与最低价相差幅度太小，不适合做高抛低吸操作")
		return
	}

	var stockHistoryDataQs []*model.StockHistoryDataQ
	if err ,stockHistoryDataQs = GetStockHistoryDataQByDate(stockCode,begin,end,db) ; err != nil{
		logger.Error(err)
		return
	}

	var isHoldStock = false
	var endDate string
	for _,stockHistoryDataQ := range stockHistoryDataQs{
		if stockHistoryDataQ.Close == 0.00{
			if isHoldStock{
				holdStockInfo.HoldDays++;
			}
			continue
		}
		var transParamInfo model.TransParamInfo
		transParamInfo.Price = stockHistoryDataQ.Close
		transParamInfo.Date = stockHistoryDataQ.Date
		transParamInfo.HighPrice = highStockHistoryDataQ.Close
		transParamInfo.LowPrice = lowStockHistoryDataQ.Close

		if BuyStock(&holdStockInfo,transParamInfo,&isFirstBuy,allMoney,surplusMoney,db) != nil{
			return
		}
		if SellStock(&holdStockInfo,transParamInfo,allMoney,surplusMoney,db) != nil{
			return
		}

		if holdStockInfo.AllCount != 0{
			isHoldStock = true
		}else{
			isHoldStock = false
		}
		if isHoldStock{
			holdStockInfo.HoldDays++;
		}

		endDate = stockHistoryDataQ.Date

		if(stockHistoryDataQ.Close > highStockHistoryDataQ.Close){
			highStockHistoryDataQ = *stockHistoryDataQ
		}
		if(stockHistoryDataQ.Close < lowStockHistoryDataQ.Close){
			lowStockHistoryDataQ = *stockHistoryDataQ
		}
	}

	if !isFirstBuy || holdStockInfo.AllCount == 0{
		return
	}
	if err = GetNowProfitInfo(&holdStockInfo,endDate,db) ; err != nil{
		logger.Error(err)
		return
	}
	holdStockInfo.Date = endDate
	holdStockInfo.TransType = "total"

	PrintProfit(holdStockInfo)
}

func GetRealTimeStockTransDateByStockCode(stockCode string,db *gorm.DB)(error,[]string){
	var transDates []string
	if err := db.Order("deal_date asc").Table("real_time_stock").Where("stock_code = ?",stockCode).Pluck("deal_date",&transDates).Error ; err != nil{
		logger.Error(err)
		return err,transDates
	}
	return nil,transDates
}


func GetNowProfitInfo(holdStockInfo *model.HoldStockInfo,date string,db *gorm.DB)(error){
	var transDates []string
	var err error
	if err ,transDates = GetRealTimeStockTransDateByStockCode(holdStockInfo.Code,db) ; err != nil{
		logger.Error(err)
		return err
	}
	var realTimeStock model.RealTimeStock
	if err := db.Where("stock_code = ? and deal_date = ?",holdStockInfo.Code,transDates[len(transDates) - 1]).First(&realTimeStock).Error ; err != nil{
		logger.Error(err)
		return err
	}
	if realTimeStock.NowPrice == 0.00{
		return nil
	}
	holdStockInfo.Date = date
	holdStockInfo.Price = realTimeStock.NowPrice
	holdStockInfo.HoldMoney = holdStockInfo.Price * float64(holdStockInfo.AllCount)
	holdStockInfo.FloatProfit = holdStockInfo.HoldMoney - holdStockInfo.Cost - holdStockInfo.Fee
	if holdStockInfo.HoldMoney > 0.00{
		holdStockInfo.FloatProfitRate = holdStockInfo.FloatProfit / holdStockInfo.HoldMoney * 100.00
	}else{
		holdStockInfo.FloatProfitRate = 0.00
	}
	if holdStockInfo.AllCount == 0{
		holdStockInfo.AvgPrice = 0.0
	}else{
		holdStockInfo.AvgPrice = (holdStockInfo.HoldMoney - holdStockInfo.FloatProfit) / float64(holdStockInfo.AllCount)
	}

	if err = db.Save(&holdStockInfo).Error ; err != nil{
		logger.Error(err)
		return err
	}
	return nil
}

func GetTransDaysByDatePre(stockCode,date string,db *gorm.DB)(int){
	var count int
	if err := db.Table("stock_history_data_q").Where("code = ? and date < ?",stockCode,date).Count(&count).Error ; err != nil{
		logger.Error(err)
		return 0
	}
	return count
}

func FloatNumToStockCount(num float64)(result int64){
	result,_ = strconv.ParseInt(strings.Split(fmt.Sprint(num),".")[0],10,64 )
	return result
}


func BuyStock(holdStockInfo *model.HoldStockInfo,stockParamInfo model.TransParamInfo,isFirstBuy *bool,allMoney,surplusMoney float64,db *gorm.DB)(error){
	var rate float64
	if stockParamInfo.Price < stockParamInfo.LowPrice{
		rate = GetRose(stockParamInfo.HighPrice - stockParamInfo.LowPrice,stockParamInfo.HighPrice - stockParamInfo.Price)
		if rate < 0.07{
			rate = 0.07
		}
	}else{
		rate = GetRose(stockParamInfo.HighPrice - stockParamInfo.Price , stockParamInfo.HighPrice - stockParamInfo.LowPrice)
	}
	var money float64
	if !(*isFirstBuy){
		if rate > 5{
			money = allMoney * 0.5
		}else if rate > 4{
			money = allMoney * 0.5 * 0.9
		}else if rate > 3{
			money = allMoney * 0.5 * 0.8
		}else if rate > 2{
			money = allMoney * 0.5 * 0.8
		}else if rate > 1{
			money = allMoney * 0.5 * 0.7
		}else{
			money = allMoney * 0.5 * 0.6
		}
	}else if(GetRose(holdStockInfo.Price - stockParamInfo.Price,holdStockInfo.Price) > 0.15){
		if rate > 5{
			money = allMoney * 0.2
		}else if rate > 4{
			money = allMoney * 0.15
		}else if rate > 3{
			money = allMoney * 0.12
		}else if rate > 2{
			money = allMoney * 0.1
		}else if rate > 1{
			money = allMoney * 0.08
		}else{
			money = allMoney * 0.07
		}
		if money > surplusMoney{
			money = surplusMoney * 0.8
		}

	}else{
		return nil
	}

	holdStockInfo.Count = GetTransCount(money,stockParamInfo.Price)

	if holdStockInfo.FloatProfitRate < -0.2{
		holdStockInfo.Count = FloatNumToStockCount(float64(holdStockInfo.AllCount) * 0.3 / 100.00) * 100
	}

	if holdStockInfo.Count == 0{
		holdStockInfo.Count = 100
	}

	holdStockInfo.TransType = "buy"
	holdStockInfo.Date = stockParamInfo.Date
	holdStockInfo.Price = stockParamInfo.Price

	if err := CacluTransInfoByTransCount(holdStockInfo,db) ; err != nil{
		logger.Error(err)
		return err
	}

	PrintProfit(*holdStockInfo)

	if !(*isFirstBuy){
		if err := db.Create(&holdStockInfo).Error ; err != nil {
			logger.Error(err)
			return err
		}
	}else{
		if err := db.Save(&holdStockInfo).Error ; err != nil {
			logger.Error(err)
			return err
		}
	}

	*isFirstBuy = true
	return nil
}


func SellStock(holdStockInfo *model.HoldStockInfo,stockParamInfo model.TransParamInfo,allMoney,surplusMoney float64,db *gorm.DB)(error) {
	if holdStockInfo.AllCount == 0{
		return nil
	}
	if GetRose(stockParamInfo.Price - holdStockInfo.Price,holdStockInfo.Price) < 0.1 || stockParamInfo.Price < holdStockInfo.Price{
		return nil
	}else{
		var transDatePreCount int
		transDatePreCount = GetTransDaysByDatePre(holdStockInfo.Code,stockParamInfo.Date,db)
		var stock60DaysPreInfo model.StockNDaysPreInfo
		if transDatePreCount > 60{
			stock60DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,stockParamInfo.Date,60,db)
		}

		var count int64
		var rose float64

		if stockParamInfo.Price > stock60DaysPreInfo.HighPrice{
			count = holdStockInfo.AllCount / 2
		}else if (stockParamInfo.Price > stock60DaysPreInfo.LowPrice){
			rose = GetRose(stockParamInfo.Price - stock60DaysPreInfo.LowPrice,stock60DaysPreInfo.HighPrice - stock60DaysPreInfo.LowPrice)
			count,_ = strconv.ParseInt(strings.Split(fmt.Sprint(float64(holdStockInfo.AllCount) * rose),".")[0],10,64)
		}else{
			return nil
		}
		//var stock3DaysPreInfo,stock5DaysPreInfo,stock7DaysPreInfo,stock15DaysPreInfo,stock30DaysPreInfo model.StockNDaysPreInfo
		//var stock60DaysPreInfo,stock120DaysPreInfo,stock180DaysPreInfo,stock300DaysPreInfo model.StockNDaysPreInfo
		//var transDatePreCount int
		//transDatePreCount = GetTransDaysByDatePre(holdStockInfo.Code,date,db)
		//if transDatePreCount > 300{
		//	stock300DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,300,db)
		//}
		//if transDatePreCount > 180{
		//	stock180DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,180,db)
		//}
		//if transDatePreCount > 120{
		//	stock120DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,120,db)
		//}
		//if transDatePreCount > 60{
		//	stock60DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,60,db)
		//}
		//if transDatePreCount > 30{
		//	stock30DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,30,db)
		//}
		//if transDatePreCount > 15{
		//	stock15DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,15,db)
		//}
		//if transDatePreCount > 7{
		//	stock7DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,7,db)
		//}
		//if transDatePreCount > 5{
		//	stock5DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,5,db)
		//}
		//if transDatePreCount > 3{
		//	stock3DaysPreInfo = GetStockInfoByDatePre(holdStockInfo.Code,date,3,db)
		//}

		count = int64(count / 100) * 100
		if count >= holdStockInfo.Count{
			count = holdStockInfo.Count - 100
		}

		if count < 100{
			count = 100
		}
		holdStockInfo.Count = count
		if holdStockInfo.AllCount < count{
			count = holdStockInfo.AllCount
		}
	}

	holdStockInfo.TransType = "sell"
	holdStockInfo.Date = stockParamInfo.Date
	holdStockInfo.Price = stockParamInfo.Price

	if err := CacluTransInfoByTransCount(holdStockInfo,db) ; err != nil{
		logger.Error(err)
		return err
	}

	PrintProfit(*holdStockInfo)

	if err := db.Save(&holdStockInfo).Error ; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func CacluTransInfoByTransCount(holdStockInfo *model.HoldStockInfo,db *gorm.DB)(error){
	var fee float64
	if holdStockInfo.Price * float64(holdStockInfo.Count) * 0.0002 < 5{
		fee = 5
	}else{
		fee = holdStockInfo.Price * float64(holdStockInfo.Count) * 0.0002
	}

	fee += holdStockInfo.Price * float64(holdStockInfo.Count) * 0.00002
	if holdStockInfo.TransType == "sell"{
		fee += holdStockInfo.Price * float64(holdStockInfo.Count) * 0.001
		holdStockInfo.AllCount -= holdStockInfo.Count
		holdStockInfo.Cost -= holdStockInfo.Price * float64(holdStockInfo.Count)
	}else{
		holdStockInfo.AllCount += holdStockInfo.Count
		holdStockInfo.Cost += holdStockInfo.Price * float64(holdStockInfo.Count)
	}

	holdStockInfo.Fee += fee
	holdStockInfo.HoldMoney = holdStockInfo.Price * float64(holdStockInfo.AllCount)

	//logger.Debug("holdStockInfo.Price:",holdStockInfo.Price," holdStockInfo.AvgPrice:",holdStockInfo.AvgPrice)
	holdStockInfo.Profit += holdStockInfo.HoldMoney - holdStockInfo.Cost - holdStockInfo.Fee
	holdStockInfo.FloatProfit = holdStockInfo.HoldMoney - holdStockInfo.Cost - holdStockInfo.Fee
	if holdStockInfo.HoldMoney > 0.00{
		holdStockInfo.FloatProfitRate = holdStockInfo.FloatProfit / holdStockInfo.HoldMoney * 100.00
	}else{
		holdStockInfo.FloatProfitRate = 0.00
	}

	var stockTransInfo model.TransStockInfo
	stockTransInfo.Code = holdStockInfo.Code
	stockTransInfo.Count = holdStockInfo.Count
	stockTransInfo.Date = holdStockInfo.Date
	stockTransInfo.Fee = fee
	stockTransInfo.Price = holdStockInfo.Price
	stockTransInfo.TransType = holdStockInfo.TransType
	stockTransInfo.Cost = holdStockInfo.Price * float64(holdStockInfo.Count)
	stockTransInfo.Profit = (holdStockInfo.Price - holdStockInfo.AvgPrice) * float64(holdStockInfo.Count)
	if err := db.Save(&stockTransInfo).Error ; err != nil{
		logger.Error(err)
		return err
	}

	if holdStockInfo.AllCount == 0{
		holdStockInfo.AvgPrice = 0.0
	}else{
		holdStockInfo.AvgPrice = (holdStockInfo.HoldMoney - holdStockInfo.FloatProfit) / float64(holdStockInfo.AllCount)
	}

	return nil
}

func GetTransCount(money,price float64)(result int64){
	result,_ = strconv.ParseInt(strings.Split(fmt.Sprint( money / price / 100 ),".")[0],10,64)
	return result * 100
}

func PrintProfit(holdStockInfo model.HoldStockInfo){
	var dealCount int64 = 0
	if holdStockInfo.TransType != "total"{
		dealCount = holdStockInfo.Count
	}

	logger.Debug("交易类型:",holdStockInfo.TransType," 当前日期:",holdStockInfo.Date," 交易量:",
		dealCount,"股 当前价格:",holdStockInfo.Price,
		" 当前股票持仓:",holdStockInfo.AllCount,"股 持股市值:",holdStockInfo.HoldMoney," 总成本:",holdStockInfo.Cost,
		" 持股天数:",holdStockInfo.HoldDays," 平均成本:",holdStockInfo.AvgPrice,
		" 浮动盈亏:",holdStockInfo.FloatProfit," 盈亏率:",holdStockInfo.FloatProfitRate ,"%")

	//logger.Debug("交易类型:",holdStockInfo.TransType," 当前日期:",holdStockInfo.Date," 交易量:",dealCount,"股 当前价格:",holdStockInfo.Price,
	//	" 当前股票持仓:",holdStockInfo.AllCount,"股 总成本:",holdStockInfo.Cost," 持股市值:",holdStockInfo.HoldMoney,
	//	" 持股天数:",holdStockInfo.HoldDays," 总手续费:",holdStockInfo.Free," 盈亏:",holdStockInfo.Profit,
	//	" 平均成本:",holdStockInfo.AvgPrice," 浮动盈亏:",holdStockInfo.FloatProfit)
}