package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Episode struct {
	Id           uint32         `gorm:"primary_key" json:"id"`
	VideoId      uint32         `gorm:"index" json:"video_id"`
	Title        string         `gorm:"size:255" json:"title" valid:"Str" name:"title" len:"1,255" translated:"true"`
	Pinyin       string         `gorm:"size:32;index" json:"pinyin"`
	Description  string         `gorm:"type:longtext" json:"description" translated:"true"`
	Sort         uint32         `json:"sort"`
	Number       string         `gorm:"size:64" json:"number"`
	Score        float64        `gorm:"type:float(3,1)" json:"score"`
	ThumbX       string         `gorm:"size:255;column:thumb_x" json:"thumb_x" valid:"Str" name:"thumb_x" len:"0,255"`
	ThumbY       string         `gorm:"size:255;column:thumb_y" json:"thumb_y" valid:"Str" name:"thumb_y" len:"0,255"`
	ThumbOttX    string         `gorm:"size:255;column:thumb_ott_x" json:"thumb_ott_x" valid:"Str" name:"thumb_ott_x" len:"0,255"`
	ThumbOttY    string         `gorm:"size:255;column:thumb_ott_y" json:"thumb_ott_y" valid:"Str" name:"thumb_ott_y" len:"0,255"`
	Duration     uint32         `json:"duration"`
	PublishDate  string         `gorm:"size:20" json:"publish_date"`
	PlayUrls     []*PlayUrl     `gorm:"-" json:"play_urls"`
	CreatedAt    time.Time      `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt    time.Time      `json:"updated_at"`       // 更新时间，utc格式
}

func (Episode) TableName() string {
	return "episode"
}

func initEpisode(db *gorm.DB) error {
	var err error
	if db.HasTable(&Episode{}) {
		err = db.AutoMigrate(&Episode{}).Error
	} else {
		err = db.CreateTable(&Episode{}).Error
	}
	return err
}

func dropEpisode(db *gorm.DB) {
	db.DropTableIfExists(&Episode{})
}
