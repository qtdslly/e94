package model

import (
  "github.com/jinzhu/gorm"
  "time"
)

type TonghuashunSstd struct {
  Id        uint32     `gorm:"primary_key" json:"id"`
  Date      string     `gorm:"date" json:"date"`    //日期
  Code      string     `gorm:"code" json:"code"`    //股票代码
  Name      string     `gorm:"name" json:"name"`    //股票名称
  Jys       string     `gorm:"jys" json:"jys"`      //交易所
  Price     string     `gorm:"price" json:"price"`  //现价
  Zdf       string     `gorm:"zdf" json:"zdf"`      //涨跌幅
  Mrxh      string     `gorm:"mrxh" json:"mrxh"`    //买入信号
  Jsxt      string     `gorm:"jsxt" json:"jsxt"`    //技术形态


  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

func (TonghuashunSstd) TableName() string {
  return "tonghuashun_sstd"
}

func initTonghuashunSstd(db *gorm.DB) error {
  var err error

  if db.HasTable(&TonghuashunSstd{}) {
    err = db.AutoMigrate(&TonghuashunSstd{}).Error
  } else {
    err = db.CreateTable(&TonghuashunSstd{}).Error
  }
  return err
}

func dropTonghuashunSstd(db *gorm.DB) {
  db.DropTableIfExists(&TonghuashunSstd{})
}
