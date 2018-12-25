package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	Id           uint32         `gorm:"primary_key" json:"id"`
	DoctorId     uint32         `json:"doctor_id"`
	UserId       uint32         `json:"user_id"`
	Content      string         `gorm:"size:255" json:"content"`
	Reply        string         `gorm:"size:255" json:"reply"`
	State        int            `json:"state"`

	CreatedAt    time.Time      `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt    time.Time      `json:"updated_at"`       // 更新时间，utc格式
}

const (
	CommentStatePublish = 1
	CommentStateUnknow  = 2
	CommentStateUnPass  = 3
)
func (Comment) TableName() string {
	return "comment"
}

func initComment(db *gorm.DB) error {
	var err error
	if db.HasTable(&Comment{}) {
		err = db.AutoMigrate(&Comment{}).Error
	} else {
		err = db.CreateTable(&Comment{}).Error
	}
	return err
}

func dropComment(db *gorm.DB) {
	db.DropTableIfExists(&Comment{})
}
