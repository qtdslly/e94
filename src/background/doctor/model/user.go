package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	Id           uint32         `gorm:"primary_key" json:"id"`
	OpenId       string         `gorm:"size:60" json:"open_id"`
	Avtar        string         `gorm:"size:255" json:"avtar"`
	Nick         string         `gorm:"size:60" json:"nick"`
	Country      string         `gorm:"size:60" json:"country"`
	Province     string         `gorm:"size:60" json:"province"`
	City         string         `gorm:"size:60" json:"city"`
	Language     string         `gorm:"size:60" json:"language"`
	Gender       string         `gorm:"size:2" json:"gender"`

	CreatedAt    time.Time      `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt    time.Time      `json:"updated_at"`       // 更新时间，utc格式
}

func (User) TableName() string {
	return "user"
}

func initUser(db *gorm.DB) error {
	var err error
	if db.HasTable(&User{}) {
		err = db.AutoMigrate(&User{}).Error
	} else {
		err = db.CreateTable(&User{}).Error
	}
	return err
}

func dropUser(db *gorm.DB) {
	db.DropTableIfExists(&User{})
}
