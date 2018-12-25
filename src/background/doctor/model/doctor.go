package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Doctor struct {
	Id           uint32         `gorm:"primary_key" json:"id"`
	Name         string         `json:"name"`
	Avtar        string         `json:"avtar"`
	Hospital     string         `json:"hospital"`    // 医院
	Department   string         `json:"department"`  // 科室
	StartWork    int            `json:"start_work"`  // 开始工作时间
	Position     string         `json:"position"`    // 职位
	Title        string         `json:"Title"`       // 职称
	Nick         string         `json:"nick"`        // 昵称

	CreatedAt    time.Time      `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt    time.Time      `json:"updated_at"`       // 更新时间，utc格式
}


func (Doctor) TableName() string {
	return "doctor"
}

func initDoctor(db *gorm.DB) error {
	var err error
	if db.HasTable(&Doctor{}) {
		err = db.AutoMigrate(&Doctor{}).Error
	} else {
		err = db.CreateTable(&Doctor{}).Error
	}
	return err
}

func dropDoctor(db *gorm.DB) {
	db.DropTableIfExists(&Doctor{})
}
