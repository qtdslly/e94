package model

import (
  "github.com/jinzhu/gorm"
  //"time"
  "time"
)

type BaiduZncp struct {
  Id             uint32     `gorm:"primary_key" json:"id"`
  Date           string     `gorm:"date" json:"date"`                         //股票代码
  Code           string     `gorm:"code" json:"code"`                         //股票代码
  Name           string     `gorm:"name" json:"name"`                         //股票代码
  Jys            string     `gorm:"jys" json:"jys"`                           //交易所
  Price          float64    `gorm:"price" json:"price"`                       //当前价
  Zlw            string    `gorm:"wlw" json:"wlw"`                           //阻力位
  Zcw            string    `gorm:"zcw" json:"zcw"`                           //支撑位
  Szgl           string    `gorm:"szgl" json:"szgl"`                         //上涨概率
  Kpld           string     `gorm:"kpld" json:"kpld"`                        //控盘力度
  Close          float64    `gorm:"close" json:"close"`                       //收盘价
  High           float64    `gorm:"high" json:"high"`                         //最高价
  Low            float64    `gorm:"low" json:"low"`                           //最低价
  Open           float64    `gorm:"open" json:"open"`                         //开盘价
  Volume         uint64    `gorm:"volume" json:"volume"`                      //开盘价
  Capitalization uint64    `gorm:"capitalization" json:"capitalization"`      //开盘价
  NetChange      float64    `gorm:"net_change" json:"net_change"`             //开盘价
  NetChangeRatio float64    `gorm:"net_change_ratio" json:"net_change_ratio"` //开盘价
  AmplitudeRatio float64    `gorm:"amplitude_ratio" json:"amplitude_ratio"`   //开盘价
  TurnoverRatio  float64    `gorm:"turnover_ratio" json:"turnover_ratio"`     //开盘价


  CreatedAt      time.Time `json:"created_at"`
  UpdatedAt      time.Time `json:"updated_at"`
}

func (BaiduZncp) TableName() string {
  return "baidu_zncp"
}

func initBaiduZncp(db *gorm.DB) error {
  var err error

  if db.HasTable(&BaiduZncp{}) {
    err = db.AutoMigrate(&BaiduZncp{}).Error
  } else {
    err = db.CreateTable(&BaiduZncp{}).Error
  }
  return err
}

func dropBaiduZncp(db *gorm.DB) {
  db.DropTableIfExists(&BaiduZncp{})
}
