package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initPicture(db)
	if err != nil {
		logger.Fatal("Init db picture failed, ", err)
		return err
	}

	err = initMove(db)
	if err != nil {
		logger.Fatal("Init db move failed, ", err)
		return err
	}

	err = initTag(db)
	if err != nil {
		logger.Fatal("Init db tag failed, ", err)
		return err
	}

	err = initCategory(db)
	if err != nil {
		logger.Fatal("Init db category failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropPicture(db)
	dropMove(db)
	dropTag(db)
	dropCategory(db)
	InitModel(db)
}
