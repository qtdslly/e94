package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Person struct {
	Id           uint32         `gorm:"primary_key" json:"id"`
	Provider     uint32         `json:"provider"`
	SourceId     string         `gorm:"size:64" json:"source_id"`
	Name         string         `gorm:"size:64" json:"name"`
	LocalName    string         `gorm:"size:64" json:"local_name"`
	Nickname     string         `gorm:"size:128" json:"nickname"`
	Figure       string         `gorm:"size:256" json:"figure"`
	Description  string         `gorm:"type:longtext" json:"description"`
	Role         uint8          `json:"role"`
	Gender       uint8          `json:"gender"` // refer to constant/typ.go Gender*
	Country      string         `gorm:"size:8" json:"country"`
	VideoCount   uint32         `gorm:"-" json:"video_count"`
	Birthday     *time.Time     `json:"birthday"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	SyncedAt     *time.Time     `json:"synced_at"` // 同步时间，utc格式
}


const (
	PersonRoleTypeUper     = 1
	PersonRoleTypeActor    = 2
	PersonRoleTypeWriter   = 4
	PersonRoleTypeArtist   = 8
)

func (Person) TableName() string {
	return "person"
}

func initPerson(db *gorm.DB) error {
	var err error

	if db.HasTable(&Person{}) {
		err = db.AutoMigrate(&Person{}).Error
	} else {
		err = db.CreateTable(&Person{}).Error
	}

	return err
}

func dropPerson(db *gorm.DB) {
	db.DropTableIfExists(&Person{})
}
