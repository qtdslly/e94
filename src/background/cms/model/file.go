package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type File struct {
	Id          uint64           `gorm:"primary_key;auto_increment" json:"id"`
	RelPath     string           `gorm:"size:255" json:"rel_path"`
	Filename    string           `gorm:"size:128" json:"filename"`
	Groups      []*UserOpinion   `gorm:"many2many:user_thumb" json:"user_thumb"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (File) TableName() string {
	return "file"
}

func initFile(db *gorm.DB) error {
	var err error

	if db.HasTable(&File{}) {
		err = db.AutoMigrate(&File{}).Error
	} else {
		err = db.CreateTable(&File{}).Error
	}

	return err
}

func dropFile(db *gorm.DB) {
	db.DropTableIfExists(&File{})
}
