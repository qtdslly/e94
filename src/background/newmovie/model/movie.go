package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Movie struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Title       string           `gorm:"size:255" json:"title" translated:"true"`
	Description string           `gorm:"type:longtext" json:"description" translated:"true"`
	ThumbX      string           `gorm:"size:255;column:thumb_x" json:"thumb_x"`
	ThumbY      string           `gorm:"size:255;column:thumb_y" json:"thumb_y"`
	Url                    string    `json:"url"`
	Directors		       string     `json:"directors"`
	Actors		       string     `json:"actors"`
	Score		       string     `json:"score"`
	PublishDate            string     `json:"publish_date"`

	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
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
