package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type StreamSource struct {
	Id        uint32    `gorm:"primary_key" json:"id"`
	StreamId  uint32    `json:"stream_id"`
	Title     string    `gorm:"size:255" json:"title"`
	Ratio     string    `json:"ratio"`
	Width     uint32    `json:"width"`
	Height    uint32    `json:"height"`
	Bitrate   uint32    `json:"bitrate"`
	Url       string    `gorm:"size:255" json:"url"`
	P2PUrl    string    `gorm:"column:p2p_url" json:"p2p_url"`
	NewP2PUrl string    `gorm:"column:new_p2p_url" json:"new_p2p_url"`
	Cp        string    `gorm:"size:64" json:"cp"` // content_provider id
	Quality   uint8     `json:"quality"`           // refer to constant/typ.go VideoQuality*
	CreatedAt time.Time `json:"created_at"`        // 创建时间，utc格式
	UpdatedAt time.Time `json:"updated_at"`        // 更新时间，utc格式
}

type Stream struct {
	Id             uint32           `gorm:"primary_key" json:"id"`
	Title          string           `gorm:"size:60" json:"title" valid:"Str" name:"title" len:"1,255"`
	Area           string           `gorm:"size:16" json:"area"`
	Pinyin         string           `gorm:"size:32;index" json:"pinyin"`
	Description    string           `gorm:"type:longtext" json:"description"`
	Sort           uint32           `json:"sort"`
	Icon           string           `gorm:"size:255" json:"icon" valid:"Str" name:"icon" len:"0,255"`
	Thumb          string           `gorm:"size:255" json:"thumb" valid:"Str" name:"thumb" len:"0,255"`
	PlayUrls       []*PlayUrl       `json:"play_urls"`
	OnLine         bool             `json:"on_line"`
	Disable        bool             `json:"disable"`
	HasEpg         bool             `json:"has_epg"`
	Category       string           `json:"category"`
	EpgSyncedAt    *time.Time       `json:"epg_synced_at"`
	CreatedAt      time.Time        `json:"created_at"`       // 创建时间，utc格式
	UpdatedAt      time.Time        `json:"updated_at"`       // 更新时间，utc格式
}

func (StreamSource) TableName() string {
	return "stream_source"
}

func (Stream) TableName() string {
	return "stream"
}

func initStream(db *gorm.DB) error {
	var err error
	if db.HasTable(&Stream{}) {
		err = db.AutoMigrate(&Stream{}).Error
	} else {
		err = db.CreateTable(&Stream{}).Error
		if err == nil {
			err = db.Exec("alter table stream add unique (title) ;").Error
		}
	}
	return err
}

func dropStream(db *gorm.DB) {
	db.DropTableIfExists(&Stream{})
}
