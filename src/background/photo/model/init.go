package model

import (
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error


	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	InitModel(db)
}
