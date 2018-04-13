package task

import (
	"background/stock/model"
	"background/common/logger"
	apiths "background/stock/api/tonghuashun"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)
func GetTonghuashun(db *gorm.DB){
	var err error
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
}
