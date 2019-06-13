package main

import (
  "flag"

  "background/common/logger"
  "background/stock/config"
  cc "background/stock/common"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jinzhu/gorm"
  "background/stock/model"
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

  var stocks []model.WangYiStock
  if err = db.Find(&stocks).Error; err != nil {
    logger.Error(err)
    return
  }

  cc.Sstd(db)

}

