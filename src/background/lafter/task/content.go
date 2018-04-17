package task

import (

	"background/lafter/model"
	"background/common/logger"

	"time"
	"fmt"
	"sync"

	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	//"github.com/axgle/mahonia"
)

func GetContentByProvider(provider model.Provider,db *gorm.DB){
	GetContent(provider,db)
}

func GetContent(provider model.Provider,db *gorm.DB){
	var err error
	var pageUrls []model.PageUrl

	if err = db.Where("provider_id = ? and url_status = 0 and url like ?",provider.Id,provider.Url + "%").Find(&pageUrls).Error ; err != nil{
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

	var process *sync.Mutex
	process = new(sync.Mutex)
	var Count int = 0

	for _ , p := range pageUrls{
		for{
			if Count > 5{
				time.Sleep(time.Millisecond * 100)
			}else{
				break
			}
		}

		go func(){
			process.Lock()
			Count++
			process.Unlock()

			query := GetPageInfo(p,db)

			if query == nil{
				process.Lock()
				Count--
				process.Unlock()
				return
			}

			if p.ProviderId == 4{
				GetWaiJiContent(&p,query,db)
			}else if p.ProviderId == 1 {
				GetJokejiContent(&p,query,db)
			}else{
				process.Lock()
				Count--
				process.Unlock()
				return
			}

			p.PageStatus = 1
			if err = db.Save(&p).Error ; err != nil{
				logger.Error(err)
				process.Lock()
				Count--
				process.Unlock()
				return
			}

			process.Lock()
			Count--
			process.Unlock()
		}()
		time.Sleep(time.Millisecond * 100)
	}
}

func GetPageInfo(p model.PageUrl,db *gorm.DB)(*goquery.Document){
	logger.Debug("url1:",p.Url)
	query, err := goquery.NewDocument(p.Url)
	logger.Debug("url2:",p.Url)
	if err != nil {
		logger.Error(err)
		p.PageStatus = 2
		p.UpdatedAt = time.Now()
		if err = db.Save(&p).Error ; err != nil{
			logger.Error(err)
			return nil
		}
		return nil
	}
	return query
}
func GetWaiJiContent(pageUrl *model.PageUrl,document *goquery.Document,db *gorm.DB){

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

func GetJokejiContent(pageUrl *model.PageUrl,document *goquery.Document,db *gorm.DB){

	contentId := "#text110"
	doc := document.Find(contentId)
	text := doc.Text()
	title := document.Find(".main").Find(".left").Find(".left_up").Find("h1").Eq(0).Text()
	if len(text) != 0{

		//logger.Debug("title:",title)
		//logger.Debug("text:",text)

		title, _ = DecodeToGBK(title)
		text, _ = DecodeToGBK(text)
		//enc := mahonia.NewEncoder("utf-8")
		//title = enc.ConvertString(title)
		//text = enc.ConvertString(text)

		logger.Debug("title:",title)
		logger.Debug("text:",text)

		titles := strings.Split(title,"->")
		title = titles[len(titles) - 1]

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