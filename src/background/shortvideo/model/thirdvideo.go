package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ThirdVideo struct {
	Id                 uint32  `gorm:"primary_key" json:"id"`
	Provider           string  `json:"provider"`
	ThirdVideoId       string  `json:"third_video_id"`
	Title              string  `gorm:"size:255" json:"title" translated:"true"`
	Description        string  `gorm:"type:longtext" json:"description" translated:"true"`
	ThumbX             string  `gorm:"size:255;column:thumb_x" json:"thumb_x"`
	ThumbY             string  `gorm:"size:255;column:thumb_y" json:"thumb_y"`
	ThirdAuthorId      string  `json:"third_author_id"`
	ThirdAuthorShortId string  `json:"third_author_short_id"`
	NickName           string  `json:"nick_name"`
	PlayCount          uint32  `json:"play_count"`
	CommentCount       uint32  `json:"comment_count"`
	DiggCount          uint32  `json:"digg_count"`
	ShareCount         uint32  `json:"share_count"`
	Filesize           uint64  `json:"file_size"`
	Height             uint32  `json:"height"`
	Width              uint32  `json:"width"`
	Duration           string  `json:"duration"` //时长
	ThirdId            string  `json:"third_id"`
	Playurl            string  `json:"play_url"`
	SourceUrl          string  `json:"source_url"`
	FileName           string  `json:"file_name"`
	Location           string  `json:"location"`
	AuthorThumb        string  `json:"author_thumb"`
	ShareTitle         string  `json:"share_title"`
	ShareDescription   string  `json:"share_description"`
	Birthday           string  `json:"birthday"`
	HasWaterMark       bool    `json:"has_water_mark"`
	Country            string  `json:"province"`
	Province           string  `json:"province"`
	City               string  `json:"city"`
	SimpleAddr         string  `json:"simple_addr"`
	District           string  `json:"district"`
	Address            string  `json:"address"`
	Longitude          float64 `json:"longitude"`
	Latitude           float64 `json:"latitude"`
	Tag                string  `json:"tag"`
	IsVerticalScreen   bool    `json:"is_vertical_screen"`
	Category           string  `json:"category"`

	ThirdCreatedAt string    `json:"third_created_at"`
	CreatedAt      time.Time `json:"created_at"` // 创建时间，utc格式
	UpdatedAt      time.Time `json:"updated_at"` // 更新时间，utc格式

}

func (ThirdVideo) TableName() string {
	return "third_video"
}

func InitThirdVideo(db *gorm.DB) error {
	var err error
	if db.HasTable(&ThirdVideo{}) {
		err = db.AutoMigrate(&ThirdVideo{}).Error
	} else {
		err = db.CreateTable(&ThirdVideo{}).Error
	}
	return err
}

func dropThirdVideo(db *gorm.DB) {
	db.DropTableIfExists(&ThirdVideo{})
}
