package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initPageUrl(db)
	if err != nil {
		logger.Fatal("Init db page_url failed, ", err)
		return err
	}

	err = initProvider(db)
	if err != nil {
		logger.Fatal("Init db provider failed, ", err)
		return err
	}

	err = initContent(db)
	if err != nil {
		logger.Fatal("Init db content failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropPageUrl(db)
	dropProvider(db)
	dropContent(db)
	InitModel(db)
}
