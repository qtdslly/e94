package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initPage(db)
	if err != nil {
		logger.Fatal("Init db page failed, ", err)
		return err
	}

	err = initVideo(db)
	if err != nil {
		logger.Fatal("Init db movie failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropVideo(db)
	dropPage(db)
	InitModel(db)
}
