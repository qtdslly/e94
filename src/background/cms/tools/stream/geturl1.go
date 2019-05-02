package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"

	"strings"
	"fmt"
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
	db.LogMode(false)
	model.InitModel(db)
	f, err := os.Open("/home/lyric/Git/e94/src/background/newmovie/tools/succ.txt")
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

		//logger.Debug(line)


		url := line

		var play model.PlayUrl
		play.Url = url
		if err := db.Where("url = ?",play.Url).First(&play).Error ; err == gorm.ErrRecordNotFound{
			fmt.Println(url)
		}
	}


}


