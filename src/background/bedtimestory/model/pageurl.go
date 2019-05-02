package model



import (
	"time"

	"github.com/jinzhu/gorm"
)

type PageUrl struct {
	Id              uint32     `gorm:"primary_key" json:"id"`
	ProviderId      uint8      `gorm:"provider_id" json:"provider_id"`
	Url             string     `gorm:"url" json:"url"`
	PageStatus      uint8      `gorm:"page_status" json:"page_status"`
	UrlStatus       uint8      `gorm:"url_status" json:"url_status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (PageUrl) TableName() string {
	return "page_url"
}

func initPageUrl(db *gorm.DB) error {
	var err error

	if db.HasTable(&PageUrl{}) {
		err = db.AutoMigrate(&PageUrl{}).Error
	} else {
		err = db.CreateTable(&PageUrl{}).Error
	}
	return err
}

func dropPageUrl(db *gorm.DB) {
	db.DropTableIfExists(&PageUrl{})
}
