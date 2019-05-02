package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Top struct {
	Id            uint32    `gorm:"primary_key" json:"id"`
	ContentType   uint32    `json:"content_type"`
	ContentId     uint32    `json:"content_id"`
	Name          string    `gorm:"size:128" json:"name" valid:"Str" name:"name" len:"1,128"`
	Sort          uint32    `json:"sort"`
	State         uint32    `json:"state"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Top) TableName() string {
	return "top"
}

func initTop(db *gorm.DB) error {
	var err error

	if db.HasTable(&Top{}) {
		err = db.AutoMigrate(&Top{}).Error
	} else {
		err = db.CreateTable(&Top{}).Error
	}
	return err
}

func dropTop(db *gorm.DB) {
	db.DropTableIfExists(&Top{})
}
