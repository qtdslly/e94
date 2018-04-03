package model

import (
	"github.com/jinzhu/gorm"
)

/*
模拟信息
*/
type Simulation struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	Value                  string     `gorm:"value" json:"value"`
}

func (Simulation) TableName() string {
	return "simulation"
}

func initSimulation(db *gorm.DB) error {
	var err error

	if db.HasTable(&Simulation{}) {
		err = db.AutoMigrate(&Simulation{}).Error
	} else {
		err = db.CreateTable(&Simulation{}).Error
	}
	return err
}

func dropSimulation(db *gorm.DB) {
	db.DropTableIfExists(&Simulation{})
}
