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
		groupMap[group.Name] = group.Id
	}

	var streams []model.Stream
	if err := db.Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,stream := range streams{
		if strings.Contains(stream.Title,"体育") || strings.Contains(stream.Title,"运动") ||
			strings.Contains(stream.Title,"足球") || strings.Contains(stream.Title,"篮球") ||
			strings.Contains(stream.Title,"CCTV5") || strings.Contains(stream.Title,"乒羽") ||
			strings.Contains(stream.Title,"垂钓") || strings.Contains(stream.Title,"钓鱼") ||
			strings.Contains(stream.Title,"网球") || strings.Contains(stream.Title,"台球"){
			if err := db.Exec("insert into stream_group(stream_id,resource_group_id) values(?,?)",stream.Id,groupMap["体育"]).Error ; err != nil{
				logger.Error(err)
				return
			}
			continue
		}

		if strings.Contains(stream.Title,"电影") || strings.Contains(stream.Title,"影视") {
			if err := db.Exec("insert into stream_group(stream_id,resource_group_id) values(?,?)",stream.Id,groupMap["电影"]).Error ; err != nil{
				logger.Error(err)
				return
			}
		}

		if stream.Category == "体育" || stream.Category == "电影"{
			continue
		}

		if stream.Category == "墨西哥" || stream.Category == "捷克" ||
			stream.Category == "希腊" || stream.Category == "意大利" ||
			stream.Category == "马来西亚" || stream.Category == "加拿大" ||
			stream.Category == "澳大利亚" || stream.Category == "印度" ||
			stream.Category == "日本" || stream.Category == "卡塔尔" ||
			stream.Category == "英国" || stream.Category == "新西兰" ||
			stream.Category == "印度尼西亚"{
			stream.Category = "其他外国频道"
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

		if stream.Category == "卫视"{
			stream.Title = strings.Replace(stream.Title,"卫视","",-1)
			stream.Title = strings.Replace(stream.Title,"藏语","",-1)

			if stream.Title == "农林"{
				stream.Title = "甘肃"
			}else if stream.Title == "南方"{
				stream.Title = "广东"
			}else if stream.Title == "东方"{
				stream.Title = "上海"
			}else if stream.Title == "东南" ||stream.Title == "厦门" || stream.Title == "海峡"{
				stream.Title = "福建"
			}else if stream.Title == "兵团"{
				stream.Title = "新疆"
			}else if stream.Title == "安多"{
				stream.Title = "青海"
			}else if stream.Title == "延边"{
				stream.Title = "吉林"
			}else if stream.Title == "旅游"{
				stream.Title = "海南"
			}else if stream.Title == "山东教育"{
				stream.Title = "山东"
			}else if stream.Title == "康巴"{
				stream.Title = "四川"
			}

			if err := db.Exec("insert into stream_group(stream_id,resource_group_id) values(?,?)",stream.Id,groupMap[stream.Title]).Error ; err != nil{
				logger.Error(err)
				return
			}
		}

	}
}


