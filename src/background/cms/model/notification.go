package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Notification struct {
	Id            uint32     `gorm:"primary_key" json:"id"`
	Title         string     `gorm:"size:255" json:"title"`
	Thumb         string     `gorm:"size:255" json:"thumb"`
	Description   string     `gorm:"size:255" json:"description"`
	ContentType   uint32     `json:"content_type"`
	ContentId     uint32     `json:"content_id"`
	OnLine        bool       `json:"on_line"`                         //是否已推送
	ExpireTime    int64      `json:"expire_time"`                    //消息离线存储有效期，单位：ms
	Clicks        uint32     `json:"clicks"`                         //点击数
	Impressions   uint32     `json:"impressions"`                    //展示数
	CreatedAt     time.Time  `json:"created_at"`                     // 创建时间，utc格式
	UpdatedAt     time.Time  `json:"updated_at"`                     // 更新时间，utc格式
}

func (Notification) TableName() string {
	return "notification"
}

func initNotification(db *gorm.DB) error {
	var err error
	if db.HasTable(&Notification{}) {
		err = db.AutoMigrate(&Notification{}).Error
	} else {
		err = db.CreateTable(&Notification{}).Error
	}
	return err
}

func dropNotification(db *gorm.DB) {
	db.DropTableIfExists(&Notification{})
}
