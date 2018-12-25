package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Duty struct {
	Id           uint32         `gorm:"primary_key" json:"id"`
	DoctorId     uint32         `json:"doctor_id"`
	Date         string         `json:"date"`
	Morning      bool           `json:"morning"`
	Afternoon    bool           `json:"afternoon"`
	Night        bool           `json:"night"`

	CreatedAt    time.Time      `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt    time.Time      `json:"updated_at"`       // 更新时间，utc格式
}

func (Duty) TableName() string {
	return "duty"
}

func initDuty(db *gorm.DB) error {
	var err error
	if db.HasTable(&Duty{}) {
		err = db.AutoMigrate(&Duty{}).Error
	} else {
		err = db.CreateTable(&Duty{}).Error
	}
	return err
}

func dropDuty(db *gorm.DB) {
	db.DropTableIfExists(&Duty{})
}
