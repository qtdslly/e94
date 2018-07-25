package main

import (
	"background/newmovie/config"
	"background/common/util"
	"background/common/logger"
	uutil "background/newmovie/util"
	"background/newmovie/model"
	"strings"
	"fmt"
	"github.com/PuerkitoBio/goquery"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"flag"
	"background/common/constant"
	"strconv"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

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

	model.InitModel(db)
	i := 1
	for {
		url := "http://www.hanju.cc/hanju/list_8_" + fmt.Sprint(i) + ".html"
		GetHanJuInfo(url,db)
		i = i + 1
		if i == 68{
			break
		}
	}
}


func GetHanJuInfo(url string,db *gorm.DB){
	document, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return
	}
	movieDoc := document.Find(".sdlist").Eq(1).Find(".sdlist-l").Eq(0).Find(".sdlist-lcon").Eq(0).Find(".stab-con").Eq(0).Find("dl").Eq(0).Find("dd").Eq(0).Find("li")

	var directors ,writer,totalEpisode,description string
	movieDoc.Each(func(i int, s *goquery.Selection) {


		apiUrl1,_ := s.Find("a").Eq(0).Attr("href")
		if !strings.Contains(apiUrl1,"http"){
			apiUrl1 = "http://www.hanju.cc" + apiUrl1
		}
		logger.Debug(apiUrl1)
		doc, err := goquery.NewDocument(apiUrl1)
		if err != nil {
			logger.Debug(apiUrl1)
			logger.Error(err)
			return
		}

		videoType := document.Find("#sdlist").Find(".sdlist").Eq(0).Find(".pleft").Eq(0).Find("a").Eq(2).Text()
		if strings.Contains(videoType,"电影") || strings.Contains(videoType,"综艺"){
			return
		}

		mDoc := doc.Find(".vothercon").Eq(0).Text()
		mDoc,_ = uutil.DecodeToGBK(mDoc)
		ss := strings.Split(mDoc,"[")
		for _,text := range ss{
			if strings.Contains(text,"导 演"){
				directors = text
				directors = strings.Replace(directors,"导 演]: ","",-1)
				if strings.Contains(directors,"（"){
					directors = directors[:strings.Index(directors,"（")]
				}
			}

			if strings.Contains(text,"编 剧"){
				writer = text
				writer = strings.Replace(writer,"编 剧]: ","",-1)
				if strings.Contains(writer,"（"){
					writer = writer[:strings.Index(writer,"（")]
				}
			}

			if strings.Contains(text,"集 数"){
				totalEpisode = text
				totalEpisode = strings.Replace(totalEpisode,"集 数]: ","",-1)
				totalEpisode = uutil.OnlyNumber(totalEpisode)
			}

			if strings.Contains(text,"简 介"){
				description = text
				description = strings.Replace(description,"简 介]: ","",-1)
			}
		}


		title := s.Find(".jd_info").Eq(0).Find(".tit").Eq(0).Find("a").Eq(0).Text()

		currentEpisode := s.Find(".lionhover").Eq(0).Find(".left_info").Text()
		currentEpisode = uutil.OnlyNumber(currentEpisode)
		thumb_y,_ := s.Find("img").Eq(0).Attr("src")
		if !strings.Contains(thumb_y,"http"){
			thumb_y = "http:" + thumb_y
		}

		actors := s.Find(".lionhover").Eq(0).Find(".right_info").Eq(0).Find(".actor").Text()
		//description := s.Find(".lionhover").Eq(0).Find(".right_info").Eq(0).Find(".descr").Text()

		actors = strings.Replace(actors,"主演","",-1)
		actors = strings.TrimLeft(actors,":")
		actors = strings.TrimLeft(actors,"：")

		//description = strings.Replace(description,"剧情：","",-1)

		title ,_ = uutil.DecodeToGBK(title)
		actors ,_ = uutil.DecodeToGBK(actors)

		actors = strings.Replace(actors,"\n","",-1)
		directors = strings.Replace(directors,"\n","",-1)
		writer = strings.Replace(writer,"\n","",-1)

		actors = strings.Replace(actors,"主演：","",-1)

		actors = strings.Replace(actors," ","/",-1)
		writer = strings.Replace(writer," ","/",-1)
		directors = strings.Replace(directors," ","/",-1)

		logger.Debug(title)
		logger.Debug(currentEpisode)
		logger.Debug(totalEpisode)
		logger.Debug(thumb_y)
		logger.Debug(actors)
		logger.Debug(directors)
		logger.Debug(writer)
		logger.Debug(description)
		var video  model.Video
		video.Actors = actors
		video.Title = title
		video.Description = description
		video.Actors = actors
		video.Writer = writer
		video.ThumbY = thumb_y
		video.Pinyin = util.TitleToPinyin(video.Title)
		video.OnLine = constant.MediaStatusOnLine
		num,_ := strconv.Atoi(currentEpisode)
		current := uint32(num)
		video.CurrentEpisode = current
		num,_ = strconv.Atoi(totalEpisode)
		total := uint32(num)
		video.TotalEpisode = total

		if err := db.Where("title = ? and category = '韩剧'",video.Title).First(&video).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&video)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["description"] = description
			}
			if len(thumb_y) > 0{
				updateMap["thumb_y"] = thumb_y
			}
			if total > 0{
				updateMap["total_episode"] = total
			}
			if current > 0{
				updateMap["current_episode"] = current
			}

			if err = db.Model(model.Video{}).Where("id=?", video.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}

		episodeDoc := doc.Find(".playbox").Eq(0).Find(".list").Eq(0).Find(".abc").Eq(0).Find("a")
		episodeDoc.Each(func(i int, s1 *goquery.Selection) {
			episodeTitle := s1.Text()
			episodeTitle,_ = uutil.DecodeToGBK(episodeTitle)
			pUrl , _ := s1.Attr("href")
			if !strings.Contains(pUrl,"http"){
				pUrl = "http://www.hanju.cc" + pUrl
			}
			logger.Debug(episodeTitle,"\t",pUrl)

			if strings.Contains(episodeTitle,"预告"){
				return
			}
			var episode model.Episode
			episode.Title = video.Title + episodeTitle
			episode.VideoId = video.Id
			episode.Description = description
			num := uutil.OnlyNumber(episodeTitle)
			sort,_ := strconv.Atoi(num)
			episode.Sort = uint32(sort)
			episode.Number = num
			episode.Pinyin = util.TitleToPinyin(episode.Title)

			if err := db.Where("video_id = ? and sort = ?",video.Id,episode.Sort).First(&episode).Error ; err == gorm.ErrRecordNotFound{
				db.Create(&episode)
			}

			var playUrl model.PlayUrl
			playUrl.Title = episode.Title
			playUrl.ContentType = constant.MediaTypeEpisode
			playUrl.ContentId = episode.Id
			playUrl.Provider = constant.ContentProviderHanju
			playUrl.Url = pUrl
			playUrl.OnLine = constant.MediaStatusOnLine

			if err := db.Where("content_id = ? and content_type = ? and provider = ?",episode.Id,playUrl.ContentType,playUrl.Provider).First(&playUrl).Error ; err == gorm.ErrRecordNotFound{
				db.Create(&playUrl)
			}
		})
		logger.Debug("===========================================================")
	})
}
