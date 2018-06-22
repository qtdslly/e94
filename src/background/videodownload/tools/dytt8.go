package main

import (
	"background/common/logger"
	"background/videodownload/config"
	"background/videodownload/model"
	"background/common/constant"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"time"
	"golang.org/x/text/encoding/simplifiedchinese"
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

	found1 := false
	found2 := false
	go func(){
		for{
			if !GetDytt8MovieUrls(db){
				found1 = true
				return
			}
		}

	}()

	go func(){
		for{
			if !GetDytt8MovieInfos(db){
				found2 = true
				return
			}
		}

	}()

	for{
		if found1 && found2{
			break
		}
		time.Sleep(time.Minute)
	}
}

func GetDytt8MovieInfos(db *gorm.DB)bool{
	var pages []model.Page
	var err error
	if err = db.Where("url_status = ?",constant.DouBanCrawlStatusReady).First(&pages).Error ; err != nil{
		logger.Error(err)
		if err == gorm.ErrRecordNotFound{
			time.Sleep(time.Second * 60)
			return true
		}
		return false
	}

	for _ , page := range pages {
		time.Sleep(time.Second * 30)
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

func SaveMovieInfo(query *goquery.Document,db *gorm.DB)(bool){

	baseInfo1 := query.Find(".co_content8")
	if baseInfo1 == nil{
		return false
	}

	baseInfo := baseInfo1.Text()

	info,err := DecodeToGBK(baseInfo)
	if err != nil{
		logger.Error(err)
		return false
	}
	//fmt.Println(info)

	var entitle,title,year,country,types,language,subTitle,releaseDate,imdb,score,fileFormat,ratio,size,duration,directors,actors,description string
	info1 := strings.Split(info,"◎")
	for _ , value := range info1{
		value = strings.Replace(value," ","",-1)
		value = strings.Replace(value,"　","",-1)
		if strings.HasPrefix(value,"译名"){
			entitle = strings.Replace(value,"译名","",-1)
		}else if strings.HasPrefix(value,"片名"){
			title = strings.Replace(value,"片名","",-1)
		}else if strings.HasPrefix(value,"年代"){
			year = strings.Replace(value,"年代","",-1)
		}else if strings.HasPrefix(value,"产地"){
			country = strings.Replace(value,"产地","",-1)
		}else if strings.HasPrefix(value,"类别"){
			types = strings.Replace(value,"类别","",-1)
		}else if strings.HasPrefix(value,"语言"){
			language = strings.Replace(value,"语言","",-1)
		}else if strings.HasPrefix(value,"字幕"){
			subTitle = strings.Replace(value,"字幕","",-1) + "字幕"
		}else if strings.HasPrefix(value,"上映日期"){
			releaseDate = strings.Replace(value,"上映日期","",-1)
		}else if strings.HasPrefix(value,"IMDb评分"){
			imdb = strings.Replace(value,"IMDb评分","",-1)
		}else if strings.HasPrefix(value,"豆瓣评分"){
			score = strings.Replace(value,"豆瓣评分","",-1)
		}else if strings.HasPrefix(value,"文件格式"){
			fileFormat = strings.Replace(value,"文件格式","",-1)
		}else if strings.HasPrefix(value,"视频尺寸"){
			ratio = strings.Replace(value,"视频尺寸","",-1)
		}else if strings.HasPrefix(value,"文件大小"){
			size = strings.Replace(value,"文件大小","",-1)
		}else if strings.HasPrefix(value,"片长"){
			duration = strings.Replace(value,"片长","",-1)
		}else if strings.HasPrefix(value,"导演"){
			directors = strings.Replace(value,"导演","",-1)
		}else if strings.Contains(value,"主演"){
			actors = strings.Replace(value,"主演","",-1)
		}else if strings.HasPrefix(value,"简介"){
			description = strings.Replace(value,"简介","",-1)
			description = strings.Replace(description,"【下载地址】磁力链下载点击这里","",-1)
		}
		//logger.Debug(value)
	}

	table:= query.Find("table")

	var downloadUrls []string
	var sorts []uint32
	var urlTitles []string

	var sort int = 0
	table.Each(func(i int, ss *goquery.Selection) {
		if i != 0{
			s := ss.Find("a").Eq(0)
			downloadUrl,found := s.Attr("href")
			if found{
				urlTitle := s.Text()
				downloadUrl,err := DecodeToGBK(urlTitle)
				if err == nil{
					sort++
					downloadUrl = strings.Replace(downloadUrl," ","",-1)
					downloadUrl = strings.Replace(downloadUrl,"\n","",-1)
					downloadUrls = append(downloadUrls,downloadUrl)
					sorts = append(sorts,uint32(sort))
					urlTitles = append(urlTitles,urlTitle)
					logger.Debug(downloadUrl)

				}
			}else {
				logger.Debug(downloadUrl)
			}
		}
	})
	if err != nil{
		logger.Error(err)
		return false
	}

	entitle = strings.Replace(entitle,"\n","",-1)
	title = strings.Replace(title,"\n","",-1)
	year = strings.Replace(year,"\n","",-1)
	country = strings.Replace(country,"\n","",-1)
	types = strings.Replace(types,"\n","",-1)
	language = strings.Replace(language,"\n","",-1)
	subTitle = strings.Replace(subTitle,"\n","",-1)
	releaseDate = strings.Replace(releaseDate,"\n","",-1)
	imdb = strings.Replace(imdb,"\n","",-1)
	score = strings.Replace(score,"\n","",-1)
	fileFormat = strings.Replace(fileFormat,"\n","",-1)
	ratio = strings.Replace(ratio,"\n","",-1)
	size = strings.Replace(size,"\n","",-1)
	duration = strings.Replace(duration,"\n","",-1)
	directors = strings.Replace(directors,"\n","",-1)
	actors = strings.Replace(actors,"\n","/",-1)
	description = strings.Replace(description,"\n","",-1)

	logger.Debug(entitle)
	logger.Debug(title)
	logger.Debug(year)
	logger.Debug(country)
	logger.Debug(types)
	logger.Debug(language)
	logger.Debug(subTitle)
	logger.Debug(releaseDate)
	logger.Debug(imdb)
	logger.Debug(score)
	logger.Debug(fileFormat)
	logger.Debug(ratio)
	logger.Debug(size)
	logger.Debug(duration)
	logger.Debug(directors)
	logger.Debug(actors)
	logger.Debug(description)

	var video model.Video
	video.Title = title
	video.Year = year
	if err := db.Where("title = ?",video.Title).First(&video).Error ; err == nil{
		return true
	}

	video.EnglishTitle = entitle
	video.SubTitle = subTitle
	video.Description = description
	video.Country = country
	video.FileFormat = fileFormat
	video.Ratio = ratio
	video.Actors = actors
	video.Directors = directors
	video.Duration = duration
	video.Imdb = imdb
	video.Size = size
	video.Language = language
	video.ReleaseDate = releaseDate
	video.Score = score
	video.Types = types

	tx := db.Begin()
	if err := tx.Create(&video).Error ; err != nil{
		logger.Error(err)
		tx.Rollback()
		return false
	}
	for k, down := range downloadUrls{
		if k == len(downloadUrls) - 1{
			break
		}
		var download model.DownloadUrl
		download.Title = urlTitles[k]
		download.Sort = sorts[k]
		download.Url = down
		download.VideoId = video.Id
		download.Provider = constant.DownloadUrlProviderDytt8
		if err = tx.Create(&download).Error ; err != nil{
			logger.Error(err)
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}

func TrimNoNumber(str string)string{//只保留0-9和斜杠(/)
	rTitle := ([]rune)(str)
	result := ""
	for _, m := range rTitle {
		if m >= 47 && m <= 57{
			result += string(m)
		}
	}
	return result
}

func TrimStr(str string)string{
	str = strings.TrimLeft(str,"\n")
	str = strings.TrimLeft(str,"\r")
	str = strings.TrimRight(str,"\n")
	str = strings.TrimRight(str,"\r")
	str = strings.Replace(str,"\n\n","\n",-1)
	return str
}
func GetDytt8MovieUrls(db *gorm.DB)(bool){
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
				return false
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
		time.Sleep(time.Second * 25)
	}
	return true
}

func SaveUrls(document *goquery.Document,db *gorm.DB)(bool){
	query := document.Find("a")
	query.Each(func(i int, s *goquery.Selection) {
		url, found := s.Attr("href")
		if found{
			if (strings.Contains(url,"http") && strings.Contains(url,"www.dytt8.net")) || (!strings.Contains(url,".css") && !strings.Contains(url,".js") && !strings.Contains(url,".jpg") && !strings.Contains(url,".jpeg")&& !strings.Contains(url,".gif")&& !strings.Contains(url,".webp") && !strings.Contains(url,"ftp")){
				if !strings.Contains(url,"http") && !strings.Contains(url,"dytt8.net"){
					url = "http://www.dytt8.net" + url
				}else if !strings.Contains(url,"http:") && strings.Contains(url,"dytt8.net"){
					url = "http:" + url
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

func DecodeToGBK(text string) (string, error) {

	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}