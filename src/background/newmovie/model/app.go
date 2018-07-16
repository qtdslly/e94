package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type App struct {
	Id           uint32     `gorm:"primary_key" json:"id"`
	Name         string     `gorm:"size:64" json:"name" valid:"Str" name:"name" len:"1,64"`
	AppKey       string     `gorm:"size:64;" json:"app_key"`
	Type         uint32     `json:"type" valid:"AppType" name:"type"`
	OsType       uint32     `json:"os_type" valid:"OsType" name:"os_type"`
	Disabled     bool       `json:"disabled"`
	Versions     []*Version `json:"versions"`
	IosCount     uint32     `gorm:"-" json:"ios_count"`     //IOS版本数量
	AndroidCount uint32     `gorm:"-" json:"android_count"` //安卓版本数量
	Sort         uint32     `json:"sort"`
	Creator      string     `gorm:"size:255" json:"creator"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

const (
	AppTypeApp = 1 // app
	AppTypeSdk = 2 // sdk
)

func (App) TableName() string {
	return "app"
}

func initApp(db *gorm.DB) error {
	var err error
	if db.HasTable(&App{}) {
		err = db.AutoMigrate(&App{}).Error
	} else {
		err = db.CreateTable(&App{}).Error
	}

	return err
}

func dropApp(db *gorm.DB) {
	db.DropTableIfExists(&App{})
}
