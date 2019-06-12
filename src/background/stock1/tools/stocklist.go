package main

import (

	"log"
	"fmt"
	"flag"
	"io/ioutil"

	"background/stock1/model"
	"background/common/logger"
	"background/stock1/config"
	"background/common/systemcall"
	"background/common/util"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
  "time"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	systemcall.SetFileLimit()

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.DB().SetMaxIdleConns(10)

	db.LogMode(true)

	model.InitModel(db)

	GetStockList(db)
}

func GetStockList(db *gorm.DB)(error){

  p := time.Now()
  date := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
	var page int64 = 0
	for{
		url := "http://quotes.money.163.com/hs/service/diyrank.php?host=http://quotes.money.163.com/hs/service/diyrank.php&page=" + fmt.Sprintf("%d",page) + "&query=STYPE:EQA&fields=NO,SYMBOL,NAME,PRICE,PERCENT,UPDOWN,FIVE_MINUTE,OPEN,YESTCLOSE,HIGH,LOW,VOLUME,TURNOVER,HS,LB,WB,ZF,PE,MCAP,TCAP,MFSUM,MFRATIO.MFRATIO2,MFRATIO.MFRATIO10,SNAME,CODE,ANNOUNMT,UVSNEWS&sort=PERCENT&order=desc&count=24&type=query"
		resp, err := req.Get(url)
		if err != nil {
			logger.Error(err)
			return err
		}

		recv,err := ioutil.ReadAll(resp.Response().Body)
		if err != nil{
			logger.Error(err)
			return err
		}

		data,_ := util.DecodeToGBK(string(recv))

		logger.Debug(data)

		pageCount := gjson.Get(data, "pagecount").Int()

		list := gjson.Get(data, "list")

		if list.Exists() {
			items := list.Array()

			for _, item := range items {
				var stock model.Stock
				code := item.Get("SYMBOL").String()
				if err := db.Where("code = ?",code).First(&stock).Error ; err != nil{
          logger.Error(err)
          return err
				}

				stock.Code = item.Get("SYMBOL").String()
				stock.Jys = item.Get("CODE").String()[0:1]
				stock.FiveMinute = item.Get("FIVE_MINUTE").String()
				stock.High = item.Get("HIGH").String()
				stock.Hs = item.Get("HS").String()
				stock.Lb = item.Get("LB").String()
				stock.Low = item.Get("LOW").String()
				stock.Mcap = item.Get("MCAP").String()
				stock.Eps = item.Get("MFSUM").String()
				stock.Name = item.Get("NAME").String()
				stock.Open = item.Get("OPEN").String()
				stock.Percent = item.Get("PERCENT").String()
				stock.Price = item.Get("PRICE").String()
				stock.Pe = item.Get("PE").String()
				stock.Code = item.Get("SYMBOL").String()
				stock.Tcap = item.Get("TCAP").String()
				stock.Turnover = item.Get("TURNOVER").String()
				stock.UpDown = item.Get("UPDOWN").String()
				stock.Volume = item.Get("VOLUME").String()
				stock.Wb = item.Get("WB").String()
				stock.YestClose = item.Get("YESTCLOSE").String()
				stock.Zf = item.Get("ZF").String()
				stock.NetProfit = item.Get("MFRATIO.MFRATIO2").String()
				stock.TotalRevenue = item.Get("MFRATIO.MFRATIO10").String()
        stock.Date = date

        if err := db.Create(&stock).Error ; err != nil{
          logger.Error(err)
          return err
        }
			}
		}

		page++
		if page >= pageCount{
			break
		}

	}

	return nil
}
