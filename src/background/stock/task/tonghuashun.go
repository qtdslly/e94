package task

import (
  "time"
  "fmt"

	"background/stock/model"
	"background/common/logger"
	apiths "background/stock/api/tonghuashun"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)
func GetTonghuashun(db *gorm.DB){
  logger.Debug("开始抓取同花顺股票数据")

  var err error

	var p = time.Now()
  today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())

  dd, _ := time.ParseDuration("24h")
  to := p.Add(dd)
  tomorry := fmt.Sprintf("%04d-%02d-%02d",to.Year(),to.Month(),to.Day())

  var stocks []model.StockBasic
  if err = db.Order("code asc").Find(&stocks).Error ; err != nil{
    logger.Error(err)
    return
  }

	var task model.StockTask
	if err := db.Where("`key` = 'tonghuashun_control'").First(&task).Error ; err != nil{
		logger.Error(err)
		return
	}

	if task.Date > today{
		logger.Debug("今日数据已抓取")
		return
	}else{
    for _,stock := range stocks{
      apiths.GetControlInfo(stock.Code,db)
    }

    if err := db.Model(model.StockTask{}).Where("`key` = 'tonghuashun_control'").Update("date", tomorry).Error; err != nil {
      logger.Error(err)
      return
    }
  }

  var task1 model.StockTask
  if err := db.Where("`key` = 'tonghuashun_comprehensive'").First(&task1).Error ; err != nil{
    logger.Error(err)
    return
  }

  if task1.Date > today{
    logger.Debug("今日数据已抓取")
    return
  }else{
    for _,stock := range stocks{
      apiths.GetComprehensive(stock.Code,db)
    }
    if err := db.Model(model.StockTask{}).Where("`key` = 'tonghuashun_comprehensive'").Update("date", tomorry).Error; err != nil {
      logger.Error(err)
      return
    }
  }



}
