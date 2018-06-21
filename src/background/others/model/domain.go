package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Domain struct {
	Id                 uint32     `gorm:"primary_key" json:"id"`
	Url                string     `gorm:"size:255" json:"url"`
	Status             uint32     `json:"status"`
	ExpirationDate     string     `gorm:"size:10" json:"expiration_ate"`
	RegisterDate       string     `gorm:"size:10" json:"register_ate"`
	Reseller           string     `gorm:"size:255" json:"reseller"`
	RegistrantName     string     `gorm:"size:255" json:"registrant_name"`
	Sponsoring         string     `gorm:"size:255" json:"sponsoring"`
	RegistrantCountry  string     `gorm:"size:255" json:"registrant_country"`
	RegistrantProvince string     `gorm:"size:255" json:"registrant_province"`
	RegistrantCity     string     `gorm:"size:255" json:"registrant_city"`
	RegistrantStreet   string     `gorm:"size:255" json:"registrant_street"`
	RegistrantEmail    string     `gorm:"size:100" json:"registrant_email"`
	RegistrantPhone    string     `gorm:"size:100" json:"registrant_phone"`
	RegistrarUrl       string     `gorm:"size:255" json:"registrar_url"`

	Chinese            string     `gorm:"size:100" json:"chinese"`

	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Domain) TableName() string {
	return "domain"
}

func initDomain(db *gorm.DB) error {
	var err error

	if db.HasTable(&Domain{}) {
		err = db.AutoMigrate(&Domain{}).Error
	} else {
		err = db.CreateTable(&Domain{}).Error
	}
	return err
}

func dropDomain(db *gorm.DB) {
	db.DropTableIfExists(&Domain{})
}
