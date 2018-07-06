package model

import (
	"time"

	"github.com/jinzhu/gorm"
)


type UserStream struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	InstallationId uint64           `json:"user_id"`
	Title          string           `gorm:"size:60" json:"title" valid:"Str" name:"title" len:"1,255"`
	Pinyin         string           `gorm:"size:32;index" json:"pinyin"`
	Sort           uint32           `json:"sort"`
	Icon           string           `gorm:"size:255" json:"icon" valid:"Str" name:"icon" len:"0,255"`
	Thumb          string           `gorm:"size:255" json:"thumb" valid:"Str" name:"thumb" len:"0,255"`
	Url            string           `gorm:"size:1024" json:"url"`
	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (UserStream) TableName() string {
	return "user_stream"
}


func initUserStream(db *gorm.DB) error {
	var err error
	if db.HasTable(&UserStream{}) {
		err = db.AutoMigrate(&UserStream{}).Error
	} else {
		err = db.CreateTable(&UserStream{}).Error
	}
	return err
}

func dropUserStream(db *gorm.DB) {
	db.DropTableIfExists(&UserStream{})
}
