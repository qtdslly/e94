package model

//type HoldStockInfo struct{
//	Code 	               string     //股票代码
//	Price	               float64    //上笔交易股价
//	Count                  int64      //上笔交易量
//	Date                   string     //上笔交易日期
//	AvgPrice               float64    //持仓均价
//	AllCount               int64      //持仓总股数
//	HoldMoney              float64    //持仓市值
//	Free                   float64    //手续费
//	Cost                   float64    //总成本
//	Profit	               float64    //盈亏
//	FloatProfit            float64    //浮动盈亏
//	FloatProfitRate	       float64    //浮动盈亏率
//	HoldDays               int64      //持股天数
//	TransType              string     //交易类型
//}


type TransParamInfo struct {
	Date	       string	    //交易日期
        Price          float64      //交易价格
        HighPrice      float64      //历史最高价
        LowPrice       float64      //历史最低价
}

type StockNDaysPreInfo struct {
	HighPrice	float64
	HighPriceDate	string
	LowPrice	float64
	LowPriceDate	string
}
