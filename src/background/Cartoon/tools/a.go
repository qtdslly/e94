package main

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"background/common/logger"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"background/videodownload/config"
)

func main(){


	logger.SetLevel(config.GetLoggerLevel())

	url := "http://www.dytt8.net/html/tv/hytv/20180602/56938.html"
	getDytt8Info(url)
}

func getDytt8Info(url string)(bool){
	query, err := goquery.NewDocument(url)
	if err != nil {
		logger.Error(err)
		return false
	}

	baseInfo := query.Find(".co_content8").Text()

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
	var sorts []string
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
					sorts = append(sorts,fmt.Sprint(sort))
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
	for _, downloadUrl := range downloadUrls{
		logger.Debug(downloadUrl)
	}

	return true
}


