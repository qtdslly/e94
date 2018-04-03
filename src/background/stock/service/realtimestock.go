package service

import (
	"background/stock/model"
	"strings"
	"background/common/logger"
	"errors"
	"strconv"
	"net/http"
	"io/ioutil"
)

const SINA_REALTIME_STOCK_URL = "http://hq.sinajs.cn/"
func GetRealTimeStockInfoByStockCode(jysCode,stockCode string)(error,*model.RealTimeStock){
	var realTimeStock *model.RealTimeStock
	url := SINA_REALTIME_STOCK_URL + "list=" + jysCode + stockCode
	resp, _ := http.Get(url)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		logger.Error("获取股票实时信息失败:",jysCode,stockCode,"err:",err)
		return err,realTimeStock
	}
	if len(string(data)) < 30{
		return errors.New("股票代码不存在"),realTimeStock
	}

	if err ,realTimeStock = GetRealTimeStockObject(jysCode,stockCode,string(data)) ; err != nil{
		logger.Error("解析股票信息失败:",jysCode,stockCode,"err:",err)
		return err,realTimeStock
	}
	return nil ,realTimeStock
}

func GetRealTimeStockObject(jysCode,stockCode,strObj string)(error ,*model.RealTimeStock){
	var realTimeStock model.RealTimeStock
	strObj = strObj[21:len(strObj)-3]
	realTimeStockInfo := strings.Split(strObj,",")
	if len(realTimeStockInfo) == 0 {
		logger.Error("分割字段，获取股票信息出错")
		return errors.New("分割字段，获取股票信息出错"), nil
	}
	if len(realTimeStockInfo) != 33{
		logger.Error("股票信息不完整" )
		return errors.New("股票信息不完整") , nil
	}
	realTimeStock.JysCode = jysCode
	realTimeStock.StockCode = stockCode
	realTimeStock.TodayOpenPrice, _ = strconv.ParseFloat(realTimeStockInfo[1], 64) // 今日开盘价    1
	realTimeStock.YestdayClosePrice, _ = strconv.ParseFloat(realTimeStockInfo[2], 64) // 昨日收盘价    2
	realTimeStock.NowPrice, _ = strconv.ParseFloat(realTimeStockInfo[3], 64) // 当前价格    3
	realTimeStock.TodayHighPrice, _ = strconv.ParseFloat(realTimeStockInfo[4], 64) // 今日最高价   4
	realTimeStock.TodayLowPrice, _ = strconv.ParseFloat(realTimeStockInfo[5], 64) // 今日最低价   5
	realTimeStock.BuyPrice, _ = strconv.ParseFloat(realTimeStockInfo[6], 64) // 竞买价(买一) 6
	realTimeStock.SellPrice, _ = strconv.ParseFloat(realTimeStockInfo[7], 64) // 竞卖价(卖一) 7
	realTimeStock.DealCount,_ = strconv.Atoi(realTimeStockInfo[8]) // 成交数(单位为1股)       8
	realTimeStock.DealMoney, _ = strconv.ParseFloat(realTimeStockInfo[9], 64) // 成交金额(单位为元)      9
	realTimeStock.BuyCount1,_ = strconv.Atoi(realTimeStockInfo[10]) // 买一申请股数(单位为1股) 10
	realTimeStock.BuyPrice1,_ = strconv.ParseFloat(realTimeStockInfo[11], 64) // 买一报价 11
	realTimeStock.BuyCount2,_ = strconv.Atoi(realTimeStockInfo[12]) // 买二申请股数(单位为1股) 12
	realTimeStock.BuyPrice2,_ = strconv.ParseFloat(realTimeStockInfo[13], 64) // 买二报价 13
	realTimeStock.BuyCount3,_ = strconv.Atoi(realTimeStockInfo[14]) // 买三申请股数(单位为1股) 14
	realTimeStock.BuyPrice3,_ = strconv.ParseFloat(realTimeStockInfo[15], 64) // 买三报价 15
	realTimeStock.BuyCount4,_ = strconv.Atoi(realTimeStockInfo[16]) // 买四申请股数(单位为1股) 16
	realTimeStock.BuyPrice4,_ = strconv.ParseFloat(realTimeStockInfo[17], 64) // 买四报价 17
	realTimeStock.BuyCount5,_ = strconv.Atoi(realTimeStockInfo[18]) // 买五申请股数(单位为1股) 18
	realTimeStock.BuyPrice5,_ = strconv.ParseFloat(realTimeStockInfo[19], 64) // 买五报价 19
	realTimeStock.SellCount1,_ = strconv.Atoi(realTimeStockInfo[20]) // 卖一申请股数(单位为1股) 20
	realTimeStock.SellPrice1,_ = strconv.ParseFloat(realTimeStockInfo[21], 64) // 卖一报价 21
	realTimeStock.SellCount2,_ = strconv.Atoi(realTimeStockInfo[22]) // 卖二申请股数(单位为1股) 22
	realTimeStock.SellPrice2,_ = strconv.ParseFloat(realTimeStockInfo[23], 64) // 卖二报价 23
	realTimeStock.SellCount3,_ = strconv.Atoi(realTimeStockInfo[24]) // 卖三申请股数(单位为1股) 24
	realTimeStock.SellPrice3,_ = strconv.ParseFloat(realTimeStockInfo[25], 64) // 卖三报价 25
	realTimeStock.SellCount4,_ = strconv.Atoi(realTimeStockInfo[26]) // 卖四申请股数(单位为1股) 26
	realTimeStock.SellPrice4,_ = strconv.ParseFloat(realTimeStockInfo[27], 64) // 卖四报价 27
	realTimeStock.SellCount5,_ = strconv.Atoi(realTimeStockInfo[28]) // 卖五申请股数(单位为1股) 28
	realTimeStock.SellPrice5,_ = strconv.ParseFloat(realTimeStockInfo[29], 64) // 卖五报价 29

	realTimeStock.DealDate = realTimeStockInfo[30] // 日期 30
	realTimeStock.DealTime = realTimeStockInfo[31] // 时间 31

	return nil,&realTimeStock
}
