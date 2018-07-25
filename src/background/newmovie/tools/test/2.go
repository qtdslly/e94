package main

import (
	"background/newmovie/config"
	"background/common/logger"
	uutil "background/newmovie/util"
	"strings"
	"fmt"
	"github.com/PuerkitoBio/goquery"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	url := "http://www.hanju.cc/hanju/list_8_" + fmt.Sprint(3) + ".html"
	GetHanJuInfo(url)

}


func GetHanJuInfo(url string){
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
					logger.Debug(text)
					logger.Debug(strings.Index(writer,"（"))
					writer = directors[:strings.Index(writer,"（")]
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

		logger.Debug("===========================================================")
	})
}
