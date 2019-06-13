package task

import (
  "fmt"
  "time"

  "background/common/logger"
  cc "background/stock/common"

  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
  "background/stock/model"
)

func Zncp(db *gorm.DB)error{
  logger.Debug("开始抓取百度股票数据")
  p := time.Now()

  dd, _ := time.ParseDuration("24h")
  to := p.Add(dd)
  tomorry := fmt.Sprintf("%04d-%02d-%02d",to.Year(),to.Month(),to.Day())

  var task model.StockTask
  if err := db.Where("`key` = 'baidu_zncp'").First(&task).Error ; err != nil{
    logger.Error(err)
    return err
  }

  today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
  if task.Date > today{
    logger.Debug("今日数据已抓取")
    return nil
  }

  var stocks []model.StockBasic
  if err := db.Order("code asc").Where("date = ?",today).Find(&stocks).Error ; err != nil{
    logger.Error(err)
    return err
  }

  for _ ,stock := range stocks{
    jys := cc.GetJysCodeByStockCode(stock.Code)
    cc.Zncp(db,jys,stock.Code)
  }

  if err := db.Model(model.StockTask{}).Where("`key` = 'baidu_zncp'").Update("date", tomorry).Error; err != nil {
    logger.Error(err)
    return err
  }

  return nil

}
