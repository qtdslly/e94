package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Picture struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	Provider       uint32           `json:"provider"`
	SourceId       string           `gorm:"size:60" json:"source_id"`
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
	IsMove         bool             `json:"is_move"`
	WaterMark      bool             `json:"water_mark"`
	Vertical       bool             `json:"vertical"` // 表示是否是竖屏
	Like           uint32           `json:"like"`
	Click          uint32           `json:"click"`
	Show           uint32           `json:"show"`
	Share          uint32           `json:"share"`
	Tags           []*Tag           `gorm:"many2many:picture_tag" json:"tags"`
	MoveId         uint32           `json:"move_id"`
	Move           *Move            `gorm:"-" json:"move"`
	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (Picture) TableName() string {
	return "picture"
}

func initPicture(db *gorm.DB) error {
	var err error
	if db.HasTable(&Picture{}) {
		err = db.AutoMigrate(&Picture{}).Error
	} else {
		err = db.CreateTable(&Picture{}).Error
	}
	return err
}

func dropPicture(db *gorm.DB) {
	db.DropTableIfExists(&Picture{})
}
