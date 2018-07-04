package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	Id                  uint32      `gorm:"primary_key" json:"id"`
	Nickname            string      `gorm:"size:30;index" json:"nickname"`
	Avatar              string      `json:"avatar"`
	Gender              uint8       `json:"gender"`
	InstallationId      uint64      `json:"installation_id"`
	Bean                uint32      `json:"bean"`
	Birthday            string      `grom:"size:10" json:"birthday"`
	CheckinDays         uint32      `json:"checkin_days"` // 连续签到天数
	LastCheckin         time.Time   `json:"last_checkin"` // 最后签到时间
	Status              uint32      `json:"status"`
	Token               string      `gorm:"size:100" json:"token"`
	Laravel             uint32      `json:"laravel"`   // 用户等级
	Longitude           float32     `json:"longitude"` // 经度
	Latitude            float32     `json:"latitude"`  // 纬度
	LastUseIp           string      `gorm:"size:96" json:"last_use_ip"`
	LastUseAt           time.Time   `json:"last_use_at"` // 最后使用时间
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

const (
	UserStatusLoginBanned   = 0 // 用户状态：禁止登录
	UserStatusUseBanned     = 1 // 用户状态：禁止使用app
	UserStatusWhiteList     = 2 // 用户状态：白名单

	UserOrdinary = 0    //普通用户   可看所有国内频道
	UserPublic   = 1    //大众会员   可看澳门频道，可看最新电影
	UserBronze   = 2    //青铜会员   可看香港频道
	UserSilver   = 3    //白银会员   可看台湾频道
	UserGold     = 4    //黄金会员   可看海外频道
	UserPlatinum = 5    //铂金会员   可看韩国视频
	UserDiamonds = 6    //钻石会员   可看福利视频

)

func (User) TableName() string {
	return "user"
}

func initUser(db *gorm.DB) error {
	var err error

	if db.HasTable(&User{}) {
		err = db.AutoMigrate(&User{}).Error
	} else {
		err = db.CreateTable(&User{}).Error
	}
	return err
}

func dropUser(db *gorm.DB) {
	db.DropTableIfExists(&User{})
}
