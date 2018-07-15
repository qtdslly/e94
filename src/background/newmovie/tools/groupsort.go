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

	var resourceGroups []model.ResourceGroup
	if err := db.Where("content_type = 4").Find(&resourceGroups).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,group := range resourceGroups{
		var count uint32
		if err = db.Table("stream").Joins("inner join stream_group where stream.on_line = 1 and stream.id = stream_group.stream_id and stream_group.resource_group_id = ?",group.Id).Count(&count).Error ; err != nil{
			logger.Error(err)
			return
		}
		group.Count = count
		group.Sort = group.Count + 100
		if group.Name == "央视"{
			group.Sort = 1000
		}else if group.Name == "卫视"{
			group.Sort = 999
		}else if group.Name == "NewTV"{
			group.Sort = 998
		}else if group.Name == "台湾"{
			group.Sort = 997
		}else if group.Name == "香港"{
			group.Sort = 996
		}else if group.Name == "澳门"{
			group.Sort = 995
		}
		if group.Count == 0{
			group.OnLine = false
		}
		if err := db.Save(group).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}


