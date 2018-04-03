package model

import (
	"github.com/jinzhu/gorm"
)
type StockList struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Code                   string     `gorm:"code" json:"code"`                            //股票代码
	Name                   string     `gorm:"name" json:"name"`                            //股票名称
	Industry               string     `gorm:"industry" json:"industry"`                    //所属行业
	Area                   string     `gorm:"area" json:"area"`                            //地区
	Pe                     string     `gorm:"pe" json:"pe"`                                //市盈率
	Outstanding            string     `gorm:"outstanding" json:"outstanding"`              //流通股本(亿)
	Totals                 string     `gorm:"totals" json:"totals"`                        //总股本(亿)
	TotalAssets            string     `gorm:"totalAssets" json:"totalAssets"`              //总资产(万)
	LiquidAssets           string     `gorm:"liquidAssets" json:"liquidAssets"`            //流动资产
	FixedAssets            string     `gorm:"fixedAssets" json:"fixedAssets"`              //固定资产
	Reserved               string     `gorm:"reserved" json:"reserved"`                    //公积金
	ReservedPerShare       string     `gorm:"reservedPerShare" json:"reservedPerShare"`    //每股公积金
	Esp                    string     `gorm:"esp" json:"esp"`                              //每股收益
	Bvps                   string     `gorm:"bvps" json:"bvps"`                            //每股净资
	Pb                     string     `gorm:"pb" json:"pb"`                                //市净率
	TimeToMarket           string     `gorm:"timeToMarket" json:"timeToMarket"`            //上市日期
	Undp                   string     `gorm:"undp" json:"undp"`                            //未分利润
	Perundp                string     `gorm:"perundp" json:"perundp"`                      //每股未分配
	Rev                    string     `gorm:"rev" json:"rev"`                              //收入同比(%)
	Profit                 string     `gorm:"profit" json:"profit"`                        //利润同比(%)
	Gpr                    string     `gorm:"gpr" json:"gpr"`                              //毛利率(%)
	Npr                    string     `gorm:"npr" json:"npr"`                              //净利润率(%)
	Holders                string     `gorm:"holders" json:"holders"`                      //股东人数
}

func (StockList) TableName() string {
	return "stock_list"
}

func initStockList(db *gorm.DB) error {
	var err error

	if db.HasTable(&StockList{}) {
		err = db.AutoMigrate(&StockList{}).Error
	} else {
		err = db.CreateTable(&StockList{}).Error
	}
	return err
}

func dropStockList(db *gorm.DB) {
	db.DropTableIfExists(&StockList{})
}
