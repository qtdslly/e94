package main

import (
	"background/newmovie/service/script"
	"background/newmovie/config"
	"background/common/logger"
	"flag"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	url := "http://www.rrmeiju.com/bofang/475/1/1.html"
	logger.Debug(url)
	realUrl := script.GetRRmeijuRealPlayUrl(url)
	logger.Debug(realUrl)
}