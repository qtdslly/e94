package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initUser(db)
	if err != nil {
		logger.Fatal("Init db user failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropUser(db)

	InitModel(db)
}
