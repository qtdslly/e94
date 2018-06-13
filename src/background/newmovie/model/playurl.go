package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type PlayUrl struct {
	Id          uint32         `gorm:"primary_key" json:"id"`
	Provider    uint32         `json:"provider"`
	ContentType uint8          `gorm:"index" json:"content_type"`
	ContentId   uint32         `gorm:"index" json:"content_id"`
	Title       string         `gorm:"size:255" json:"title"`
	Ratio       string         `json:"ratio"`
	Width       uint32         `json:"width"`
	Height      uint32         `json:"height"`
	Bitrate     uint32         `json:"bitrate"`
	Url         string         `gorm:"size:255" json:"url"`
	Disabled    bool           `json:"disabled"` // 链接播放不了，临时禁止
	Quality     uint8          `json:"quality"`
	Sort        uint32         `json:"sort"`
	CreatedAt   time.Time      `json:"created_at"` // 创建时间，utc格式
	UpdatedAt   time.Time      `json:"updated_at"` // 更新时间，utc格式
}

func (PlayUrl) TableName() string {
	return "play_url"
}

func initPlayUrl(db *gorm.DB) error {
	var err error
	if db.HasTable(&PlayUrl{}) {
		err = db.AutoMigrate(&PlayUrl{}).Error
	} else {
		err = db.CreateTable(&PlayUrl{}).Error
	}
	return err
}

func dropPlayUrl(db *gorm.DB) {
	db.DropTableIfExists(&PlayUrl{})
}
