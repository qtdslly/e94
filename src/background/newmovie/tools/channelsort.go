package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	"strings"
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

	var streams []model.Stream
	if err := db.Find(&streams).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _,stream := range streams{
		if strings.Contains(stream.Title,"CCTV1"){
			stream.Sort = 1
		}else if strings.Contains(stream.Title,"CCTV2"){
			stream.Sort = 2
		}else if strings.Contains(stream.Title,"CCTV3"){
			stream.Sort = 3
		}else if strings.Contains(stream.Title,"CCTV4"){
			stream.Sort = 4
		}else if strings.Contains(stream.Title,"CCTV5"){
			stream.Sort = 5
		}else if strings.Contains(stream.Title,"CCTV6"){
			stream.Sort = 6
		}else if strings.Contains(stream.Title,"CCTV7"){
			stream.Sort = 7
		}else if strings.Contains(stream.Title,"CCTV8"){
			stream.Sort = 8
		}else if strings.Contains(stream.Title,"CCTV9"){
			stream.Sort = 9
		}else if strings.Contains(stream.Title,"CCTV10"){
			stream.Sort = 10
		}else if strings.Contains(stream.Title,"CCTV11"){
			stream.Sort = 11
		}else if strings.Contains(stream.Title,"CCTV12"){
			stream.Sort = 12
		}else if strings.Contains(stream.Title,"CCTV13"){
			stream.Sort = 13
		}else if strings.Contains(stream.Title,"CCTV14"){
			stream.Sort = 14
		}else if strings.Contains(stream.Title,"CCTV15"){
			stream.Sort = 15
		}else if strings.Contains(stream.Title,"CCTV"){
			stream.Sort = 16
		}else if strings.Contains(stream.Title,"中国教育"){
			stream.Sort = 17
		}else if strings.Contains(stream.Title,"NewTV") {
			stream.Sort = 18
		}else if strings.Contains(stream.Title,"中国黄河") || strings.Contains(stream.Title,"中国气象"){
			stream.Sort = 20
		}else if strings.Contains(stream.Title,"湖南卫视"){
			stream.Sort = 31
		}else if strings.Contains(stream.Title,"浙江卫视"){
			stream.Sort = 32
		}else if strings.Contains(stream.Title,"江苏卫视"){
			stream.Sort = 33
		}else if strings.Contains(stream.Title,"东方卫视"){
			stream.Sort = 34
		}else if strings.Contains(stream.Title,"北京卫视"){
			stream.Sort = 35
		}else if strings.Contains(stream.Title,"卫视"){
			stream.Sort = 41
		}else if strings.Contains(stream.Category,"北京"){
			stream.Sort = 51
		}else if strings.Contains(stream.Category,"上海"){
			stream.Sort = 52
		}else if strings.Contains(stream.Category,"江苏"){
			stream.Sort = 53
		}else if strings.Contains(stream.Category,"浙江"){
			stream.Sort = 54
		}else if strings.Contains(stream.Category,"四川"){
			stream.Sort = 55
		}else if strings.Contains(stream.Category,"湖北"){
			stream.Sort = 56
		}else if strings.Contains(stream.Category,"江西"){
			stream.Sort = 57
		}else if strings.Contains(stream.Category,"天津"){
			stream.Sort = 58
		}else if strings.Contains(stream.Category,"湖南"){
			stream.Sort = 59
		}else if strings.Contains(stream.Category,"广东"){
			stream.Sort = 60
		}else if strings.Contains(stream.Category,"安徽"){
			stream.Sort = 61
		}else if strings.Contains(stream.Category,"福建"){
			stream.Sort = 62
		}else if strings.Contains(stream.Category,"河北"){
			stream.Sort = 63
		}else if strings.Contains(stream.Category,"山东"){
			stream.Sort = 64
		}else if strings.Contains(stream.Category,"山西"){
			stream.Sort = 65
		}else if strings.Contains(stream.Category,"陕西"){
			stream.Sort = 66
		}else if strings.Contains(stream.Category,"新疆"){
			stream.Sort = 67
		}else if strings.Contains(stream.Category,"内蒙古"){
			stream.Sort = 68
		}else if strings.Contains(stream.Category,"辽宁"){
			stream.Sort = 69
		}else if strings.Contains(stream.Category,"吉林"){
			stream.Sort = 70
		}else if strings.Contains(stream.Category,"黑龙江"){
			stream.Sort = 71
		}else if strings.Contains(stream.Category,"甘肃"){
			stream.Sort = 72
		}else if strings.Contains(stream.Category,"宁夏"){
			stream.Sort = 73
		}else if strings.Contains(stream.Category,"贵州"){
			stream.Sort = 74
		}else if strings.Contains(stream.Category,"云南"){
			stream.Sort = 75
		}else if strings.Contains(stream.Category,"广西"){
			stream.Sort = 76
		}else if strings.Contains(stream.Category,"西藏"){
			stream.Sort = 77
		}else if strings.Contains(stream.Category,"风景"){
			stream.Sort = 81
		}else if strings.Contains(stream.Category,"香港"){
			stream.Sort = 91
		}else if strings.Contains(stream.Category,"澳门"){
			stream.Sort = 111
		}else if strings.Contains(stream.Category,"台湾"){
			stream.Sort = 131
		}else if strings.Contains(stream.Category,"韩国"){
			stream.Sort = 161
		}else if strings.Contains(stream.Category,"朝鲜"){
			stream.Sort = 231
		}else if strings.Contains(stream.Category,"日本"){
			stream.Sort = 251
		}else if strings.Contains(stream.Category,"新加坡"){
			stream.Sort = 301
		}else if strings.Contains(stream.Category,"泰国"){
			stream.Sort = 331
		}else if strings.Contains(stream.Category,"马来西亚"){
			stream.Sort = 361
		}else if strings.Contains(stream.Category,"美国"){
			stream.Sort = 401
		}else if strings.Contains(stream.Category,"加拿大"){
			stream.Sort = 481
		}else if strings.Contains(stream.Category,"俄罗斯"){
			stream.Sort = 521
		}else if strings.Contains(stream.Category,"澳大利亚"){
			stream.Sort = 601
		}else if strings.Contains(stream.Category,"沙特阿拉伯"){
			stream.Sort = 651
		}else if strings.Contains(stream.Category,"印度"){
			stream.Sort = 701
		}else if strings.Contains(stream.Category,"英国"){
			stream.Sort = 751
		}else if strings.Contains(stream.Category,"卡塔尔"){
			stream.Sort = 801
		}else if strings.Contains(stream.Category,"西班牙"){
			stream.Sort = 821
		}else if strings.Contains(stream.Category,"欧洲"){
			stream.Sort = 851
		}else if strings.Contains(stream.Category,"其他外国频道"){
			stream.Sort = 1001
		}
		if err := db.Save(&stream).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}


