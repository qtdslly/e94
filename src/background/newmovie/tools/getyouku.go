package main

import (
	"background/newmovie/config"
	"background/common/constant"
	"background/common/util"
	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"strings"
	"fmt"
	"background/newmovie/model"
	"github.com/jinzhu/gorm"
	"flag"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)

	i := 14
	for {
		url := "http://list.youku.com/category/show/c_96_pt_1_s_1_d_1_p_" + fmt.Sprint(i) +  ".html?spm=a2h1n.8251845.0.0"
		GetYouKuMovie(url,db)
		i = i + 1
		if i == 31{
			break
		}
	}

	//GetYouKuPageInfo("http://v.youku.com/v_show/id_XNzA4ODY0NzQ0.html")
}


func GetYouKuMovie(url string,db *gorm.DB){
	query := GetYouKuPageInfo(url)

	if query != nil{
		FilterYouKuMovieInfo(query,db)
	}
}

func FilterYouKuMovieInfo(document *goquery.Document,db *gorm.DB)(){
	movieDoc := document.Find("body").Find(".s-body").Find(".yk-content").Find(".vaule_main").Find(".box-series").Find(".panel").Find(".yk-col4")

	movieDoc.Each(func(i int, s *goquery.Selection) {
		a := s.Find(".p-thumb").Eq(0)
		title,_ := a.Find("a").Eq(0).Attr("title")
		url,_ := a.Find("a").Eq(0).Attr("href")
		if !strings.Contains(url,"http"){
			url = "http:" + url
		}
		thumb_y,_ := a.Find("img").Eq(0).Attr("src")

		dd := s.Find(".info-list").Eq(0).Find(".actor").Eq(0).Find("a")

		actors := ""
		dd.Each(func(i int, t *goquery.Selection) {
			actors += t.Text() + ","
		})

		if len(actors) > 0{
			actors = actors[0:len(actors) - 1]
		}

		logger.Debug(title)
		logger.Debug(url)
		logger.Debug(thumb_y)
		logger.Debug(actors)
		year,publish_date,score,directors,country,description,tags := GetYoukMovieOtherInfo(url)
		logger.Debug(year)
		logger.Debug(publish_date)
		logger.Debug(score)
		logger.Debug(directors)
		logger.Debug(country)
		logger.Debug(description)
		logger.Debug(tags)

		var video  model.Video
		video.Actors = actors
		video.Title = title
		video.Description = description
		video.Directors = directors
		video.PublishDate = publish_date
		score1,_ := strconv.ParseFloat(score,10)

		video.Score = score1
		video.ThumbY = thumb_y
		video.Year = year
		video.Country = country
		video.Tags = tags
		video.TotalEpisode = 1
		video.OnLine = constant.OnlineTypeTrue
		now := time.Now()
		video.CreatedAt = now
		video.UpdatedAt = now

		if err := db.Where("title = ?",video.Title).First(&video).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&video)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["description"] = description
			}
			if len(thumb_y) > 0{
				updateMap["thumb_y"] = thumb_y
			}

			if len(actors) > 0{
				updateMap["actors"] = actors
			}
			if len(tags) > 0{
				updateMap["tags"] = tags
			}
			if len(directors) > 0{
				updateMap["directors"] = directors
			}
			if len(publish_date) > 0{
				updateMap["publish_date"] = publish_date
			}
			if len(score) > 0{
				updateMap["score"] = score
			}

			if err = db.Model(model.Video{}).Where("id=?", video.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}

		var episode model.Episode
		episode.Title = title
		episode.VideoId = video.Id
		episode.Description = description
		episode.Score = score1
		episode.Pinyin = util.TitleToPinyin(video.Title)

		episode.CreatedAt = now
		episode.UpdatedAt = now
		if err := db.Where("video_id = ?",video.Id).First(&episode).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&episode)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["description"] = description
			}
			if len(score) > 0{
				updateMap["score"] = score1
			}

			if err = db.Model(model.Episode{}).Where("id = ?", episode.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}

		var playUrl model.PlayUrl
		playUrl.Title = episode.Title
		playUrl.ContentType = constant.MediaTypeEpisode
		playUrl.ContentId = episode.Id
		playUrl.Provider = constant.ContentProviderYouKu
		playUrl.Url = url
		playUrl.OnLine = constant.OnlineTypeTrue

		playUrl.CreatedAt = now
		playUrl.UpdatedAt = now
		if err := db.Where("content_id = ? and content_type = ? and provider = ?",episode.Id,playUrl.ContentType,playUrl.Provider).First(&playUrl).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&playUrl)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["url"] = url
			}

			if err = db.Model(model.PlayUrl{}).Where("id = ?", playUrl.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}

	})
}

func GetYoukMovieOtherInfo(url string)(uint32,string,string,string,string,string,string){
	if url == ""{
		return 0,"","","","","",""
	}

	query := GetYouKuPageInfo(url)
	newUrl,_ := query.Find("#bpmodule-playpage-righttitle-code").Find(".tvinfo").Eq(0).Find("h2").Eq(0).Find("a").Eq(0).Attr("href")
	//publish_date := query.Find("#bpmodule-playpage-lefttitle").First(".player-title").First(".title-wrap").First(".desc").First(".video-status").First(".bold").First("span").Text()
	//description := query.Find("#module_basic_intro").First(".mod").First(".c").First(".tab-c").First(".summary-wrap").First(".summary").Text()

	if newUrl == ""{
		return 0,"","","","","",""
	}
	if !strings.Contains(newUrl,"http"){
		newUrl = "http:" + newUrl
	}
	logger.Debug(newUrl)
	q := GetYouKuPageInfo(newUrl)
	base := q.Find("body").Find(".s-body").Find(".yk-content").Find(".mod").Find(".p-base").Find("ul").Eq(0)

	year := base.Find(".p-row").Eq(0).Find("span").Eq(0).Text()
	publish_date := base.Children().Eq(2).Eq(0).Find("span").Eq(0).Text()
	score := base.Find(".p-score").Eq(0).Find("span").Eq(1).Text()
	directors := base.Children().Eq(6).Find("a").Eq(0).Text()
	country := base.Children().Eq(7).Find("a").Eq(0).Text()
	description := base.Children().Eq(12).Find("span").Eq(0).Text()

	tagDoc := base.Children().Eq(8).Find("a")
	var tags string
	tagDoc.Each(func(i int, s *goquery.Selection) {
		tags += s.Text() + ","

	})

	if len(tags) > 0{
		tags = tags[0:len(tags) - 1]
	}

	publish_date = strings.Replace(publish_date,"上映：","",-1)
	description = strings.Replace(description,"简介：","",-1)
	var year1 int
	if year != ""{
		year1,_ = strconv.Atoi(year)
	}else{
		year1 = 0
	}
	return uint32(year1),publish_date,score,directors,country,description,tags
}
func GetYouKuPageInfo(url string)( *goquery.Document){
	logger.Debug(url)
	query, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return nil
	}
	return query
}