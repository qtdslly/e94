package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"github.com/PuerkitoBio/goquery"
	uutil "background/newmovie/util"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	GetHanJuInfo("http://www.hanju.cc/hanju/151646.html")
}


func GetHanJuInfo(url string){
	document, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return
	}
	videoType := document.Find("#sdlist").Find(".sdlist").Eq(0).Find(".pleft").Eq(0).Find("a").Eq(2).Text()
	videoType,_ = uutil.DecodeToGBK(videoType)

	logger.Debug(videoType)

}
