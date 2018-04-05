package task

import (

	"background/lafter/model"
	"background/common/logger"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func GetContentByProvider(provider model.Provider,db *gorm.DB){
	for{
		GetContent(provider,db)
		time.Sleep(time.Second * 10)
	}

}

func GetContent(provider model.Provider,db *gorm.DB){
	var err error
	var pageUrls []model.PageUrl

	if err = db.Where("provider_id = ? and url_status = 0",provider.Id).Find(&pageUrls).Error ; err != nil{
		logger.Error(err)
		return
	}

	if len(pageUrls) == 0{
		var pageUrl model.PageUrl
		if err = db.Where("url = ?",provider.Url).First(&pageUrl).Error ; err == nil{
			logger.Debug(err)
			return
		}
		pageUrl.PageStatus = 0
		pageUrl.ProviderId = provider.Id
		pageUrl.Url = provider.Url
		pageUrl.CreatedAt = time.Now()
		pageUrl.UpdatedAt = time.Now()
		if err = db.Save(&pageUrl).Error ; err != nil{
			logger.Error(err)
			return
		}
		return
	}

	for _ , p := range pageUrls{
		query, err := goquery.NewDocument(p.Url)
		if err != nil {
			logger.Error(err)

			p.PageStatus = 2
			p.UpdatedAt = time.Now()
			if err = db.Save(&p).Error ; err != nil{
				logger.Error(err)
				return
			}
			return

			return
		}

		if p.ProviderId == 1{
			GetXiaoHuaJiContent(p,query)
		}

		p.PageStatus = 1
		if err = db.Save(&p).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}

func GetXiaoHuaJiContent(pageUrl model.PageUrl,document *goquery.Document){
	//des := document.Find("#text110").Text()
	//title := document.Find(".main").Find("h1").
	//if len(des) != 0{
	//	var content model.Content
	//	content.Title = ""
	//	content.Content = des
	//}
}