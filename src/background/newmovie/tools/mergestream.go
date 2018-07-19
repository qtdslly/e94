package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")

	streamId1 := flag.Int("a", 0, "first stream id")
	streamId2 := flag.Int("b", 0, "second stream id")

	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.LogMode(true)

	tx := db.Begin()
	if err := db.Exec("update play_url set content_id = ? where content_type = 4 and content_id = ?",streamId1,streamId2).Error ; err != nil{
		logger.Error(err)
		tx.Rollback()
		return
	}

	if err := db.Where("id = ?",streamId2).Delete(model.Stream{}).Error ; err != nil{
		logger.Error(err)
		tx.Rollback()
		return
	}
	tx.Commit()
}


