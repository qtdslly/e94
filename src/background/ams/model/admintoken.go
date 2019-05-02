package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AdminToken struct {
	Id        uint32     `gorm:"primary_key" json:"id"`
	AdminId   uint32     `gorm:"index" json:"admin_id"`
	Disabled  bool       `json:"disabled"`
	Token     string     `gorm:"size:100" json:"token"`
	UserAgent string     `gorm:"size:255" json:"user_agent"`
	LoginAt   *time.Time `json:"login_at"`
	LoginIp   string     `gorm:"size:96" json:"login_ip"`
	CreatedIp string     `gorm:"size:96" json:"created_ip"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (AdminToken) TableName() string {
	return "admin_token"
}

func initAdminToken(db *gorm.DB) error {
	var err error

	if db.HasTable(&AdminToken{}) {
		err = db.AutoMigrate(&AdminToken{}).Error
	} else {
		err = db.CreateTable(&AdminToken{}).Error
	}
	return err
}

func dropAdminToken(db *gorm.DB) {
	db.DropTableIfExists(&AdminToken{})
}
