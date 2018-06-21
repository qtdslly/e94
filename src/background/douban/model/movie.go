package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Movie struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	PageId      uint32     `json:"page_id"`
	Title       string     `gorm:"size:255" json:"title"`
	Description string     `gorm:"type:longtext" json:"description"`
	Score       float32    `json:"score"`
	Directors   string     `gorm:"size:255" json:"directors"`
	Writer      string     `gorm:"size:255" json:"writer"`
	Types       string     `gorm:"size:100" json:"types"`
	Country     string     `gorm:"size:100" json:"country"`
	Language    string     `gorm:"size:255" json:"language"`
	ReleaseDate string     `gorm:"size:60" json:"release_date"`
	Duration    uint32     `json:"duration"`
	Alias       string     `gorm:"size:255" json:"alias"`
	Imdb        string     `gorm:"size:255" json:"imdb"`
	Comments    uint32     `json:"comments"`   //评论人数

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
