package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Page struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	Url         string     `gorm:"size:255;unique" json:"name"`
	PageStatus  uint32     `json:"page_status"`
	UrlStatus   uint32     `json:"url_status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Page) TableName() string {
	return "page"
}

func initPage(db *gorm.DB) error {
	var err error

	if db.HasTable(&Page{}) {
		err = db.AutoMigrate(&Page{}).Error
	} else {
		err = db.CreateTable(&Page{}).Error
	}
	return err
}

func dropPage(db *gorm.DB) {
	db.DropTableIfExists(&Page{})
}
