package kv

import (
	"time"

	"github.com/jinzhu/gorm"
)

type KvStore struct {
	AppId     uint32 `gorm:"index"`
	VersionId uint32 `gorm:"index"`
	Key       string `gorm:"index;size:64"`
	Value     string `gorm:"type:longtext"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (KvStore) TableName() string {
	return "kv_store"
}

func initKvStore(db *gorm.DB) error {
	var err error

	if db.HasTable(&KvStore{}) {
		err = db.AutoMigrate(&KvStore{}).Error
	} else {
		err = db.CreateTable(&KvStore{}).Error
		err = db.Exec("alter table kv_store add primary key(`app_id`, `version_id`, `key`);").Error
	}
	return err
}

func DropKvStore(db *gorm.DB) {
	db.DropTableIfExists(&KvStore{})
}
