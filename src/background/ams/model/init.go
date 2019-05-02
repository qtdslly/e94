package model

import (
	"github.com/jinzhu/gorm"
	"background/common/logger"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initAdmin(db)
	if err != nil {
		logger.Fatal("Init db Admin failed, ", err)
		return err
	}

	err = initAdminToken(db)
	if err != nil {
		logger.Fatal("Init db Admin failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropAdmin(db)
	dropAdminToken(db)

	InitModel(db)
}
