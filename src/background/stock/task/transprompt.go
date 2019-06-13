package task

import (
	"background/stock/model"
	"background/stock/service"
	cc "background/stock/common"
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
	if err := db.Where("status = 1").Find(&transPrompts).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _ , transPrompt := range transPrompts{
		//go TransPromptByPromptInfo(transPrompt,db)
		go RosePrompt(transPrompt,db)
	}

	for{
		time.Sleep(time.Minute)
		var p = time.Now()
		if fmt.Sprintf("%02d%02d",p.Hour(),p.Minute()) > "1500"{
			break
		}
	}
}

func TransPromptByPromptInfo(transPrompt model.TransPrompt,db *gorm.DB){
	var err error
	var jysCode string
	jysCode = cc.GetJysCodeByStockCode(transPrompt.StockCode)
	if jysCode == ""{
		logger.Error("获取交易所代码错误，程序退出!!!")
		return
	}
	var realTimeStockInfo *model.RealTimeStock
	var havePrompt bool = false
	for{
		if havePrompt{
			time.Sleep(time.Minute * 5)
		}
		if err,realTimeStockInfo = service.GetRealTimeStockInfoByStockCode(jysCode,transPrompt.StockCode) ; err != nil{
			logger.Error("获取股票实时信息失败，程序退出!!!")
			return
		}
		if realTimeStockInfo.NowPrice == 0.00{
			break
		}

		//rose := (realTimeStockInfo.NowPrice - realTimeStockInfo.YestdayClosePrice) / realTimeStockInfo.YestdayClosePrice
		var result bool
		if(realTimeStockInfo.NowPrice <= transPrompt.PromptBuyPrice){
			result = cc.SendEmail("股票交易提示",
				"<div><h2>股票代码:" + transPrompt.StockCode + "  股票名称:" + cc.GetNameByCode(transPrompt.StockCode,db) + "到达买入价格</h2></br>" +
					"<h4>设定买入价格为:" + fmt.Sprint(transPrompt.PromptBuyPrice) + "</h4></br>" +
					"<h4>当前价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>设定的交易量为:" + fmt.Sprint(transPrompt.PromptBuyCount) + "</h4></br>" +
					"<h4>当前涨幅为:" + fmt.Sprint((((realTimeStockInfo.NowPrice - realTimeStockInfo.YestdayClosePrice) / realTimeStockInfo.YestdayClosePrice) * 100.00)) + "%" + "</h4></br:>" +
					"<h1>请尽快交易!!!</h1></div>")
		}else if(realTimeStockInfo.NowPrice >= transPrompt.PromptSellPrice){
			result = cc.SendEmail("股票交易提示",
				"<div'><h2>股票代码:" + transPrompt.StockCode + "  股票名称:" + cc.GetNameByCode(transPrompt.StockCode,db) + "到达卖出价格</h2></br>" +
					"<h4>设定卖出价格为:" + fmt.Sprint(transPrompt.PromptSellPrice) + "</h4></br>" +
					"<h4>当前价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>设定的交易量为:" + fmt.Sprint(transPrompt.PromptSellCount) + "</h4></br>" +
					"<h4>当前涨幅为:" + fmt.Sprint((((realTimeStockInfo.NowPrice - realTimeStockInfo.YestdayClosePrice) / realTimeStockInfo.YestdayClosePrice) * 100.00)) + "%" + "</h4></br:>" +
					"<h1>请尽快交易!!!</h1></div>")
		}//else if( rose > 0.02){
		//	result = util.SendEmail("股票交易提示",
		//		"<div'><h2>股票代码:" + transPrompt.StockCode + "  股票名称:" + util.GetNameByCode(transPrompt.StockCode,db) + "</h2></br>" +
		//			"<h4>当前涨幅:" + fmt.Sprintf("%.2f",rose * 100) + "</h4></br></div>" )
		//}else if (rose < -0.02){
		//	result = util.SendEmail("股票交易提示",
		//		"<div'><h2>股票代码:" + transPrompt.StockCode + "  股票名称:" + util.GetNameByCode(transPrompt.StockCode,db) + "</h2></br>" +
		//		"<h4>当前跌:" + fmt.Sprintf("%.2f",rose * 100) + "</h4></br></div>")
		//}
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
	jysCode = cc.GetJysCodeByStockCode(transPrompt.StockCode)
	if jysCode == ""{
		logger.Error("获取交易所代码错误，程序退出!!!")
		return
	}
	var realTimeStockInfo *model.RealTimeStock
	var havePrompt bool = false
	for{
		if havePrompt{
			time.Sleep(time.Minute * 10)
			havePrompt = false
		}
		if err,realTimeStockInfo = service.GetRealTimeStockInfoByStockCode(jysCode,transPrompt.StockCode) ; err != nil{
			logger.Error("获取股票实时信息失败，程序退出!!!")
			return
		}
		if realTimeStockInfo.NowPrice == 0.00{
			break
		}
		rose := (realTimeStockInfo.NowPrice - realTimeStockInfo.YestdayClosePrice) / realTimeStockInfo.YestdayClosePrice
		if rose > 0.03 || rose < -0.02{
			cc.SendEmail("股票涨跌幅提示",
				"<div'><h2>股票代码:" + transPrompt.StockCode + "</h2></br>" +
					"<h2>股票名称:" + cc.GetNameByCode(transPrompt.StockCode,db) + "</h2></br>" +
					"<h4>当前股票价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>当前股票涨幅为:" + fmt.Sprintf("%.2f",rose * 100) + "</h4></br>" +
					"</div>")
			havePrompt = true
		}

		time.Sleep(time.Second)
	}

}
