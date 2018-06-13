package country

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Country struct {
	Code        string    `gorm:"primary_key;size:16" json:"code"`
	Name        string    `gorm:"size:50" json:"name"`
	DialingCode string    `gorm:"size:16" json:"dialing_code"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Country) TableName() string {
	return "country"
}

func InitCountry(db *gorm.DB) error {
	var err error
	if db.HasTable(&Country{}) {
		err = db.AutoMigrate(&Country{}).Error
	} else {
		err = db.CreateTable(&Country{}).Error
	}
	return err
}

func DropCountry(db *gorm.DB) {
	db.DropTableIfExists(&Country{})
}
