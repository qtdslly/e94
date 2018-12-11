package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Photo struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	Url         string     `gorm:"size:100" json:"url"`
	State       int        `json:"state"`
	Title       string     `gorm:"size:255" json:"title"`
	Description string     `gorm:"type:longtext" json:"description"`
	Count       int        `json:"count"`
	CategoryId  uint32     `json:"category_id"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

const (
	PhotoStateOnLine = 1
	PhotoStateOffLine = 0
)
func (Photo) TableName() string {
	return "photo"
}

func initPhoto(db *gorm.DB) error {
	var err error

	if db.HasTable(&Photo{}) {
		err = db.AutoMigrate(&Photo{}).Error
	} else {
		err = db.CreateTable(&Photo{}).Error
	}
	return err
}

func dropPhoto(db *gorm.DB) {
	db.DropTableIfExists(&Photo{})
}
