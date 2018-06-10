package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Installation struct {
	Id             uint64    `gorm:"primary_key" json:"id"`
	OsVersion      string    `gorm:"size:64" json:"os_version"`
	DeviceModel    string    `gorm:"size:64" json:"device_model"`
	DeviceId       string    `gorm:"size:64;index" json:"device_id"`
	MacAddress     string    `gorm:"size:32;index" json:"mac_address"`
	Imei           string    `gorm:"size:32;index" json:"imei"`
	Carrier        uint8     `json:"carrier" valid:"Carrier" name:"carrier"`
	//Channel        string    `gorm:"size:32;index" json:"channel"`
	Product        string    `gorm:"size:32" json:"product"`  //产品名称
	Brand          string    `gorm:"size:32" json:"brand"`   //设备品牌
	//Longitude      float32   `json:"longitude"` // 用户的经度
	//Latitude       float32   `json:"latitude"`  // 用户的纬度
	ActiveIp       string    `gorm:"size:32" json:"active_ip"`
	CreatedAt      time.Time `json:"created_at"` // 创建时间，utc格式
	UpdatedAt      time.Time `json:"updated_at"` // 更新时间，utc格式
}

/*
	global carrier type definition.
*/
const (
	CarrierTypeUnknown      = 0 // 未知类型
	CarrierTypeChinaMobile  = 1 // 中国移动
	CarrierTypeChinaUnicom  = 2 // 中国联通
	CarrierTypeChinaTelecom = 3 // 中国电信

	PushTypeunknown = 0
	PushTypeGetui   = 1
	PushTypeXiaomi  = 2
	PushTypeAPN     = 3
)

func (Installation) TableName() string {
	return "installation"
}

func initInstallation(db *gorm.DB) error {
	var err error
	if db.HasTable(&Installation{}) {
		err = db.AutoMigrate(&Installation{}).Error
	} else {
		err = db.CreateTable(&Installation{}).Error
	}
	return err
}

func dropInstallation(db *gorm.DB) {
	db.DropTableIfExists(&Installation{})
}
