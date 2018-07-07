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
		if stream.Title =="CCTV1"{
			stream.Sort = 1
		}else if stream.Title == "CCTV2"{
			stream.Sort = 2
		}else if stream.Title == "CCTV3"{
			stream.Sort = 3
		}else if stream.Title == "CCTV4"{
			stream.Sort = 4
		}else if stream.Title == "CCTV5"{
			stream.Sort = 5
		}else if stream.Title == "CCTV6"{
			stream.Sort = 6
		}else if stream.Title == "CCTV7"{
			stream.Sort = 7
		}else if stream.Title == "CCTV8"{
			stream.Sort = 8
		}else if stream.Title == "CCTV9"{
			stream.Sort = 9
		}else if stream.Title == "CCTV10"{
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
			stream.Sort = 20
		}else if strings.Contains(stream.Title,"CGTN") {
			stream.Sort = 30
		}else if strings.Contains(stream.Title,"CHC") {
			stream.Sort = 40
		}else if strings.Contains(stream.Title,"CIBN") {
			stream.Sort = 50
		}else if strings.Contains(stream.Title,"NewTV") || strings.Contains(stream.Title,"NEWTV") {
			stream.Sort = 60
		}else if strings.Contains(stream.Title,"中国黄河") || strings.Contains(stream.Title,"中国气象"){
			stream.Sort = 80
		}else if strings.Contains(stream.Category,"央视") {
			stream.Sort = 90
		}else if strings.Contains(stream.Title,"湖南卫视"){
			stream.Sort = 151
		}else if strings.Contains(stream.Title,"浙江卫视"){
			stream.Sort = 152
		}else if strings.Contains(stream.Title,"江苏卫视"){
			stream.Sort = 153
		}else if strings.Contains(stream.Title,"东方卫视"){
			stream.Sort = 154
		}else if strings.Contains(stream.Title,"北京卫视"){
			stream.Sort = 155
		}else if strings.Contains(stream.Title,"卫视") && !strings.Contains(stream.Category,"香港")&& !strings.Contains(stream.Category,"澳门")&& !strings.Contains(stream.Category,"台湾"){
			stream.Sort = 200
		}else if strings.Contains(stream.Category,"北京"){
			if strings.Contains(stream.Title,"北京"){
				stream.Sort = 300
			}else{
				stream.Sort = 320
			}
		}else if strings.Contains(stream.Category,"上海"){
			if strings.Contains(stream.Title,"东方"){
				stream.Sort = 400
			}else if strings.Contains(stream.Title,"上海"){
				stream.Sort = 420
			}else{
				stream.Sort = 440
			}
		}else if strings.Contains(stream.Category,"江苏"){
			if strings.Contains(stream.Title,"江苏"){
				stream.Sort = 500
			}else if strings.Contains(stream.Title,"南京"){
				stream.Sort = 520
			}else  if strings.Contains(stream.Title,"苏州"){
				stream.Sort = 540
			}else  if strings.Contains(stream.Title,"扬州"){
				stream.Sort = 560
			}else  if strings.Contains(stream.Title,"无锡"){
				stream.Sort = 580
			}else  if strings.Contains(stream.Title,"常州"){
				stream.Sort = 600
			}else  if strings.Contains(stream.Title,"盐城"){
				stream.Sort = 620
			}else{
				stream.Sort = 700
			}
		}else if strings.Contains(stream.Category,"浙江"){
			if strings.Contains(stream.Title,"浙江"){
				stream.Sort = 800
			}else if strings.Contains(stream.Title,"杭州"){
				stream.Sort = 820
			}else  if strings.Contains(stream.Title,"宁波"){
				stream.Sort = 840
			}else  if strings.Contains(stream.Title,"温州"){
				stream.Sort = 860
			}else  if strings.Contains(stream.Title,"绍兴"){
				stream.Sort = 880
			}else  if strings.Contains(stream.Title,"舟山"){
				stream.Sort = 900
			}else  if strings.Contains(stream.Title,"丽水"){
				stream.Sort = 920
			}else{
				stream.Sort = 1000
			}
		}else if strings.Contains(stream.Category,"四川"){
			if strings.Contains(stream.Title,"四川"){
				stream.Sort = 1100
			}else if strings.Contains(stream.Title,"成都"){
				stream.Sort = 1120
			}else  if strings.Contains(stream.Title,"乐山"){
				stream.Sort = 1140
			}else  if strings.Contains(stream.Title,"眉山"){
				stream.Sort = 1160
			}else  if strings.Contains(stream.Title,"南充"){
				stream.Sort = 1180
			}else  if strings.Contains(stream.Title,"广安"){
				stream.Sort = 1200
			}else  if strings.Contains(stream.Title,"德阳"){
				stream.Sort = 1220
			}else{
				stream.Sort = 1300
			}
		}else if strings.Contains(stream.Category,"湖北"){
			if strings.Contains(stream.Title,"湖北"){
				stream.Sort = 1400
			}else if strings.Contains(stream.Title,"武汉"){
				stream.Sort = 1420
			}else  if strings.Contains(stream.Title,"宜昌"){
				stream.Sort = 1440
			}else  if strings.Contains(stream.Title,"襄阳"){
				stream.Sort = 1460
			}else  if strings.Contains(stream.Title,"十堰"){
				stream.Sort = 1480
			}else  if strings.Contains(stream.Title,"黄石"){
				stream.Sort = 1500
			}else  if strings.Contains(stream.Title,"黄冈"){
				stream.Sort = 1520
			}else  if strings.Contains(stream.Title,"孝感"){
				stream.Sort = 1540
			}else  if strings.Contains(stream.Title,"鄂州"){
				stream.Sort = 1560
			}else  if strings.Contains(stream.Title,"咸宁"){
				stream.Sort = 1580
			}else{
				stream.Sort = 1600
			}
		}else if strings.Contains(stream.Category,"江西"){
			if strings.Contains(stream.Title,"江西"){
				stream.Sort = 1700
			}else if strings.Contains(stream.Title,"南昌"){
				stream.Sort = 1720
			}else{
				stream.Sort = 1740
			}
		}else if strings.Contains(stream.Category,"天津"){
			if strings.Contains(stream.Title,"天津"){
				stream.Sort = 2000
			}else{
				stream.Sort = 2100
			}
		}else if strings.Contains(stream.Category,"湖南"){
			if strings.Contains(stream.Title,"湖南"){
				stream.Sort = 2300
			}else if strings.Contains(stream.Title,"长沙"){
				stream.Sort = 2320
			}else{
				stream.Sort = 2400
			}
		}else if strings.Contains(stream.Category,"河南"){
			if strings.Contains(stream.Title,"河南"){
				stream.Sort = 2600
			}else if strings.Contains(stream.Title,"郑州"){
				stream.Sort = 2620
			}else{
				stream.Sort = 2700
			}
		}else if strings.Contains(stream.Category,"广东"){
			if strings.Contains(stream.Title,"广东"){
				stream.Sort = 3000
			}else if strings.Contains(stream.Title,"广州"){
				stream.Sort = 3020
			}else if strings.Contains(stream.Title,"深圳"){
				stream.Sort = 3040
			}else{
				stream.Sort = 3100
			}
		}else if strings.Contains(stream.Category,"安徽"){
			if strings.Contains(stream.Title,"安徽"){
				stream.Sort = 3500
			}else if strings.Contains(stream.Title,"合肥"){
				stream.Sort = 3520
			}else{
				stream.Sort = 3600
			}
		}else if strings.Contains(stream.Category,"福建"){
			if strings.Contains(stream.Title,"福建"){
				stream.Sort = 3800
			}else if strings.Contains(stream.Title,"福州"){
				stream.Sort = 3820
			}else if strings.Contains(stream.Title,"厦门"){
				stream.Sort = 3840
			}else if strings.Contains(stream.Title,"泉州"){
				stream.Sort = 3860
			}else{
				stream.Sort = 3900
			}
		}else if strings.Contains(stream.Category,"河北"){
			if strings.Contains(stream.Title,"河北"){
				stream.Sort = 4100
			}else if strings.Contains(stream.Title,"石家庄"){
				stream.Sort = 4120
			}else{
				stream.Sort = 4200
			}
		}else if strings.Contains(stream.Category,"山东"){
			if strings.Contains(stream.Title,"山东"){
				stream.Sort = 4700
			}else if strings.Contains(stream.Title,"济南"){
				stream.Sort = 4720
			}else{
				stream.Sort = 4800
			}
		}else if strings.Contains(stream.Category,"重庆"){
			if strings.Contains(stream.Title,"重庆"){
				stream.Sort = 5000
			}else{
				stream.Sort = 5100
			}
		}else if strings.Contains(stream.Category,"陕西"){
			if strings.Contains(stream.Title,"陕西"){
				stream.Sort = 5300
			}else if strings.Contains(stream.Title,"西安"){
				stream.Sort = 5320
			}else{
				stream.Sort = 5400
			}
		}else if strings.Contains(stream.Category,"新疆"){
			if strings.Contains(stream.Title,"新疆"){
				stream.Sort = 5600
			}else if strings.Contains(stream.Title,"乌鲁木齐"){
				stream.Sort = 5620
			}else{
				stream.Sort = 5700
			}
		}else if strings.Contains(stream.Category,"内蒙古"){
			if strings.Contains(stream.Title,"新疆"){
				stream.Sort = 5900
			}else if strings.Contains(stream.Title,"乌鲁木齐"){
				stream.Sort = 5920
			}else{
				stream.Sort = 6000
			}
		}else if strings.Contains(stream.Category,"辽宁"){
			if strings.Contains(stream.Title,"辽宁"){
				stream.Sort = 6200
			}else if strings.Contains(stream.Title,"沈阳"){
				stream.Sort = 6220
			}else if strings.Contains(stream.Title,"大连"){
				stream.Sort = 6240
			}else{
				stream.Sort = 6300
			}
		}else if strings.Contains(stream.Category,"吉林"){
			if strings.Contains(stream.Title,"吉林"){
				stream.Sort = 6500
			}else if strings.Contains(stream.Title,"长春"){
				stream.Sort = 6520
			}else{
				stream.Sort = 6600
			}
		}else if strings.Contains(stream.Category,"黑龙江"){
			if strings.Contains(stream.Title,"黑龙江"){
				stream.Sort = 6800
			}else if strings.Contains(stream.Title,"哈尔滨"){
				stream.Sort = 6820
			}else{
				stream.Sort = 6900
			}
		}else if strings.Contains(stream.Category,"甘肃"){
			if strings.Contains(stream.Title,"甘肃"){
				stream.Sort = 7100
			}else if strings.Contains(stream.Title,"兰州"){
				stream.Sort = 7120
			}else{
				stream.Sort = 7200
			}
		}else if strings.Contains(stream.Category,"宁夏"){
			if strings.Contains(stream.Title,"宁夏"){
				stream.Sort = 7400
			}else if strings.Contains(stream.Title,"银川"){
				stream.Sort = 7420
			}else{
				stream.Sort = 7500
			}
		}else if strings.Contains(stream.Category,"贵州"){
			if strings.Contains(stream.Title,"贵州"){
				stream.Sort = 7700
			}else if strings.Contains(stream.Title,"贵阳"){
				stream.Sort = 7720
			}else{
				stream.Sort = 7800
			}
		}else if strings.Contains(stream.Category,"青海"){
			if strings.Contains(stream.Title,"青海"){
				stream.Sort = 8000
			}else if strings.Contains(stream.Title,"西宁"){
				stream.Sort = 8020
			}else{
				stream.Sort = 8100
			}
		}else if strings.Contains(stream.Category,"云南"){
			if strings.Contains(stream.Title,"云南"){
				stream.Sort = 8300
			}else if strings.Contains(stream.Title,"昆明"){
				stream.Sort = 8320
			}else{
				stream.Sort = 8400
			}
		}else if strings.Contains(stream.Category,"广西"){
			if strings.Contains(stream.Title,"广西"){
				stream.Sort = 8600
			}else if strings.Contains(stream.Title,"昆明"){
				stream.Sort = 8320
			}else{
				stream.Sort = 8400
			}
		}else if strings.Contains(stream.Category,"海南"){
			stream.Sort = 8900
		}else if strings.Contains(stream.Category,"山西"){
			stream.Sort = 9200
		}else if strings.Contains(stream.Category,"西藏"){
			stream.Sort = 9500
		}else if strings.Contains(stream.Category,"特色") {
			stream.Sort = 10000
		}else if strings.Contains(stream.Category,"明星") {
			stream.Sort = 10500
		}else if strings.Contains(stream.Category,"风景"){
			stream.Sort = 11000
		}else if strings.Contains(stream.Category,"香港"){
			stream.Sort = 13000
		}else if strings.Contains(stream.Category,"澳门"){
			stream.Sort = 13300
		}else if strings.Contains(stream.Category,"台湾"){
			stream.Sort = 13600
		}else if strings.Contains(stream.Category,"韩国"){
			stream.Sort = 14600
		}else if strings.Contains(stream.Category,"朝鲜"){
			stream.Sort = 15000
		}else if strings.Contains(stream.Category,"日本"){
			stream.Sort = 16000
		}else if strings.Contains(stream.Category,"新加坡"){
			stream.Sort = 17000
		}else if strings.Contains(stream.Category,"越南"){
			stream.Sort = 17500
		}else if strings.Contains(stream.Category,"泰国"){
			stream.Sort = 18000
		}else if strings.Contains(stream.Category,"马来西亚"){
			stream.Sort = 19000
		}else if strings.Contains(stream.Category,"美国"){
			stream.Sort = 20000
		}else if strings.Contains(stream.Category,"德国") {
			stream.Sort = 21000
		}else if strings.Contains(stream.Category,"法国"){
			stream.Sort = 21500
		}else if strings.Contains(stream.Category,"加拿大"){
			stream.Sort = 22000
		}else if strings.Contains(stream.Category,"俄罗斯"){
			stream.Sort = 23000
		}else if strings.Contains(stream.Category,"澳大利亚"){
			stream.Sort = 24000
		}else if strings.Contains(stream.Category,"沙特阿拉伯"){
			stream.Sort = 25000
		}else if strings.Contains(stream.Category,"印度"){
			stream.Sort = 26000
		}else if strings.Contains(stream.Category,"英国"){
			stream.Sort = 27000
		}else if strings.Contains(stream.Category,"卡塔尔"){
			stream.Sort = 28000
		}else if strings.Contains(stream.Category,"西班牙"){
			stream.Sort = 29000
		}else if strings.Contains(stream.Category,"希腊"){
			stream.Sort = 30000
		}else if strings.Contains(stream.Category,"欧洲"){
			stream.Sort = 31000
		}else if strings.Contains(stream.Category,"其他外国频道"){
			stream.Sort = 33000
		}else{
			logger.Debug("=================================",stream.Title)
			stream.Sort = 0
		}
		if err := db.Save(&stream).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}


