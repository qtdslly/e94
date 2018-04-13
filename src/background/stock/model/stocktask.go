package model

import (
	"github.com/jinzhu/gorm"
)

type StockTask struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Key                    string     `gorm:"key" json:"key"`
	Value                  string     `gorm:"value" json:"value"`
	Date                   string     `gorm:"date" json:"date"`
}

func (StockTask) TableName() string {
	return "stock_task"
}

func initStockTask(db *gorm.DB) error {
	var err error

	if db.HasTable(&StockTask{}) {
		err = db.AutoMigrate(&StockTask{}).Error
	} else {
		err = db.CreateTable(&StockTask{}).Error
	}
	return err
}

func dropStockTask(db *gorm.DB) {
	db.DropTableIfExists(&StockTask{})
}
