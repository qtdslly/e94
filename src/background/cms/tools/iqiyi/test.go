package main

import (
	"fmt"
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

	url := "https://www.iqiyi.com/v_19rrnmdr2s.html"
	logger.Debug(url)
	realUrl := script.GetIqiyiRealPlayUrl(url)
	fmt.Println(realUrl)
}