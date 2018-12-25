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

	err = initComment(db)
	if err != nil {
		logger.Fatal("Init db comment failed, ", err)
		return err
	}

	err = initDoctor(db)
	if err != nil {
		logger.Fatal("Init db doctor failed, ", err)
		return err
	}

	err = initDuty(db)
	if err != nil {
		logger.Fatal("Init db duty failed, ", err)
		return err
	}
	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropUser(db)
	dropComment(db)
	dropDoctor(db)
	dropDuty(db)

	InitModel(db)
}
