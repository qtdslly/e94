package model

import (
	"github.com/jinzhu/gorm"
	//"time"
  "time"
)

type DeepFallStock struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Code                   string     `gorm:"code" json:"code"`              //股票代码
	HighPrice              float64    `gorm:"high_price" json:"high_price"`  //最高价
	HighPriceDate          string     `gorm:"high_price_date" json:"high_price_date"`  //最高价日期
	LowPrice               float64    `gorm:"low_price" json:"low_price"`    //最低价
	LowPriceDate           string     `gorm:"low_price_date" json:"low_price_date"`  //最低价日期
	NowPrice               float64    `gorm:"now_price" json:"now_price"`    //当前价
	Date                   string     `gorm:"date" json:"date"`              //当前日期

  CreatedAt    time.Time `json:"created_at"`
  UpdatedAt    time.Time `json:"updated_at"`
}

func (DeepFallStock) TableName() string {
	return "deep_fall_stock"
}

func initDeepFallStock(db *gorm.DB) error {
	var err error

	if db.HasTable(&DeepFallStock{}) {
		err = db.AutoMigrate(&DeepFallStock{}).Error
	} else {
		err = db.CreateTable(&DeepFallStock{}).Error
	}
	return err
}

func dropDeepFallStock(db *gorm.DB) {
	db.DropTableIfExists(&DeepFallStock{})
}
