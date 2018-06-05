package model

import (
	"github.com/jinzhu/gorm"
)

type Movie struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Title                  string     `json:"title"`
	Url                    string    `json:"price"`
}

func (Movie) TableName() string {
	return "movie"
}

func initMovie(db *gorm.DB) error {
	var err error

	if db.HasTable(&Movie{}) {
		err = db.AutoMigrate(&Movie{}).Error
	} else {
		err = db.CreateTable(&Movie{}).Error
	}
	return err
}

func dropMovie(db *gorm.DB) {
	db.DropTableIfExists(&Movie{})
}
