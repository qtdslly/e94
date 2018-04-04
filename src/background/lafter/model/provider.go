package model



import (
	"time"

	"github.com/jinzhu/gorm"
)

type Provider struct {
	Id          uint8     `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"name" json:"name"`
	Url         string     `gorm:"url" json:"url"`
	Status      string     `gorm:"status" json:"status"`  // 0不爬取 1爬取
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Provider) TableName() string {
	return "provider"
}

func initProvider(db *gorm.DB) error {
	var err error

	if db.HasTable(&PageUrl{}) {
		err = db.AutoMigrate(&Provider{}).Error
	} else {
		err = db.CreateTable(&Provider{}).Error
	}
	return err
}

func dropProvider(db *gorm.DB) {
	db.DropTableIfExists(&Provider{})
}
