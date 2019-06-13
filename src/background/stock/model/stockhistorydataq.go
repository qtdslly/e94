package model

import (
	"github.com/jinzhu/gorm"
)

type StockHistoryDataQ struct {
	Id            uint32     `gorm:"primary_key" json:"id"`
	Date          string     `gorm:"date:10" json:"date"`           //交易日期
	Open          float64    `gorm:"open" json:"open"`              //今日开盘价
	Close         float64    `gorm:"close" json:"close"`            //今日收盘价
	High          float64    `gorm:"high" json:"high"`              //今日最高价
	Low           float64    `gorm:"low" json:"low"`                //今日最低价
	Volume        float64    `gorm:"volume" json:"volume"`          //成交量
	Code          string     `gorm:"code" json:"code"`              //股票代码
}

func (StockHistoryDataQ) TableName() string {
	return "stock_history_data_q"
}

func initStockHistoryDataQ(db *gorm.DB) error {
	var err error

	if db.HasTable(&StockHistoryDataQ{}) {
		err = db.AutoMigrate(&StockHistoryDataQ{}).Error
	} else {
		err = db.CreateTable(&StockHistoryDataQ{}).Error
	}
	return err
}

func dropStockHistoryDataQ(db *gorm.DB) {
	db.DropTableIfExists(&StockHistoryDataQ{})
}
