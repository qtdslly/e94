package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TopSearch struct {
	Id              uint32           `gorm:"primary_key" json:"id"`
	Title           string           `gorm:"size:255" json:"title" translated:"true"`
	Sort            uint32           `json:"sort"`
	ContentType     uint8            `json:"content_type"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (TopSearch) TableName() string {
	return "top_search"
}

func initTopSearch(db *gorm.DB) error {
	var err error

	if db.HasTable(&TopSearch{}) {
		err = db.AutoMigrate(&TopSearch{}).Error
	} else {
		err = db.CreateTable(&TopSearch{}).Error
	}
	return err
}

func dropTopSearch(db *gorm.DB) {
	db.DropTableIfExists(&TopSearch{})
}
