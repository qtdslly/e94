package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"

	"strings"
	"background/newmovie/util"
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
	file := flag.String("file", "list.txt", "file name")

	flag.Parse()

	if len(*file) == 0{
		logger.Debug("useage: ./addstreams -file list.txt")
		return
	}

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
	f, err := os.Open(*file)
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

		logger.Debug(line)

		var title string
		var url string
		var category string
		fields := strings.Split(line, "|")
		if len(fields) == 2 {
			title = fields[0]
			url = fields[1]
		} else if len(fields) == 3 {
			title = fields[0]
			url = fields[1]
			category = fields[2]
		}else {
			continue
		}

		util.StreamAdd(title,url,category,db)
	}


}


