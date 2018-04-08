package task

import (

	"background/lafter/model"
	"background/common/logger"

	"time"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func GetContentByProvider(provider model.Provider,db *gorm.DB){
	GetContent(provider,db)
	for{

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
		}

		if p.ProviderId == 4{
			GetWaiJiContent(p,query,db)
		}else{
			continue
		}

		p.PageStatus = 1
		if err = db.Save(&p).Error ; err != nil{
			logger.Error(err)
			return
		}
	}
}

func GetWaiJiContent(pageUrl model.PageUrl,document *goquery.Document,db *gorm.DB){

	var i = 0
	for {
		i++
		contentId := "#content-" + fmt.Sprint(i)
		doc := document.Find(contentId)
		text := doc.Find(".c").Text()
		title := doc.Find("#title").Text()
		if len(text) != 0{
			var content model.Content
			content.Title = title
			content.Content = text
			content.PageId = pageUrl.Id
			content.CreatedAt = time.Now()
			content.UpdatedAt = time.Now()
			if err := db.Save(&content).Error ; err != nil{
				logger.Error(err)
				return
			}
		}

		pageUrl.UrlStatus = 1
		if err := db.Save(&pageUrl).Error; err != nil{
			logger.Error(err)
			return
		}
		
		if i == 10{
			break
		}
	}
}