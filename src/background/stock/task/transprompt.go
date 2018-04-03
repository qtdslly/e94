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

		if(realTimeStockInfo.NowPrice <= transPrompt.PromptBuyPrice){
			logger.Print("============================================================")
			util.SendEmail("股票交易提示",
				"<div style='color:#F00;background:url(https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1521732917614&di=0f4226d53420474a7630a4b49b2c26f4&imgtype=0&src=http%3A%2F%2Fb.zol-img.com.cn%2Fsjbizhi%2Fimages%2F2%2F320x510%2F1354095141925.jpg) no-repeat 0px 0px;'><h2>股票代码:" + transPrompt.StockCode + "到达买入价格</h2></br>" +
					"<h4>设定买入价格为:" + fmt.Sprint(transPrompt.PromptBuyPrice) + "</h4></br>" +
					"<h4>当前价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>设定的交易量为:" + fmt.Sprint(transPrompt.PromptBuyCount) + "</h4></br>" +
					"<h1>请尽快交易!!!</h1></div>")
			havePrompt = true
		}else if(realTimeStockInfo.NowPrice >= transPrompt.PromptSellPrice){
			logger.Print("============================================================")

			util.SendEmail("股票交易提示",
				"<div style='color:#F00;background:url(https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1521732917614&di=0f4226d53420474a7630a4b49b2c26f4&imgtype=0&src=http%3A%2F%2Fb.zol-img.com.cn%2Fsjbizhi%2Fimages%2F2%2F320x510%2F1354095141925.jpg) no-repeat 0px 0px;'><h2>股票代码:" + transPrompt.StockCode + "到达卖出价格</h2></br>" +
					"<h4>设定卖出价格为:" + fmt.Sprint(transPrompt.PromptSellPrice) + "</h4></br>" +
					"<h4>当前价格为:" + fmt.Sprint(realTimeStockInfo.NowPrice) + "</h4></br>" +
					"<h4>设定的交易量为:" + fmt.Sprint(transPrompt.PromptSellCount) + "</h4></br>" +
					"<h1>请尽快交易!!!</h1></div>")
			havePrompt = true
		}else{
			//logger.Debug("未到交易价格，暂不交易!!!!")
		}
		time.Sleep(time.Second)
	}

}
