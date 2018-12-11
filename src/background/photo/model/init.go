package model

import (
	"github.com/jinzhu/gorm"
	"background/qrcode/logger"
)

func InitModel(db *gorm.DB) error {
	var err error

	if err = initPhoto(db) ; err != nil{
		logger.Error(err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropPhoto(db)

	InitModel(db)
}
