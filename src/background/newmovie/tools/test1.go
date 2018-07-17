package main

import (
	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"background/newmovie/config"

	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	apiUrl := "http://jx.618g.com/?url=" + "http://v.youku.com/v_show/id_XMzcxMjM4Mjg0MA==.html?s=2eaf9082aa3811e68fae"

	query, err := goquery.NewDocument(apiUrl)
	if err != nil {
		logger.Debug(apiUrl)
		logger.Error(err)
		return
	}

	base,exist := query.Find("iframe").Eq(0).Attr("src")
	if !exist{
		logger.Debug(apiUrl)
		logger.Error(err)
		return
	}
	url := base[strings.Index(base,"url=") + 4:]
	logger.Debug(url)
}
