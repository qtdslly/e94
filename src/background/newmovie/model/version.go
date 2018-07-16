package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Version struct {
	Id        uint32    `gorm:"primary_key" json:"id"`
	AppId     uint32    `json:"app_id" binding:"required"`
	AppKey    string    `gorm:"size:64;unique_index:idx_version" json:"app_key"`
	Name      string    `gorm:"size:64" json:"name" valid:"Str" name:"name" len:"1,64"`
	Version   string    `gorm:"size:32;unique_index:idx_version" json:"version" valid:"Str" name:"version" len:"1,32"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Version) TableName() string {
	return "version"
}

func initVersion(db *gorm.DB) error {
	var err error
	if db.HasTable(&Version{}) {
		err = db.AutoMigrate(&Version{}).Error
	} else {
		err = db.CreateTable(&Version{}).Error
	}

	return err
}

func dropVersion(db *gorm.DB) {
	db.DropTableIfExists(&Version{})
}
