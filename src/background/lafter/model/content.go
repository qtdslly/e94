package model



import (
	"time"

	"github.com/jinzhu/gorm"
)

type Content struct {
	Id              uint32     `gorm:"primary_key" json:"id"`
	PageId          uint32     `gorm:"page_id" json:"page_id"`
	Title           string     `gorm:"title" json:"title"`
	Content         string     `gorm:"type:longtext" json:"content"`
	State           uint32     `gorm:"state" json:"state"`
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
