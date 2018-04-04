package task

import (
	"background/stock/model"
	"background/stock/service"
	"background/stock/tools/util"
	"background/common/logger"

	"time"
	"fmt"

	"github.com/jinzhu/gorm"
)
/*
股票交易提示
*/

func TransPromptAll(db *gorm.DB){
	var transPrompts []model.TransPrompt
	if err := db.Find(&transPrompts).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _ , transPrompt := range transPrompts{
		go TransPromptByPromptInfo(transPrompt)
		go RosePrompt(transPrompt,db)
	}

	for{
		time.Sleep(time.Minute)
	}
}

func TransPromptByPromptInfo(transPrompt model.TransPrompt){
	var err error
	var jysCode string
	jysCode = util.GetJysCodeByStockCode(transPrompt.StockCode)
	if jysCode == ""{
		logger.Error("获取交易所代码错误，程序退出!!!")
		return
	}
	var realTimeStockInfo *model.RealTimeStock
	var havePrompt bool = false
	for{
		if havePrompt{
			break
		}
		if err,realTimeStockInfo = service.GetRealTimeStockInfoByStockCode(jysCode,transPrompt.StockCode) ; err != nil{
			logger.Error("获取股票实时信息失败，程序退出!!!")
			return
		}
		if realTimeStockInfo.NowPrice == 0.00{
			break
		}

		var result bool
		if(realTimeStockInfo.NowPrice <= transPrompt.PromptBuyPrice){
			result = util.SendEmail("股票交易提示",
				"<div><h2>股票代码:" + transPrompt.StockCode + "到达买入价格</h2></br>" +
					"<h4>设定买入价格为:" + fmt.Sprint(transPrompt.PromptBuyPrice) + "</h4></br>" +
					"<h4>当前价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>设定的交易量为:" + fmt.Sprint(transPrompt.PromptBuyCount) + "</h4></br>" +
					"<h1>请尽快交易!!!</h1></div>")
		}else if(realTimeStockInfo.NowPrice >= transPrompt.PromptSellPrice){
			result = util.SendEmail("股票交易提示",
				"<div'><h2>股票代码:" + transPrompt.StockCode + "到达卖出价格</h2></br>" +
					"<h4>设定卖出价格为:" + fmt.Sprint(transPrompt.PromptSellPrice) + "</h4></br>" +
					"<h4>当前价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>设定的交易量为:" + fmt.Sprint(transPrompt.PromptSellCount) + "</h4></br>" +
					"<h1>请尽快交易!!!</h1></div>")
		}else{
			//logger.Debug("未到交易价格，暂不交易!!!!")
		}
		if result{
			logger.Debug("邮件发送成功!!!")
			havePrompt = true
		}
		time.Sleep(time.Second)
	}

}



func RosePrompt(transPrompt model.TransPrompt,db *gorm.DB){
	var err error
	var jysCode string
	jysCode = util.GetJysCodeByStockCode(transPrompt.StockCode)
	if jysCode == ""{
		logger.Error("获取交易所代码错误，程序退出!!!")
		return
	}
	var realTimeStockInfo *model.RealTimeStock
	var havePrompt bool = false
	for{
		if havePrompt{
			time.Sleep(time.Minute * 30)
		}
		if err,realTimeStockInfo = service.GetRealTimeStockInfoByStockCode(jysCode,transPrompt.StockCode) ; err != nil{
			logger.Error("获取股票实时信息失败，程序退出!!!")
			return
		}
		if realTimeStockInfo.NowPrice == 0.00{
			break
		}

		if ( realTimeStockInfo.NowPrice - realTimeStockInfo.YestdayClosePrice ) / realTimeStockInfo.YestdayClosePrice > 0.03 ||
			( realTimeStockInfo.NowPrice - realTimeStockInfo.YestdayClosePrice ) / realTimeStockInfo.YestdayClosePrice < -0.03{
			util.SendEmail("股票涨跌幅提示",
				"<div'><h2>股票代码:" + transPrompt.StockCode + "</h2></br>" +
					"<h2>股票名称:" + util.GetNameByCode(transPrompt.StockCode,db) + "</h2></br>" +
					"<h4>当前股票价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"</div>")
			havePrompt = true
		}

		time.Sleep(time.Second)
	}

}