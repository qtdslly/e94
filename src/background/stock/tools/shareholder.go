package main

import (
	"background/stock/model"
	"background/stock/config"

	"background/common/logger"
	"github.com/PuerkitoBio/goquery"

	"golang.org/x/text/encoding/simplifiedchinese"

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

	model.InitModel(db)

	GetShareHolders("000963",db)
}
func GetShareHolders(code string,db *gorm.DB) error {
	apiurl := "http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CirculateStockHolder/stockid/"+code+"/displaytype/30.phtml"
	logger.Debug(apiurl)
	document, err := goquery.NewDocument(apiurl)
	if err != nil {
		logger.Error(err)

		return err
	}

	trs := document.Find("#CirculateShareholderTable").Eq(0).Find("tr")


	var date string = ""
	trs.Each(func(i int, tr *goquery.Selection) {
		if i == 1{
			date = strings.Replace(tr.Text(),"截至日期","",-1)
		}

		if i > 3 && i < 14{
			var name,_ = DecodeToGBK(tr.Find("td").Eq(1).Text())
			var count,_ = DecodeToGBK(tr.Find("td").Eq(2).Text())
			var percent,_ = DecodeToGBK(tr.Find("td").Eq(3).Text())
			var property,_ = DecodeToGBK(tr.Find("td").Eq(4).Text())

			logger.Debug(name,",",count,",",percent,",",property)

			var shareHolder model.ShareHolder
			shareHolder.Date = date
			shareHolder.Code = code
			shareHolder.Name = name
			shareHolder.HoldCount = count
			shareHolder.Percent = percent
			shareHolder.Property = property

			if err := db.Create(&shareHolder).Error ; err != nil{
				logger.Error(err)
				return
			}
		}

	})

	return nil
}


func DecodeToGBK(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}

































