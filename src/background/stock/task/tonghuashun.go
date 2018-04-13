package task

import (
	"fmt"

	"background/stock/model"
	"background/common/logger"
	apiths "background/stock/api/tonghuashun"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
func GetTonghuashun(db *gorm.DB){
	var err error

	var p = time.Now()
	var task model.StockTask
	if err := db.Where("`key` = 'tonghuashun'").First(&task).Error ; err != nil{
		logger.Error(err)
		return
	}

	today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
	if task.Date > today{
		logger.Debug("今日数据已抓取")
		return
	}

	var stocks []model.StockList
	if err = db.Find(&stocks).Error ; err != nil{
		logger.Error(err)
		return
	}
	for _,stock := range stocks{
		apiths.GetComprehensive(stock.Code,db)
	}

	for _,stock := range stocks{
		if apiths.GetControlInfo(stock.Code,db) != nil{
			continue
		}
	}

	now := time.Now()
	dd, _ := time.ParseDuration("24h")
	to := now.Add(dd)
	tomorry := fmt.Sprintf("%04d-%02d-%02d",to.Year(),to.Month(),to.Day())
	if err := db.Model(model.StockTask{}).Where("`key` = 'realtimestock'").Update("date", tomorry).Error; err != nil {
		logger.Error(err)
		return
	}
}
