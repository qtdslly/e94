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

  var stocks []model.StockBasic
  if err := db.Order("code asc").Find(&stocks).Error; err != nil {
    logger.Error(err)
    return
  }

  p := time.Now()
  date := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
  for _ , stock := range stocks{
    var cyfx model.TonghuashunCyfx
    cyfx.Code = stock.Code
    cyfx.Date = date
    update := false
    if err := db.Where("date = ? and code = ?",cyfx.Date,cyfx.Code).First(&cyfx).Error ; err == nil{
      update = true
    }
    content := cc.Cyfx(stock.Code)
    if content == ""{
      continue
    }
    cyfx.Content = content

    if update{
      if err := db.Save(&cyfx).Error ; err != nil{
        logger.Error(err)
        return
      }
    }else{
      if err := db.Create(&cyfx).Error ; err != nil{
        logger.Error(err)
        return
      }
    }
  }
}































