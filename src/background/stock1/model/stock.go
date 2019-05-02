package model

import (
	"github.com/jinzhu/gorm"
	"time"
)
type Stock struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Jys                    string     `gorm:"jys" json:"jys"`                            //交易所 0 上海 1 深证
	Code                   string     `gorm:"code" json:"code"`                            //股票代码
	Name                   string     `gorm:"name" json:"name"`                            //股票名称
	Price                  string     `gorm:"price" json:"industry"`                       //价格
	Percent                   string     `gorm:"percent" json:"percent"`                            //涨跌幅
	UpDown              string     `gorm:"up_down" json:"up_down"`                //涨跌额
	FiveMinute             string     `gorm:"five_minute" json:"five_minute"`              //5分钟涨跌幅
	Open                   string     `gorm:"open" json:"open"`                            //开盘价
	YestClose             string     `gorm:"yest_close" json:"yest_close"`              //作日收盘价
	High                   string     `gorm:"high" json:"high"`                            //最高价
	Low                    string     `gorm:"low" json:"low"`                              //最低价
	Volume              string     `gorm:"volume" json:"volume"`                  //成交量
	Turnover           string     `gorm:"turnover" json:"turnover"`          //成交额
	Hs          string     `gorm:"hs" json:"hs"`        //换手率
	Lb          string     `gorm:"lb" json:"lb"`        //量比
	Wb          string     `gorm:"wb" json:"wb"`        //委比
	Zf          string     `gorm:"zf" json:"zf"`        //振幅
	Pe                     string     `gorm:"pe" json:"pe"`                                //市盈率
	Mcap          string     `gorm:"mcap" json:"mcap"`        //流通市值
	Tcap          string     `gorm:"tcap" json:"tcap"`        //总市值
	Eps          string     `gorm:"eps" json:"eps"`        //每股收益
	NetProfit          string     `gorm:"net_profit" json:"net_profit"`        //净利润
	TotalRevenue          string     `gorm:"total_revenue" json:"total_revenue"`        //总营收

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Stock) TableName() string {
	return "stock"
}

func initStock(db *gorm.DB) error {
	var err error

	if db.HasTable(&Stock{}) {
		err = db.AutoMigrate(&Stock{}).Error
	} else {
		err = db.CreateTable(&Stock{}).Error
	}
	return err
}

func dropStock(db *gorm.DB) {
	db.DropTableIfExists(&Stock{})
}
