package model

import (
	"background/common/logger"
	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error

	err = initStock(db)
	if err != nil {
		logger.Fatal("Init db stock failed, ", err)
		return err
	}

	err = initShareHolder(db)
	if err != nil {
		logger.Fatal("Init db share_holder failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	dropStock(db)
	dropShareHolder(db)

	InitModel(db)
}
