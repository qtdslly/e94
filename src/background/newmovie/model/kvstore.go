package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type KvStore struct {
	Id        uint32     `gorm:"primary_key" json:"id"`
	Key  string     `gorm:"size:64" json:"key"`
	Value    string    `gorm:"type:longtext" json:"value"`
	AppId     uint32    `json:"app_id"`
	VersionId     uint32    `json:"version_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
	}

	return err
}

func dropKvStore(db *gorm.DB) {
	db.DropTableIfExists(&KvStore{})
}
