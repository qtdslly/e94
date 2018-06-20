package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Domain struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	Url         string     `gorm:"size:255" json:"url"`
	Disabled    bool       `json:"disabled"`
	DueDate     string     `gorm:"size:10" json:"due_date"`   //到期日
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Domain) TableName() string {
	return "domain"
}

func initDomain(db *gorm.DB) error {
	var err error

	if db.HasTable(&Domain{}) {
		err = db.AutoMigrate(&Domain{}).Error
	} else {
		err = db.CreateTable(&Domain{}).Error
	}
	return err
}

func dropDomain(db *gorm.DB) {
	db.DropTableIfExists(&Domain{})
}
