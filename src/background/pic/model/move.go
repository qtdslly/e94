package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Move struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	Provider       uint32           `json:"provider"`
	Title          string           `gorm:"size:255" json:"title"`
	Pinyin         string           `gorm:"size:32;index" json:"pinyin"`
	Description    string           `gorm:"type:longtext" json:"description" translated:"true"`
	Sort           uint32           `json:"sort"`
	CategoryId     uint32           `json:"category_id"`
	OnLine         bool             `json:"on_line"`
	Url            string           `gorm:"size:255" json:"url"`
	SourceUrl      string           `gorm:"size:255" json:"source_url"`
	Filesize       uint32           `json:"filesize"`
	Height         uint32           `json:"height"`
	Width          uint32           `json:"width"`
	Duration       float64          `json:"duration"`
	WaterMark      bool             `json:"water_mark"`
	Vertical       bool             `json:"vertical"` // 表示是否是竖屏
	Like           uint32           `json:"like"`
	Click          uint32           `json:"click"`
	Show           uint32           `json:"show"`
	Share          uint32           `json:"share"`
	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (Move) TableName() string {
	return "move"
}

func initMove(db *gorm.DB) error {
	var err error
	if db.HasTable(&Move{}) {
		err = db.AutoMigrate(&Move{}).Error
	} else {
		err = db.CreateTable(&Move{}).Error
	}
	return err
}

func dropMove(db *gorm.DB) {
	db.DropTableIfExists(&Move{})
}
