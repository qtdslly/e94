package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Property struct {
	Id        uint32    `gorm:"primary_key" json:"id"`
	Type      uint32    `json:"type"`
	Name      string    `gorm:"size:64;unique" json:"name" valid:"Str" name:"name" len:"1,64"`
	Sort      uint32    `json:"sort"`
	Tags      []*Tag    `json:"tags"`
	Persons   []*Person `json:"persons"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	PropertyTypeUnknown = 0
	PropertyTypeNormal  = 1
	PropertyTypePerson  = 2
)

type PropertySlice []*Property

func (p PropertySlice) Len() int           { return len(p) }
func (p PropertySlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PropertySlice) Less(i, j int) bool { return p[i].Sort < p[j].Sort }

func (Property) TableName() string {
	return "property"
}

func initProperty(db *gorm.DB) error {
	var err error

	if db.HasTable(&Property{}) {
		err = db.AutoMigrate(&Property{}).Error
	} else {
		err = db.CreateTable(&Property{}).Error
	}
	return err
}

func dropProperty(db *gorm.DB) {
	db.DropTableIfExists(&Property{})
}
