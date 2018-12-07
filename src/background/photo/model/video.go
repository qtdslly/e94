package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Video struct {
	Id            uint32     `gorm:"primary_key" json:"id"`
	SubjectId     uint32     `json:"subject_id"`   //豆瓣subject ID
	Title         string     `gorm:"size:255" json:"title"`
	Description   string     `gorm:"type:longtext" json:"description"`
	Score         float64    `json:"score"`
	Year          uint32     `json:"year"`
	Actors        string     `gorm:"size:255" json:"actors"`
	Directors     string     `gorm:"size:255" json:"directors"`
	Writer        string     `gorm:"size:255" json:"writer"`
	Types         string     `gorm:"size:100" json:"types"`
	Country       string     `gorm:"size:100" json:"country"`
	Language      string     `gorm:"size:255" json:"language"`
	ReleaseDate   string     `gorm:"size:60" json:"release_date"`
	Duration      uint32     `json:"duration"`
	Alias         string     `gorm:"size:255" json:"alias"`
	Imdb          string     `gorm:"size:255" json:"imdb"`
	Comments      uint32     `json:"comments"`   //评论人数
	ThumbX        string     `gorm:"size:255" json:"thumb_x"`
	ThumbY        string     `gorm:"size:255" json:"thumb_y"`
	OfficialUrl   string     `gorm:"size:255" json:"offical_url"`
	TotalEpisode  uint32     `gorm:"size:255" json:"total_episode"`
	ContentType   uint32     `json:"content_type"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Video) TableName() string {
	return "video"
}

func initVideo(db *gorm.DB) error {
	var err error

	if db.HasTable(&Video{}) {
		err = db.AutoMigrate(&Video{}).Error
	} else {
		err = db.CreateTable(&Video{}).Error
	}
	return err
}

func dropVideo(db *gorm.DB) {
	db.DropTableIfExists(&Video{})
}
