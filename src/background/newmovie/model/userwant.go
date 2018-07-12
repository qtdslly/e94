package model

import (
	"time"

	"github.com/jinzhu/gorm"
)


type UserWant struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	InstallationId uint64           `json:"installation_id"`
	Title          string           `gorm:"size:60" json:"title" valid:"Str" name:"title" len:"1,255"`
	Description    string           `gorm:"type:longtext" json:"description" translated:"true"`
	Email          string           `gorm:"size:255" json:"icon" valid:"Str" name:"email" len:"0,255"`
	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (UserWant) TableName() string {
	return "user_want"
}


func initUserWant(db *gorm.DB) error {
	var err error
	if db.HasTable(&UserWant{}) {
		err = db.AutoMigrate(&UserWant{}).Error
	} else {
		err = db.CreateTable(&UserWant{}).Error
	}
	return err
}

func dropUserWant(db *gorm.DB) {
	db.DropTableIfExists(&UserWant{})
}
