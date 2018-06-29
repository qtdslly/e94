package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"

	"fmt"
	"strings"
	"background/common/constant"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
	"os"
	"bufio"
	"io"
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
	f, err := os.Open("/home/lyric/lly/dsp/tmp2/aa.txt")
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		var url string
		fields := strings.Split(line, "|")
		url = fields[0]

		var play model.PlayUrl
		play.Url = url
		play.Provider = uint32(constant.ContentProviderSystem)
		if err := db.Where("provider = ? and url = ?",play.Provider,play.Url).First(&play).Error ; err == gorm.ErrRecordNotFound{
			fmt.Println(url)
		}
	}


}


