package model



import (
	"time"

	"github.com/jinzhu/gorm"
)

type StatusConfig struct {
	Id        uint32     `gorm:"primary_key" json:"id"`
	Key       string     `gorm:"key" json:"key"`
	Value     string     `gorm:"value" json:"value"`
	Other     string     `gorm:"other" json:"other"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (StatusConfig) TableName() string {
	return "status_config"
}

func initStatusConfig(db *gorm.DB) error {
	var err error

	if db.HasTable(&StatusConfig{}) {
		err = db.AutoMigrate(&StatusConfig{}).Error
	} else {
		err = db.CreateTable(&StatusConfig{}).Error
	}
	return err
}

func dropStatusConfig(db *gorm.DB) {
	db.DropTableIfExists(&StatusConfig{})
}
