package model



import (
	"time"

	"github.com/jinzhu/gorm"
)

type PhoneAddress struct {
	Id        uint32     `gorm:"primary_key" json:"id"`
	Mts       string     `gorm:"mts" json:"mts"`
	Province  string     `gorm:"province" json:"province"`
	CatName   string     `gorm:"cat_name" json:"cat_name"`
	AreaVid   string     `json:"area_vid" json:"area_vid"`
	IspVid    string     `gorm:"isp_vid" json:"isp_vid"`
	Carrier   string     `gorm:"carrier" json:"carrier"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (PhoneAddress) TableName() string {
	return "phone_address"
}

func initPhoneAddress(db *gorm.DB) error {
	var err error

	if db.HasTable(&PhoneAddress{}) {
		err = db.AutoMigrate(&PhoneAddress{}).Error
	} else {
		err = db.CreateTable(&PhoneAddress{}).Error
	}
	return err
}

func dropPhoneAddress(db *gorm.DB) {
	db.DropTableIfExists(&PhoneAddress{})
}
