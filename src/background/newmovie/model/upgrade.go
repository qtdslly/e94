package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Upgrade struct {
	Id               uint32    `gorm:"primary_key" json:"id"`
	Title            string    `gorm:"size:128" json:"title" valid:"Str" name:"title" len:"1,128"`                 // 升级配置名称
	UpgradeVersion   string    `gorm:"size:32" json:"min_version" valid:"Str" name:"min_version" len:"0,32"`       // 升级版本号
	TargetVersion    string    `gorm:"size:32" json:"target_version" valid:"Str" name:"target_version" len:"0,32"` // 目标版本号
	Enable           bool      `json:"enable"`                                                                     // 是否启用
	ShowUpgrade      bool      `json:"show_upgrade"`                                                               // 提示更新
	ForceUpgrade     bool      `json:"force_upgrade"`                                                              // 强制更新
	CheckUpgrade     bool      `json:"check_upgrade"`                                                              // 检测更新
	UpgradeTip       string    `gorm:"size:2048" json:"upgrade_tip" valid:"Str" name:"upgrade_tip" len:"0,2048"`   // 更新提示
	UpgradeUrl       string    `gorm:"size:255" json:"upgrade_url" valid:"Str" name:"upgrade_url" len:"0,255"`     // 更新文件下载地址
	CreatedAt        time.Time `json:"created_at"`                                                                 // 创建时间，utc格式
	UpdatedAt        time.Time `json:"updated_at"`                                                                 // 更新时间，utc格式
}

func (Upgrade) TableName() string {
	return "upgrade"
}

func initUpgrade(db *gorm.DB) error {
	var err error
	if db.HasTable(&Upgrade{}) {
		err = db.AutoMigrate(&Upgrade{}).Error
	} else {
		err = db.CreateTable(&Upgrade{}).Error
	}
	return err
}

func dropUpgrade(db *gorm.DB) {
	db.DropTableIfExists(&Upgrade{})
}