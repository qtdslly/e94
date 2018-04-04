package task

import (

	"background/lafter/model"
	"background/common/logger"

	"time"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func GetPageUrlByProvider(provider model.Provider,db *gorm.DB){
	var err error
	var pageUrls []model.PageUrl

	if err = db.Where("provider_id = ? and page_status = 0",provider.Id).Find(&pageUrls).Error ; err != nil{
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
			return
		}

		bd := query.Find("a")
		bd.Each(func(i int, s *goquery.Selection) {
			url,_ := s.Attr("href")
			if !strings.Contains(url,"http") && !strings.Contains(url,"https"){
				if url[0:1] != "/"{
					url = "/" + url
				}
				url = provider.Url + url
			}

			var pageUrl model.PageUrl
			if err = db.Where("url = ?",url).First(&pageUrl).Error ; err != nil{
				if err != gorm.ErrRecordNotFound{
					logger.Error(err)
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
			}
		})
	}

	GetPageUrlByProvider(provider,db)
}