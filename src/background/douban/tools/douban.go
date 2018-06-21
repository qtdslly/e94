package main

import (
	"background/common/logger"
	"background/douban/config"
	"background/douban/model"
	"background/common/constant"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

func main(){
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

	for{
		if !GetDoubanMovieUrls(db){
			break
		}
		if !GetDoubanMovieInfos(db){
			break
		}
	}
}

func GetDoubanMovieInfos(db *gorm.DB)bool{
	var pages []model.Page
	var err error
	if err = db.Where("url_status = ?",constant.DouBanCrawlStatusReady).Find(&pages).Error ; err != nil{
		logger.Error(err)
		return false
	}

	for _ , page := range pages {
		query, err := goquery.NewDocument(page.Url)
		if err != nil {
			page.PageStatus = constant.DouBanCrawlStatusError
			if err = db.Save(&page).Error; err != nil {
				logger.Error(err)
				return false
			}
		}

		if !SaveMovieInfo(query, db ){
			page.UrlStatus = constant.DouBanCrawlStatusError
		}else{
			page.UrlStatus = constant.DouBanCrawlStatusSuccess
		}
		if err = db.Save(&page).Error ; err != nil{
			logger.Error(err)
			return false
		}
	}
	return true
}

func SaveMovieInfo(document *goquery.Document,db *gorm.DB)(bool){
	year := document.Find(".year").Eq(0).Text()
	year = strings.Replace(year,"(","",-1)
	year = strings.Replace(year,")","",-1)

	baseInfo := document.Find("#info").Find("span")
	var actors,directors,writer,types,officialUrl,country,language,releaseDate,duration,alias,imdb string

	baseInfo.Each(func(i int, s *goquery.Selection) {
		value := s.Text()
		if strings.Contains(value,"导演"){
			actors = strings.Replace(value,"导演:","",-1)
		}
		if strings.Contains(value,"编剧"){
			writer = strings.Replace(value,"编剧:","",-1)
		}
		if strings.Contains(value,"主演"){
			actors = strings.Replace(value,"主演:","",-1)
		}
		if strings.Contains(value,"类型"){
			types = strings.Replace(value,"类型:","",-1)
		}
		if strings.Contains(value,"官方网站"){
			officialUrl = strings.Replace(value,"官方网站:","",-1)
		}
		if strings.Contains(value,"制片国家/地区"){
			country = strings.Replace(value,"制片国家/地区:","",-1)
		}
		if strings.Contains(value,"语言"){
			language = strings.Replace(value,"语言:","",-1)
		}
		if strings.Contains(value,"上映日期"){
			releaseDate = strings.Replace(value,"上映日期:","",-1)
		}
		if strings.Contains(value,"片长"){
			duration = strings.Replace(value,"片长:","",-1)
			duration = strings.Replace(duration,"分钟","",-1)
		}
		if strings.Contains(value,"又名"){
			alias = strings.Replace(value,"又名:","",-1)
		}
		if strings.Contains(value,"IMDb链接"){
			imdb = strings.Replace(value,"IMDb链接:","",-1)
		}
	})

	score := document.Find("strong[property='v:average']").Eq(0).Text()
	comments := document.Find("strong[property='v:votes']").Eq(0).Text()
	description := document.Find("span[property='v:summary']").Eq(0).Text()
	title := document.Find("title").Eq(0).Text()
	title = strings.Replace(title,"(豆瓣)","",-1)

	subjectId,_ := document.Find("input[name='target-id']").Eq(0).Attr("value")

	thumb_y,_ := document.Find("#mainpic").Find("img").Eq(0).Attr("src")
	if strings.Contains(thumb_y,"https"){
		thumb_y = "https:" + thumb_y
	}

	actors = strings.Replace(actors," ","",-1)
	directors = strings.Replace(directors," ","",-1)
	writer = strings.Replace(writer," ","",-1)
	types = strings.Replace(types," ","",-1)
	officialUrl = strings.Replace(officialUrl," ","",-1)
	country = strings.Replace(country," ","",-1)
	language = strings.Replace(language," ","",-1)
	releaseDate = strings.Replace(releaseDate," ","",-1)
	duration = strings.Replace(duration," ","",-1)
	alias = strings.Replace(alias," ","",-1)
	imdb = strings.Replace(imdb," ","",-1)
	score = strings.Replace(score," ","",-1)
	comments = strings.Replace(comments," ","",-1)
	description = strings.Replace(description," ","",-1)
	title = strings.Replace(title," ","",-1)

	var movie model.Movie
	movie.Title = title
	num,_ := strconv.Atoi(year)
	movie.Year = uint32(num)
	if err := db.Where("title = ? and year = ?",movie.Title,movie.Year).First(&movie).Error ; err == nil{
		return true
	}
	movie.Description = description
	movie.Alias = alias
	num,_ = strconv.Atoi(year)
	movie.Comments = uint32(num)
	movie.Country = country
	movie.Directors = directors
	num,_ = strconv.Atoi(duration)
	movie.Duration = uint32(num)
	movie.Imdb = imdb
	movie.Language = language
	movie.OfficialUrl = officialUrl
	movie.ReleaseDate = releaseDate
	sc,_ := strconv.ParseFloat(score,10)
	movie.Score = sc
	num,_ = strconv.Atoi(subjectId)
	movie.SubjectId = uint32(num)
	if thumb_y != "https:"{
		movie.ThumbY = thumb_y
	}
	if err := db.Create(&movie).Error ; err != nil{
		logger.Error(err)
		return false
	}
	return true
}
func GetDoubanMovieUrls(db *gorm.DB)(bool){
	var pages []model.Page
	var err error
	if err = db.Where("page_status = ?",constant.DouBanCrawlStatusReady).Find(&pages).Error ; err != nil{
		logger.Error(err)
		return false
	}

	for _ , page := range pages{
		query , err := goquery.NewDocument(page.Url)
		if err != nil{
			page.PageStatus = constant.DouBanCrawlStatusError
			if err = db.Save(&page).Error ; err != nil{
				logger.Error(err)
				return
			}
			return false
		}

		if !SaveUrls(query, db ){
			page.PageStatus = constant.DouBanCrawlStatusError
		}else{
			page.PageStatus = constant.DouBanCrawlStatusSuccess
		}
		if err = db.Save(&page).Error ; err != nil{
			logger.Error(err)
			return false
		}
	}
	return true
}

func SaveUrls(document *goquery.Document,db *gorm.DB)(bool){
	query := document.Find("a")
	query.Each(func(i int, s *goquery.Selection) {
		url, found := s.Attr("href")
		if found{
			if (strings.Contains(url,"http") && strings.Contains(url,"douban.com")) && !strings.Contains(url,".css") && !strings.Contains(url,".js") && !strings.Contains(url,".jpg") && !strings.Contains(url,".jpeg")&& !strings.Contains(url,".gif")&& !strings.Contains(url,".webp") {
				if !strings.Contains(url,"http") && !strings.Contains(url,"movie.douban.com"){
					url = "https://movie.douban.com" + url
				}else if !strings.Contains(url,"https:") && strings.Contains(url,"movie.douban.com"){
					url = "https:" + url
				}

				var page model.Page
				page.Url = url
				if err := db.Where("url = ?",page.Url).First(&page).Error ; err == gorm.ErrRecordNotFound{
					page.PageStatus = constant.DouBanCrawlStatusReady
					page.UrlStatus = constant.DouBanCrawlStatusReady
					if err := db.Create(&page).Error ; err != nil{
						logger.Error(err)
					}
				}
			}
		}
	})
	return true
}