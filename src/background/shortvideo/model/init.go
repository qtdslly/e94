package model

import (
	"common/logger"
	"component/kv"

	"github.com/jinzhu/gorm"
)

func InitModel(db *gorm.DB) error {
	var err error



	err = kv.InitKv(db)
	if err != nil {
		logger.Fatal("Init db kv failed, ", err)
		return err
	}

	return err
}

// Do not call this method!!!!
func rebuildModel(db *gorm.DB) {
	kv.DropKvStore(db)
	InitModel(db)
}
