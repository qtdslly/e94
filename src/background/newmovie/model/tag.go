package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	Id            uint32    `gorm:"primary_key" json:"id"`
	ContentType   uint32    `json:"content_type"`
	ContentId     uint32    `json:"content_id"`
	Name          string    `gorm:"size:128" json:"name" valid:"Str" name:"name" len:"1,128"`
	Sort          uint32    `json:"sort"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tag"
}

func initTag(db *gorm.DB) error {
	var err error

	if db.HasTable(&Tag{}) {
		err = db.AutoMigrate(&Tag{}).Error
	} else {
		err = db.CreateTable(&Tag{}).Error
	}
	return err
}

func dropTag(db *gorm.DB) {
	db.DropTableIfExists(&Tag{})
}
