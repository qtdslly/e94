package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type DownloadUrl struct {
	Id            uint32     `gorm:"primary_key" json:"id"`
	Provider      uint32     `json:"provider"`
	Title         string     `gorm:"size:255" json:"title"`
	VideoId       uint32     `json:"video_id"`
	Sort          uint32     `json:"sort"`
	Url           string     `gorm:"size:255" json:"url"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (DownloadUrl) TableName() string {
	return "download_url"
}

func initDownloadUrl(db *gorm.DB) error {
	var err error

	if db.HasTable(&DownloadUrl{}) {
		err = db.AutoMigrate(&DownloadUrl{}).Error
	} else {
		err = db.CreateTable(&DownloadUrl{}).Error
	}
	return err
}

func dropDownloadUrl(db *gorm.DB) {
	db.DropTableIfExists(&DownloadUrl{})
}
