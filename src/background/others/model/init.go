package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initStatusConfig(db)
	if err != nil {
		logger.Fatal("Init db status_config failed, ", err)
		return err
	}

	err = initPhoneAddress(db)
	if err != nil {
		logger.Fatal("Init db phone_address failed, ", err)
		return err
	}

	err = initDomain(db)
	if err != nil {
		logger.Fatal("Init db domain failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropPhoneAddress(db)
	dropStatusConfig(db)
	dropDomain(db)
	InitModel(db)
}
