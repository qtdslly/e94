package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
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
	model.InitModel(db)

	if err := db.Exec("delete from stream_group").Error ; err != nil{
		logger.Error(err)
		return
	}

	var resourceGroups []model.ResourceGroup
	if err := db.Where("type = 4").Find(&resourceGroups).Error ; err != nil{
		logger.Error(err)
		return
	}

	groupMap := make(map[string]uint32)

	for _,group := range resourceGroups{
		if strings.Contains(group.Name,"@"){
			group.Name = strings.Replace(group.Name,"@","QTDSLLY",-1)
		}
		groupMap[group.Name] = group.Id
	}

	var streams []model.Stream
	if err := db.Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,stream := range streams{
		if strings.Contains(stream.Category,"@"){
			stream.Category = strings.Replace(stream.Category,"@","QTDSLLY",-1)
		}
		groupId, ok := groupMap[stream.Category]
		if !ok{
			logger.Debug("NOT BIND : ",stream.Category,"|",stream.Id,"|" ,stream.Title)
			continue
		}
		if err := db.Exec("insert into stream_group(stream_id,resource_group_id) values(?,?)",stream.Id,groupId).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}


