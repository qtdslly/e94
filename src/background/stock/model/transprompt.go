package model

import (
	"github.com/jinzhu/gorm"
)

type TransPrompt struct {
	Id                     uint32     `gorm:"primary_key" json:"id"`
	StockCode              string     `gorm:"stock_code" json:"stock_code"`                            //股票代码
	AvgPrice               float64    `gorm:"avg_price" json:"avg_price"`
	LastTransPrice         float64    `gorm:"last_trans_price" json:"last_trans_price"`
	LastTransDate          string     `gorm:"last_trans_date" json:"last_trans_date"`
	LastTransCount         int64      `gorm:"last_trans_count" json:"last_trans_count"`
	LastTransType          string     `gorm:"last_trans_type" json:"last_trans_type"`
	PromptBuyPrice         float64    `gorm:"prompt_buy_price" json:"prompt_buy_price"`
	PromptBuyCount         int64      `gorm:"prompt_buy_count" json:"prompt_buy_count"`
	PromptSellPrice        float64    `gorm:"prompt_sell_price" json:"prompt_sell_price"`
	PromptSellCount        int64      `gorm:"prompt_sell_count" json:"prompt_sell_count"`
}

func (TransPrompt) TableName() string {
	return "trans_prompt"
}

func initTransPrompt(db *gorm.DB) error {
	var err error

	if db.HasTable(&TransPrompt{}) {
		err = db.AutoMigrate(&TransPrompt{}).Error
	} else {
		err = db.CreateTable(&TransPrompt{}).Error
	}
	return err
}

func dropTransPrompt(db *gorm.DB) {
	db.DropTableIfExists(&TransPrompt{})
}
