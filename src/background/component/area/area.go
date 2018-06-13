package area

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Area struct {
	Code       string    `gorm:"primary_key;size:16" json:"code"`
	Country    string    `gorm:"size:16" json:"country"`
	Name       string    `gorm:"size:128" json:"name"`
	ParentCode string    `gorm:"size:16" json:"parent_code"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Area) TableName() string {
	return "area"
}

func InitArea(db *gorm.DB) error {
	var err error
	if db.HasTable(&Area{}) {
		err = db.AutoMigrate(&Area{}).Error
	} else {
		err = db.CreateTable(&Area{}).Error
	}
	return err
}

func DropArea(db *gorm.DB) {
	db.DropTableIfExists(&Area{})
}
