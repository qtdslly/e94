package model



import (
	"time"

	"github.com/jinzhu/gorm"
)

type Content struct {
	Id              uint32     `gorm:"primary_key" json:"id"`
	PageId          uint32     `gorm:"provider_id" json:"provider_id"`
	Title           string     `gorm:"type:title" json:"title"`
	Content         string     `gorm:"type:content" json:"content"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Content) TableName() string {
	return "content"
}

func initContent(db *gorm.DB) error {
	var err error

	if db.HasTable(&Content{}) {
		err = db.AutoMigrate(&Content{}).Error
	} else {
		err = db.CreateTable(&Content{}).Error
	}
	return err
}

func dropContent(db *gorm.DB) {
	db.DropTableIfExists(&Content{})
}
