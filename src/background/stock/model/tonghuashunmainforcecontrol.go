package model

import (
	"github.com/jinzhu/gorm"
)
type TonghuashunMainForceControl struct {
	Id            uint32     `gorm:"primary_key" json:"id"`
	Code	      string     `gorm:"code" json:"code"`                    //股票名称
	FundAnalyse   string     `gorm:"fund_analyse" json:"fundAnalyse"`     //控盘信息
	CurrentFund   string     `gorm:"current_fund" json:"currentFund"`     //当前净流入
	State         string     `gorm:"state" json:"state"`                  //流入流出状态
	Amount        float32    `gorm:"amount" json:"amount"`                //流入流出金额
	ControlValue  string     `gorm:"control_value" json:"controlvalue"`  //控盘度
	Date          string     `gorm:"date" json:"date"`                    //日期
}

func (TonghuashunMainForceControl) TableName() string {
	return "tonghuashun_main_force_control"
}

func initTonghuashunMainForceControl(db *gorm.DB) error {
	var err error

	if db.HasTable(&TonghuashunMainForceControl{}) {
		err = db.AutoMigrate(&TonghuashunMainForceControl{}).Error
	} else {
		err = db.CreateTable(&TonghuashunMainForceControl{}).Error
	}
	return err
}

func dropTonghuashunMainForceControl(db *gorm.DB) {
	db.DropTableIfExists(&TonghuashunMainForceControl{})
}
