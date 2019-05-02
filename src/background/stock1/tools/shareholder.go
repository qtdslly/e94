package main

import (
	"background/stock1/model"
	"background/stock1/config"

	"background/common/logger"
	"github.com/PuerkitoBio/goquery"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"flag"
	"log"
	"background/common/systemcall"
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

	// set file descriptor limit
	systemcall.SetFileLimit()

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	model.InitModel(db)

	var stocks []model.Stock
	if err := db.Order("code desc").Find(&stocks).Error ; err != nil{
		logger.Error(err)
		return
	}

	for _, stock := range stocks{
		if err := GetShareHolders(stock.Code,db) ; err != nil{
			logger.Error(err)
			continue
		}
	}
}
func GetShareHolders(code string,db *gorm.DB) error {
	apiurl := "http://quotes.money.163.com/f10/gdfx_" + code + ".html"
	document, err := goquery.NewDocument(apiurl)
	if err != nil {
		logger.Error(err)
		return err
	}

	trs := document.Find("#ltdateTable").Eq(0).Find("table").Eq(0).Find("tr")

	var date string = document.Find("#ltdate").Eq(0).Find("option").Eq(0).Text()
	trs.Each(func(i int, tr *goquery.Selection) {
		if i > 1 && i < 11{
			var name = tr.Find("td").Eq(0).Text()
			var count = strings.Replace(tr.Find("td").Eq(2).Text(),",","",-1)
			var percent = strings.Replace(tr.Find("td").Eq(1).Text(),"%","",-1)
			var remark = tr.Find("td").Eq(3).Text()

			//logger.Debug(name,",",count,",",percent,",",remark)

			var shareHolder model.ShareHolder
			shareHolder.Date = date
			shareHolder.Code = code
			shareHolder.Category = 0
			shareHolder.Name = name
			shareHolder.HoldCount = count
			shareHolder.Percent = percent
			shareHolder.Remark = remark

			updated := false
			if err := db.Where("date = ? and code = ? and category = 0 and name = ?",date,code,name).First(&shareHolder).Error ; err == nil{
				updated = true
			}

			if updated{
				//if err := db.Save(&shareHolder).Error ; err != nil{
				//	logger.Error(err)
				//	return
				//}
			}else{
				if err := db.Create(&shareHolder).Error ; err != nil{
					logger.Error(err)
					return
				}
			}

		}

	})

	trs1 := document.Find("#dateTable").Eq(0).Find("table").Eq(0).Find("tr")

	var date1 string = document.Find("#date").Eq(0).Find("option").Eq(0).Text()
	trs1.Each(func(i int, tr *goquery.Selection) {
		if i > 1 && i < 11{
			var name = tr.Find("td").Eq(0).Text()
			var count = strings.Replace(tr.Find("td").Eq(2).Text(),",","",-1)
			var percent = strings.Replace(tr.Find("td").Eq(1).Text(),"%","",-1)
			var remark = tr.Find("td").Eq(3).Text()

			//logger.Debug(name,",",count,",",percent,",",remark)

			var shareHolder model.ShareHolder
			shareHolder.Date = date1
			shareHolder.Code = code
			shareHolder.Category = 1
			shareHolder.Name = name
			shareHolder.HoldCount = count
			shareHolder.Percent = percent
			shareHolder.Remark = remark

			updated := false
			if err := db.Where("date = ? and code = ? and category = 1 and name = ?",date1,code,name).First(&shareHolder).Error ; err == nil{
				updated = true
			}

			if updated{
				//if err := db.Save(&shareHolder).Error ; err != nil{
				//	logger.Error(err)
				//	return
				//}
			}else{
				if err := db.Create(&shareHolder).Error ; err != nil{
					logger.Error(err)
					return
				}
			}

		}

	})

	return nil
}






























