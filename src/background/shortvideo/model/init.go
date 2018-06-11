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

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropTag(db)
	InitModel(db)
}
