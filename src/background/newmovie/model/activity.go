package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Activity struct {
	Id             uint32      `gorm:"primary_key" json:"id"`
	AppId          uint32      `json:"app_id"`
	VersionId      uint32      `json:"version_id"`
	Channel        uint32      `json:"channel"`
	Account        string      `json:"account"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Thumb          string      `json:"thumb"`
	Enable         bool        `json:"enable"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

const (
	ChannelZhiFuBao    = 1
	ChannelWeiXin      = 2
)

func (Activity) TableName() string {
	return "activity"
}

func initActivity(db *gorm.DB) error {
	var err error
	if db.HasTable(&Activity{}) {
		err = db.AutoMigrate(&Activity{}).Error
	} else {
		err = db.CreateTable(&Activity{}).Error
	}

	return err
}

func dropActivity(db *gorm.DB) {
	db.DropTableIfExists(&Activity{})
}
