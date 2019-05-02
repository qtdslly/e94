package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

/*
流通股東
*/
type ShareHolder struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Code                   string     `gorm:"code" json:"code"`
	Date                   string     `gorm:"date" json:"date"`
	Name                   string     `gorm:"name" json:"name"`
	HoldCount              string     `gorm:"hold_count" json:"hold_count"`
	Percent                string     `gorm:"percent" json:"percent"`
	Property               string     `gorm:"property" json:"property"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ShareHolder) TableName() string {
	return "share_holder"
}

func initShareHolder(db *gorm.DB) error {
	var err error

	if db.HasTable(&ShareHolder{}) {
		err = db.AutoMigrate(&ShareHolder{}).Error
	} else {
		err = db.CreateTable(&ShareHolder{}).Error
	}
	return err
}

func dropShareHolder(db *gorm.DB) {
	db.DropTableIfExists(&ShareHolder{})
}
