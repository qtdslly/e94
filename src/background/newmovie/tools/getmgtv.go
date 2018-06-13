package main

import (
	"background/newmovie/config"
	"background/common/util"
	"github.com/PuerkitoBio/goquery"
	"background/common/logger"
	"strings"
	"fmt"
	"background/newmovie/model"
	"time"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	"flag"
	"background/common/constant"
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

	model.InitModel(db)
	i := 1
	for {
		url := "https://list.mgtv.com/3/a4-a3-------a7-1-" + fmt.Sprint(i) +  "--b1-.html?channelId=3"
		GetMgMovie(url,db)
		i = i + 1
		if i == 42{
			break
		}
	}
}

func GetMgMovie(url string,db *gorm.DB){
	query := GetMgPageInfo(url)

	if query != nil{
		FilterMgMovieInfo(query,db)
	}
}

type MgtvVideoInfo struct {
	Area	             string             `gorm:"area" json:"area"`
	Duration                 string        `gorm:"duration" json:"duration"`
	Release                 string        `gorm:"release" json:"release"`
}
type MgtvDataInfo struct {
	Info                 MgtvVideoInfo          `gorm:"info" json:"info"`
}
type MgtvData struct{
	Code                int             `gorm:"code" json:"code"`
	Data                 MgtvDataInfo    `gorm:"data" json:"data"`
	Msg                string             `gorm:"msg" json:"msg"`
	Seqid                string             `gorm:"seqid" json:"seqid"`
}

func FilterMgMovieInfo(document *goquery.Document,db *gorm.DB)(){
	movieDoc := document.Find("body").Find(".m-main").Find(".m-content").Find(".m-content-wrapper").Find(".m-result").Find(".m-result-list").Find("ul").Find("li")

	movieDoc.Each(func(i int, s *goquery.Selection) {
		title := s.Find(".u-title").Eq(0).Text()
		url,_ := s.Find(".u-title").Eq(0).Attr("href")
		if !strings.Contains(url,"http"){
			url = "http:" + url
		}
		thumb_y,_ := s.Find("img").Eq(0).Attr("src")
		if !strings.Contains(thumb_y,"http"){
			thumb_y = "http:" + thumb_y
		}

		score := s.Find("a").Eq(0).Find("em").Eq(0).Text()

		logger.Debug(title)
		logger.Debug(url)
		logger.Debug(thumb_y)
		logger.Debug(score)

		directors,actors,country,tags,description := GetMgInfoByUrl(url)

		logger.Debug(country)
		logger.Debug(score)
		logger.Debug(directors)
		logger.Debug(description)
		logger.Debug(tags)

		logger.Debug(url)
		tmp := strings.Replace(url,"http://www.mgtv.com/b/","",-1)
		tmp = strings.Replace(tmp,".html","",-1)
		videoIds := strings.Split(tmp,"/")
		if len(videoIds) != 2{
			logger.Error("获取videoid错误!!!")
			return
		}

		p := time.Now()
		t := fmt.Sprint(p.Unix()) + "000"
		apiurl := "https://pcweb.api.mgtv.com/movie/list?video_id=" + videoIds[1] + "&cxid=&version=5.5.35&callback=jQuery18205330103375947974_1528805597948&_support=10000000&_=" + t

		requ, err := http.NewRequest("GET",apiurl,nil)

		resp, err := http.DefaultClient.Do(requ)
		if err != nil {
			logger.Debug("Proxy failed!")
			return
		}

		recv,err := ioutil.ReadAll(resp.Body)
		if err != nil{
			logger.Error(err)
			return
		}

		data := string(recv)
		if data == ""{
			return
		}

		data = strings.Replace(data,"jQuery18205330103375947974_1528805597948(","",-1)
		data = data[0:len(data) - 1]

		//data, _ = service.DecodeToGBK(data)
		logger.Debug(data)

		var mgData MgtvData
		if err = json.Unmarshal([]byte(data), &mgData); err != nil {
			logger.Error(err)
			return
		}

		publishDate := mgData.Data.Info.Release
		duration := mgData.Data.Info.Duration
		duration = strings.Replace(duration,"分钟","",-1)
		logger.Debug(publishDate)
		logger.Debug(duration)

		var video  model.Video
		video.Title = title
		video.Description = description
		video.ThumbY = thumb_y
		//movie.Year = year
		video.Country = country
		video.Actors = actors
		video.Tags = tags
		video.Directors = directors
		video.PublishDate = publishDate
		score1,_ := strconv.ParseFloat(score,10)
		video.Score = score1
		video.Pinyin = util.TitleToPinyin(video.Title)
		video.Status = constant.MediaStatusReleased

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
			if len(country) > 0{
				updateMap["country"] = country
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
			if len(publishDate) > 0{
				updateMap["publish_date"] = publishDate
			}
			if len(score) > 0{
				updateMap["score"] = score1
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
		dur,_ := strconv.Atoi(duration)
		episode.Duration = uint32(dur) * 60
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
				updateMap["score"] = score
			}
			if len(country) > 0{
				updateMap["duration"] = duration
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
		playUrl.Provider = constant.ContentProviderMgtv
		playUrl.Url = url
		playUrl.Disabled = false

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


func GetMgInfoByUrl(url string)(string,string,string,string,string){
	query := GetMgPageInfo(url)

	base := query.Find(".play-container").Find(".play-primary").Find(".play-primary").Find(".v-panel").Find(".v-panel-box").Find(".v-panel-info").Find(".extend").Find(".v-panel-extend").Find(".v-panel-meta")

	base = base.Find("p")
	directors := base.Eq(0).Text()
	directors = strings.Replace(directors,"导演：","",-1)
	directors = strings.Replace(directors," ","",-1)
	directors = strings.Replace(directors,"\n","",-1)

	actors := base.Eq(1).Text()
	actors = strings.Replace(actors,"主演：","",-1)
	actors = strings.Replace(actors," ","",-1)
	actors = strings.Replace(actors,"\n","",-1)

	country := base.Eq(2).Text()
	country = strings.Replace(country,"地区：","",-1)
	country = strings.Replace(country," ","",-1)
	country = strings.Replace(country,"\n","",-1)


	tags := base.Eq(3).Text()
	tags = strings.Replace(tags,"类型：","",-1)
	tags = strings.Replace(tags," ","",-1)
	tags = strings.Replace(tags,"\n","",-1)

	description := base.Eq(4).Text()
	description = strings.Replace(description,"简介：","",-1)
	description = strings.Replace(description," ","",-1)
	description = strings.Replace(description,"\n","",-1)

	return directors,actors,country,tags,description
}

func GetMgPageInfo(url string)( *goquery.Document){
	logger.Debug(url)
	query, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return nil
	}
	return query
}