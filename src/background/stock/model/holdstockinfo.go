package model

import (
	"github.com/jinzhu/gorm"
)

type HoldStockInfo struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Code                   string     `gorm:"code" json:"code"`                            //股票代码
	Price                  float64    `gorm:"price" json:"price"`                          //上笔交易股价
	Count                  int64      `gorm:"count" json:"count"`                          //上笔交易量
	Date                   string     `gorm:"date" json:"date"`                            //上笔交易日期
	AvgPrice               float64    `gorm:"avg_price" json:"avg_price"`                  //持仓均价
	AllCount               int64      `gorm:"all_count" json:"all_count"`                  //持仓总股数
	HoldMoney              float64    `gorm:"hold_money" json:"hold_money"`                //持仓市值
	Fee                    float64    `gorm:"fee" json:"fee"`                              //手续费
	Cost                   float64    `gorm:"cost" json:"cost"`                            //总成本
	Profit	               float64    `gorm:"profit" json:"profit"`                        //盈亏
	FloatProfit            float64    `gorm:"float_profit" json:"float_profit"`            //浮动盈亏
	FloatProfitRate	       float64    `gorm:"float_profit_rate" json:"float_profit_rate"`  //浮动盈亏率
	HoldDays               int64      `gorm:"hold_days" json:"hold_days"`                  //持股天数
	TransType              string     `gorm:"trans_type" json:"trans_type"`                //交易类型
	SimulationId           uint32     `gorm:"simulation_id" json:"simulation_id"`          //模拟ID
}

func (HoldStockInfo) TableName() string {
	return "hold_stock_info"
}

func initHoldStockInfo(db *gorm.DB) error {
	var err error

	if db.HasTable(&HoldStockInfo{}) {
		err = db.AutoMigrate(&HoldStockInfo{}).Error
	} else {
		err = db.CreateTable(&HoldStockInfo{}).Error
	}
	return err
}

func dropHoldStockInfo(db *gorm.DB) {
	db.DropTableIfExists(&HoldStockInfo{})
}
