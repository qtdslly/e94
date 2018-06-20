package main

import (
	"flag"
	"background/pic/config"
	"background/pic/model"
	"background/common/logger"

	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
)
func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")

	flag.Parse()

	logger.SetLevel(config.GetLoggerLevel())

	err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Config Failed!!!!", err)
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)
	model.InitModel(db)
	logger.SetLevel(config.GetLoggerLevel())

	url := "https://page.appdao.com/forward?link=16195339&style=160105&item=883716&page=1&limit=25&after=57210369&screen_w=1242&screen_h=2208&ir=0&app=1P_ElfWallpapers&v=1.4&lang=zh-Hans-CN&it=1529458959.610182&ots=3&jb=0&as=0&mobclix=0&deviceid=replaceudid&macaddr=&idv=4ADB9AA0-A063-492D-8FD6-632E835F7AF4&idvs=&ida=F7835587-9CBD-4B92-8C6E-2813811E7B5F&phonetype=iphone&model=iphone8%2C2&osn=iOS&osv=11.3.1&tz=8"
	GetBizhi(url,db)

}


func GetBizhi(url string,db *gorm.DB)bool{
	requ, err := http.NewRequest("GET", url, nil)
	requ.Header.Add("Host", "page.appdao.com")
	requ.Header.Add("User-Agent", "LiVideoIOS/4.3.6 (iPhone; iOS 11.3.1; Scale/3.00)")

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		logger.Debug("Proxy failed!")
		return false
	}

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return false
	}

	logger.Debug(string(recv))

	return true
}