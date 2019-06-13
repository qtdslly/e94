package main

import (
  "flag"

  "background/common/logger"
  "background/stock/config"
  cc "background/stock/common"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jinzhu/gorm"
  "background/stock/model"
  "time"
  "fmt"
)

func main() {
  var err error
  configPath := flag.String("conf", "../config/config.json", "Config file path")
  flag.Parse()

  err = config.LoadConfig(*configPath)
  if err != nil {
    logger.Error("Config Failed!!!!", err)
    return
  }

  logger.SetLevel(config.GetLoggerLevel())

  db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
  if err != nil {
    logger.Fatal("Open db Failed!!!!", err)
    return
  }

  db.LogMode(true)

  model.InitModel(db)

  p := time.Now()
  today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
  var stocks []model.WangYiStock
  if err = db.Order("code asc").Where("date = ?",today).Find(&stocks).Error; err != nil {
    logger.Error(err)
    return
  }

  for _,stock := range stocks{
    cc.GetWangyiStockBasic(stock.Jys,stock.Code,db)
  }

}































