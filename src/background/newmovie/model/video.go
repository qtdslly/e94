package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Video struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	Title          string           `gorm:"size:255" json:"title" valid:"Str" name:"title" len:"1,255" translated:"true"`
	Pinyin         string           `gorm:"size:32;index" json:"pinyin"`
	Description    string           `gorm:"type:longtext" json:"description" translated:"true"`
	Sort           uint32           `json:"sort"`
	TotalEpisode   uint32           `json:"total_episode"`
	CurrentEpisode uint32           `json:"current_episode"` // 当前播放最新的集数
	OnlineEpisode  uint32           `gorm:"-" json:"online_episode"`
	Episodes       []*Episode       `gorm:"-" json:"episodes"`
	Score          float64          `gorm:"type:float(3,1)" json:"score"`
	ThumbX         string           `gorm:"size:255;column:thumb_x" json:"thumb_x" valid:"Str" name:"thumb_x" len:"0,255"`
	ThumbY         string           `gorm:"size:255;column:thumb_y" json:"thumb_y" valid:"Str" name:"thumb_y" len:"0,255"`
	PublishDate    string           `json:"publish_date"`
	Category       string           `json:"category"`
	Status         uint32           `json:"status"`
	Year           uint32           `json:"year"`
	Language        string          `gorm:"size:60" json:"language"`
	Country        string           `gorm:"size:20" json:"country"`
	Directors      string           `gorm:"size:255" json:"directors"`
	Actors         string           `gorm:"size:255" json:"actors"`
	Tags           string           `gorm:"size:255" json:"tags"`

	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (Video) TableName() string {
	return "video"
}

func initVideo(db *gorm.DB) error {
	var err error
	if db.HasTable(&Video{}) {
		err = db.AutoMigrate(&Video{}).Error
	} else {
		err = db.CreateTable(&Video{}).Error
	}
	return err
}

func dropVideo(db *gorm.DB) {
	db.DropTableIfExists(&Video{})
}
