package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Cartoon struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	Url         string     `gorm:"size:255" json:"url"`
	PageStatus  uint32     `json:"page_status"`
	UrlStatus   uint32     `json:"url_status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Cartoon) TableName() string {
	return "cartoon"
}

func initCartoon(db *gorm.DB) error {
	var err error

	if db.HasTable(&Cartoon{}) {
		err = db.AutoMigrate(&Cartoon{}).Error
	} else {
		err = db.CreateTable(&Cartoon{}).Error
	}
	return err
}

func dropCartoon(db *gorm.DB) {
	db.DropTableIfExists(&Cartoon{})
}
