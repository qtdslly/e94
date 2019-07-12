package model

import (
  "github.com/jinzhu/gorm"
  "time"
)

//同花顺承压分析
type TonghuashunCyfx struct {
  Id        uint32     `gorm:"primary_key" json:"id"`
  Date      string     `gorm:"date" json:"date"` //日期
  Code      string     `gorm:"code" json:"code"` //股票代码
  Content   string     `gorm:"content" json:"content"`

  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

func (TonghuashunCyfx) TableName() string {
  return "tonghuashun_cyfx"
}

func initTonghuashunCyfx(db *gorm.DB) error {
  var err error

  if db.HasTable(&TonghuashunCyfx{}) {
    err = db.AutoMigrate(&TonghuashunCyfx{}).Error
  } else {
    err = db.CreateTable(&TonghuashunCyfx{}).Error
  }
  return err
}

func dropTonghuashunCyfx(db *gorm.DB) {
  db.DropTableIfExists(&TonghuashunCyfx{})
}
