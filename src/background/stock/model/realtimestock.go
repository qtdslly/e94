package model

import (
	"github.com/jinzhu/gorm"
	//"time"
  "time"
)

type RealTimeStock struct {
	Id            uint32    `gorm:"primary_key" json:"id"`
	JysCode       string     `gorm:"jys_code:2" json:"jys_code"`           //交易所编码 :0 上海 1 深圳
	StockCode     string    `gorm:"stock_code:6" json:"stock_code"`     //股票代码
	TodayOpenPrice float64  `gorm:"today_open_price" json:"today_open_price"`     //今日开盘价
	YestdayClosePrice float64  `gorm:"yestday_close_price" json:"yestday_close_price"`     //昨日收盘价
	NowPrice float64  `gorm:"now_price" json:"now_price"`     //当前价
	TodayHighPrice float64  `gorm:"today_high_price" json:"today_high_price"`     //今日最高价
	TodayLowPrice float64  `gorm:"today_low_price" json:"today_low_price"`     //今日最低价
	BuyPrice float64  `gorm:"buy_price" json:"buy_price"`     //竞买价(买一)
	SellPrice float64  `gorm:"sell_price" json:"sell_price"`     //竞卖价(卖一)
	DealCount  int  `gorm:"deal_count" json:"deal_count"`     //成交数(单位为1股)
	DealMoney  float64  `gorm:"deal_money" json:"deal_money"`     //成交金额(单位为元)
	BuyCount1  int  `gorm:"buy_count1" json:"buy_count1"`     //买一申请股数(单位为1股)
	BuyPrice1  float64  `gorm:"buy_price1" json:"buy_price1"`     //买一报价
	BuyCount2  int  `gorm:"buy_count2" json:"buy_count2"`     //买一申请股数(单位为1股)
	BuyPrice2  float64  `gorm:"buy_price2" json:"buy_price2"`     //买一报价
	BuyCount3  int  `gorm:"buy_count3" json:"buy_count3"`     //买一申请股数(单位为1股)
	BuyPrice3  float64  `gorm:"buy_price3" json:"buy_price3"`     //买一报价
	BuyCount4  int  `gorm:"buy_count4" json:"buy_count4"`     //买一申请股数(单位为1股)
	BuyPrice4  float64  `gorm:"buy_price4" json:"buy_price4"`     //买一报价
	BuyCount5  int  `gorm:"buy_count5" json:"buy_count5"`     //买一申请股数(单位为1股)
	BuyPrice5  float64  `gorm:"buy_price5" json:"buy_price5"`     //买一报价
	SellCount1  int  `gorm:"sell_count1" json:"sell_count1"`     //买一申请股数(单位为1股)
	SellPrice1  float64  `gorm:"sell_price1" json:"sell_price1"`     //买一报价
	SellCount2  int  `gorm:"sell_count2" json:"sell_count2"`     //买一申请股数(单位为1股)
	SellPrice2  float64  `gorm:"sell_price2" json:"sell_price2"`     //买一报价
	SellCount3  int  `gorm:"sell_count3" json:"sell_count3"`     //买一申请股数(单位为1股)
	SellPrice3  float64  `gorm:"sell_price3" json:"sell_price3"`     //买一报价
	SellCount4  int  `gorm:"sell_count4" json:"sell_count4"`     //买一申请股数(单位为1股)
	SellPrice4  float64  `gorm:"sell_price4" json:"sell_price4"`     //买一报价
	SellCount5  int  `gorm:"sell_count5" json:"sell_count5"`     //买一申请股数(单位为1股)
	SellPrice5  float64  `gorm:"sell_price5" json:"sell_price5"`     //买一报价
	DealDate    string  `gorm:"deal_date:10" json:"deal_date"`     //交易日期
	DealTime    string  `gorm:"deal_time:30" json:"deal_time"`     //交易时间


  CreatedAt    time.Time `json:"created_at"`
  UpdatedAt    time.Time `json:"updated_at"`
}

func (RealTimeStock) TableName() string {
	return "real_time_stock"
}

func initRealTimeStock(db *gorm.DB) error {
	var err error

	if db.HasTable(&RealTimeStock{}) {
		err = db.AutoMigrate(&RealTimeStock{}).Error
	} else {
		err = db.CreateTable(&RealTimeStock{}).Error
	}
	return err
}

func dropRealTimeStock(db *gorm.DB) {
	db.DropTableIfExists(&RealTimeStock{})
}
