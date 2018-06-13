package model

import (
	"background/common/logger"

	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = InitThirdVideo(db)
	if err != nil {
		logger.Fatal("Init db kv failed, ", err)
		return err
	}

	err = initPerson(db)
	if err != nil {
		logger.Fatal("Init db person failed, ", err)
		return err
	}

	err = initTag(db)
	if err != nil {
		logger.Fatal("Init db tag failed, ", err)
		return err
	}

	err = initProperty(db)
	if err != nil {
		logger.Fatal("Init db property failed, ", err)
		return err
	}

	err = initVideo(db)
	if err != nil {
		logger.Fatal("Init db video failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropTag(db)
	InitModel(db)
}
