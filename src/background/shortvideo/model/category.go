package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Category struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	ContentType uint8      `json:"content_type"`
	Name        string     `gorm:"size:64;unique" json:"name" valid:"Str" name:"name" len:"1,64"`
	Sort        uint32     `json:"sort"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	SyncedAt    *time.Time `json:"synced_at"` // 同步时间，utc格式
}

func (Category) TableName() string {
	return "category"
}

func initCategory(db *gorm.DB) error {
	var err error

	if db.HasTable(&Category{}) {
		err = db.AutoMigrate(&Category{}).Error
	} else {
		err = db.CreateTable(&Category{}).Error
	}
	return err
}

func dropCategory(db *gorm.DB) {
	db.DropTableIfExists(&Category{})
}
