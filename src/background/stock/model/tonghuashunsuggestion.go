package model

import (
	"github.com/jinzhu/gorm"
)
type TonghuashunSuggestion struct {
	Id                 uint32     `gorm:"primary_key" json:"id"`
	Code               string     `gorm:"code" json:"code"`                            //股票代码
	TotalScore	   float32    `gorm:"total_score" json:"totalScore"`                //综合得分
	TotalAnalyseInfo   string     `gorm:"total_analyse_info" json:"totalAnalyseInfo"`   //买卖信号
	ClassNumber        float32    `gorm:"class_number" json:"classnumber"`              //
	Suggestion	   string     `gorm:"suggestion" json:"suggestion"`                 //建议
	TotalAnalyse	   string     `gorm:"total_analyse" json:"totalAnalyse"`            //综合分析
	StockName	   string     `gorm:"stock_name" json:"stockname"`                  //股票名称
	Date               string     `gorm:"date" json:"date"`                             //建议日期
}

func (TonghuashunSuggestion) TableName() string {
	return "tonghuashun_suggestion"
}

func initTonghuashunSuggestion(db *gorm.DB) error {
	var err error

	if db.HasTable(&TonghuashunSuggestion{}) {
		err = db.AutoMigrate(&TonghuashunSuggestion{}).Error
	} else {
		err = db.CreateTable(&TonghuashunSuggestion{}).Error
	}
	return err
}

func dropTonghuashunSuggestion(db *gorm.DB) {
	db.DropTableIfExists(&TonghuashunSuggestion{})
}
