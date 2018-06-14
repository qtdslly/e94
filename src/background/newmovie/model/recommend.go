package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Recommend struct {
	Id              uint32     `gorm:"primary_key" json:"id"`
	ContentType     uint32     `json:"content_type"`
	ContentId       uint32     `json:"content_id"`
	Title           string     `gorm:"size:60" json:"title"`
	Focus           string     `gorm:"size:60" json:"focus"`
	ThumbX          string     `gorm:"size:60" json:"thumb_x"`
	ThumbY          string     `gorm:"size:60" json:"thumb_y"`
	Status          uint32     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Recommend) TableName() string {
	return "recommend"
}

func initRecommend(db *gorm.DB) error {
	var err error

	if db.HasTable(&Recommend{}) {
		err = db.AutoMigrate(&Recommend{}).Error
	} else {
		err = db.CreateTable(&Recommend{}).Error
	}
	return err
}

func dropRecommend(db *gorm.DB) {
	db.DropTableIfExists(&Recommend{})
}
