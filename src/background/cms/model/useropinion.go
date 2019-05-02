package model

import (
	"time"

	"github.com/jinzhu/gorm"
)


type UserOpinion struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	InstallationId uint64           `json:"user_id"`
	Category       uint8            `json:"category"`
	Description    string           `gorm:"type:longtext" json:"description"`
	Thumb          string           `gorm:"size:255" json:"thumb" valid:"Str" name:"thumb" len:"0,255"`
	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (UserOpinion) TableName() string {
	return "user_opinion"
}


func initUserOpinion(db *gorm.DB) error {
	var err error
	if db.HasTable(&UserOpinion{}) {
		err = db.AutoMigrate(&UserOpinion{}).Error
	} else {
		err = db.CreateTable(&UserOpinion{}).Error
	}
	return err
}

func dropUserOpinion(db *gorm.DB) {
	db.DropTableIfExists(&UserOpinion{})
}
