package model

import (
	"github.com/jinzhu/gorm"
)

type TransStockInfo struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Code                   string     `gorm:"code" json:"code"`                            //股票代码
	Price                  float64    `gorm:"price" json:"price"`                          //上笔交易股价
	Count                  int64      `gorm:"count" json:"count"`                          //上笔交易量
	Date                   string     `gorm:"date" json:"date"`                            //上笔交易日期
	Fee                    float64    `gorm:"fee" json:"fee"`                              //手续费
	Cost                   float64    `gorm:"cost" json:"cost"`                            //总成本
	Profit	               float64    `gorm:"profit" json:"profit"`                        //盈亏
	TransType              string     `gorm:"trans_type" json:"trans_type"`                //交易类型
	SimulationId           uint32     `gorm:"simulation_id" json:"simulation_id"`             //模拟ID
}

func (TransStockInfo) TableName() string {
	return "trans_stock_info"
}

func initTransStockInfo(db *gorm.DB) error {
	var err error

	if db.HasTable(&TransStockInfo{}) {
		err = db.AutoMigrate(&TransStockInfo{}).Error
	} else {
		err = db.CreateTable(&TransStockInfo{}).Error
	}
	return err
}

func dropTransStockInfo(db *gorm.DB) {
	db.DropTableIfExists(&TransStockInfo{})
}
