package main

import (
	//"background/stock/model"
	"background/common/logger"
	"background/stock/config"

	"flag"
	//"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"background/stock/model"
	"net/http"
	"fmt"
	"time"
	"os"
	"io"
)

func main(){
	var err error
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err = config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Config Failed!!!!", err)
		return
	}

	logger.SetLevel(config.GetLoggerLevel())

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)

	model.InitModel(db)

	var stocks []model.StockBasic
	if err = db.Find(&stocks).Error ; err != nil{
		logger.Error(err)
		return
	}

	p := time.Now()
	end :=  fmt.Sprintf("%04d%02d%02d",p.Year(),p.Month(),p.Day())
	for _,stock := range stocks{
		url := "http://quotes.money.163.com/service/chddata.html?code=1" + stock.Code
		start := stock.ToMarketDate
		if len(start) != 8{
			logger.Debug("股票代码:",stock.Code," 股票名称:",stock.Name," 数据获取失败")
			continue
		}
		url = url + "&start=" + start + "&end=" + end
		url = url + "&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"
		logger.Debug(url)
		resp, err := http.Get(url)
		if err != nil {
			logger.Error(err)
			return
		}
		defer resp.Body.Close()
		f, err := os.Create("/root/data/wangyi/" + stock.Code + ".csv")
		if err != nil {
			panic(err)
		}
		io.Copy(f, resp.Body)
	}


}


