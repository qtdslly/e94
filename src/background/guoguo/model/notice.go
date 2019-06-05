package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Notice struct {
	Id        uint32  `gorm:"primary_key" json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
	BuyCount  int     `json:"buy_count"`
	SellCount int     `json:"sell_count"`
	State     bool    `json:"state"`
	Frequency int     `json:"frequency"`     //频率 0 高频 1 低频  高频实时刷新预警 低频每小时刷新一次预警

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Notice) TableName() string {
	return "notice"
}

func initNotice(db *gorm.DB) error {
	var err error
	if db.HasTable(&Notice{}) {
		err = db.AutoMigrate(&Notice{}).Error
	} else {
		err = db.CreateTable(&Notice{}).Error
	}
	return err
}

func dropNotice(db *gorm.DB) {
	db.DropTableIfExists(&Notice{})
}
