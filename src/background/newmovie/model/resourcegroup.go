package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ResourceGroup struct {
	Id              uint32           `gorm:"primary_key" json:"id"`
	Type            uint32           `json:"type"` //对应media类型 global media content type
	Name            string           `gorm:"size:64" json:"name" valid:"Str" name:"name" len:"1,64"`
	Icon            string           `gorm:"size:255" json:"icon" valid:"Str" name:"icon" len:"0,255"`
	Sort            uint32           `json:"sort"`
	Count           uint32           `json:"count"`
	Videos          []*Video         `gorm:"many2many:video_group" json:"videos"`
	Streams         []*Stream        `gorm:"many2many:stream_group" json:"streams"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

const (
	ThirdResourceGroup = "第三方资源组"
)

type ResourceGroupSlice []*ResourceGroup

func (c ResourceGroupSlice) Len() int           { return len(c) }
func (c ResourceGroupSlice) Less(i, j int) bool { return c[i].Id < c[j].Id }
func (c ResourceGroupSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (ResourceGroup) TableName() string {
	return "resource_group"
}

func initResourceGroup(db *gorm.DB) error {
	var err error

	if db.HasTable(&ResourceGroup{}) {
		err = db.AutoMigrate(&ResourceGroup{}).Error
	} else {
		err = db.CreateTable(&ResourceGroup{}).Error
	}
	return err
}

func dropResourceGroup(db *gorm.DB) {
	db.DropTableIfExists(&ResourceGroup{})
}
