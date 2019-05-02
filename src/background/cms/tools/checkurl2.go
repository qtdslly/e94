package main

import (
	"fmt"

	"background/newmovie/config"
	"background/common/logger"
	"background/common/util"
	"os"
	"bufio"
	"io"
	"strings"
	//"background/newmovie/model"
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

	//f, err := os.Open("/home/lyric/Git/e94/src/background/newmovie/tools/2.txt")
	f, err := os.Open("C:/work/code/e94/src/background/newmovie/tools/2.txt")
	//f, err := os.Open("f:/Git/e94/src/background/newmovie/tools/2.txt")

	if err != nil {
		logger.Error(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	i := 0
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)

		logger.Debug(line)

		//var title string
		//var url string
		//fields := strings.Split(line, "|")
		//if len(fields) == 2{
		//	title = fields[0]
		//	url = fields[1]
		//}else{
		//	url = fields[0]
		//}

		url := line

		//var playUrl model.PlayUrl
		//if err = db.Where("url = ?",url).First(&playUrl).Error ; err == nil{
		//	continue
		//}
		if util.CheckStream(url,"c:/work/photo/red/" + fmt.Sprint(i) + ".jpg"){
			fmt.Println(fmt.Sprint(i) + "|SUCCESS|" + "" + "|" + url)
		}
		i++
	}
}


