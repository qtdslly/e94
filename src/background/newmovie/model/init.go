package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initMovie(db)
	if err != nil {
		logger.Fatal("Init db movie failed, ", err)
		return err
	}
	err = initTopSearch(db)
	if err != nil {
		logger.Fatal("Init db top_search failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropMovie(db)
	dropTopSearch(db)

	InitModel(db)
}
