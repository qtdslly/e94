package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ContentAction struct {
	Id           uint32    `gorm:"primary_key" json:"id"`
	InstallationId       uint64    `json:"installation_id"`
	Action       uint8    `json:"action"`
	ContentType  uint8     `json:"content_type"`
	ContentId    uint32    `json:"content_id"`
	Title        string    `json:"title"`
	Thumb        string    `json:"thumb"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

const (
	ContentActionDigg   = 1 //点赞
	ContentActionFollow = 2 //关注
)

func (ContentAction) TableName() string {
	return "content_action"
}

func initContentAction(db *gorm.DB) error {
	var err error

	if db.HasTable(&ContentAction{}) {
		err = db.AutoMigrate(&ContentAction{}).Error
	} else {
		err = db.CreateTable(&ContentAction{}).Error
	}

	return err
}

func dropContentAction(db *gorm.DB) {
	db.DropTableIfExists(&ContentAction{})
}
