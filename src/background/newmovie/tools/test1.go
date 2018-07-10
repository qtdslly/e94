package main

import (
	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"background/newmovie/config"

	"strings"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	apiUrl := "http://jx.618g.com/?url=" + "http://www.iqiyi.com/v_19rrbxpprk.html"

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
