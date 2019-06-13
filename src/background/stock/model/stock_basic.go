package model

import (
  "github.com/jinzhu/gorm"
  "time"
)

type StockBasic struct {
  Id           uint32     `gorm:"primary_key" json:"id"`
  Jys          string     `gorm:"jys" json:"jys"`                       //交易所 0 上海 1 深证
  Code         string     `gorm:"code" json:"code"`                     //股票代码
  Name         string     `gorm:"name" json:"name"`                     //股票名称
  ToMarketDate string     `gorm:"to_market_date" json:"to_market_date"` //上市日期
  Company      string     `gorm:"company" json:"company"`               //公司名称
  Zyyw         string     `gorm:"zyyw" json:"zyyw"`                     //主营业务
  Dsz          string     `gorm:"dsz" json:"dsz"`                       //董事长
  Zjl          string     `gorm:"zjl" json:"zjl"`                       //总经理
  Dmdh         string     `gorm:"dmdh" json:"dmdh"`                     //董秘电话
  Address      string     `gorm:"address" json:"address"`               //办公地址
  WebAddress   string     `gorm:"web_address" json:"web_address"`       //公司网址
  Zczb         string     `gorm:"zczb" json:"zczb"`                     //注册资本
  Zgb          string     `gorm:"zgb" json:"zgb"`                       //总股本
  Ltgb         string     `gorm:"ltgb" json:"ltgb"`                     //流通股本


  CreatedAt    time.Time `json:"created_at"`
  UpdatedAt    time.Time `json:"updated_at"`
}

func (StockBasic) TableName() string {
  return "stock_basic"
}

func initStockBasic(db *gorm.DB) error {
  var err error

  if db.HasTable(&StockBasic{}) {
    err = db.AutoMigrate(&StockBasic{}).Error
  } else {
    err = db.CreateTable(&StockBasic{}).Error
  }
  return err
}

func dropStockBasic(db *gorm.DB) {
  db.DropTableIfExists(&StockBasic{})
}
