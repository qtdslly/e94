package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Video struct {
	Id          uint32           `gorm:"primary_key" json:"id"`
	Provider    uint32           `json:"provider"`
	SourceId    string           `gorm:"size:255" json:"source_id"`
	Title       string           `gorm:"size:255" json:"title" translated:"true"`
	Description string           `gorm:"type:longtext" json:"description" translated:"true"`
	ThumbX      string           `gorm:"size:255;column:thumb_x" json:"thumb_x"`
	ThumbY      string           `gorm:"size:255;column:thumb_y" json:"thumb_y"`
	Plays       uint32           `json:"plays"`
	Diggs       uint32           `json:"diggs"`
	Filesize    uint32           `json:"filesize"`
	Height      uint32           `json:"height"`
	Width       uint32           `json:"width"`
	Duration    uint32           `json:"duration"` //时长
	Url         string           `gorm:"size:255" json:"url"`
	Status      uint32           `json:"status"`
	Watermark   bool             `json:"watermark"`
	Country     string           `gorm:"size:255" json:"country"`
	Province    string           `gorm:"size:255" json:"province"`
	City        string           `gorm:"size:255" json:"city"`
	District    string           `gorm:"size:255" json:"district"`
	Address     string           `gorm:"size:255" json:"address"`
	Longitude   float64          `json:"longitude"`
	Latitude    float64          `json:"latitude"`
	Vertical    bool             `json:"vertical"` // 表示是否是竖屏
	Tags        []*Tag           `gorm:"many2many:video_tag" json:"tags"`
	PersonId    uint32           `json:"person_id"`
	CategoryId  uint32           `json:"category_id"`
	Author      *Person          `gorm:"-" json:"author"`
	ReleasedAt  *time.Time       `json:"released_at"` // 发布时间，utc格式
	SyncedAt    *time.Time       `json:"created_at"`  // 创建时间，utc格式
	CreatedAt   time.Time        `json:"created_at"`  // 创建时间，utc格式
	UpdatedAt   time.Time        `json:"updated_at"`  // 更新时间，utc格式
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
