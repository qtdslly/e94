package model

import (
	"background/common/logger"
	"background/common/util"
	"time"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	Id        uint32     `gorm:"primary_key" json:"id"`
	Username  string     `gorm:"size:64;unique" json:"username"`
	Mobile    string     `gorm:"size:16;unique" json:"mobile"`
	Email     string     `gorm:"size:64;unique" json:"email"`
	Password  string     `gorm:"size:255" json:"password"`
	LoginAt   *time.Time `json:"login_at"`
	Locked    bool       `json:"locked"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Admin) TableName() string {
	return "admin"
}

func initAdmin(db *gorm.DB) error {
	var err error
	if db.HasTable(&Admin{}) {
		err = db.AutoMigrate(&Admin{}).Error
	} else {
		err = db.CreateTable(&Admin{}).Error
	}

	var admin Admin
	admin.Username = "lly"
	if err = db.Where("username=?", admin.Username).First(&admin).Error; err != nil {
		logger.Debug(err)
		admin.Password = util.SHA512("abc123")
		db.Create(&admin)
	}

	return err
}

func dropAdmin(db *gorm.DB) {
	db.DropTableIfExists(&Admin{})
}
