package service

import (
	"background/common/constant"
	"background/common/logger"
	apimodel "background/wechart/controller/api/model"
)

func SearchVideo(title string)(*apimodel.Video){

	var count int = 0

	var youkuVideo apimodel.Video

	err,title,description,actors,directors,thumb,pageUrl,publishDate := GetYoukuVideoInfoByTitle(title)
	if err == nil{
		youkuVideo.Title = title
		youkuVideo.Description = description
		youkuVideo.Actors = actors
		youkuVideo.Directors = directors
		youkuVideo.ThumbY = thumb
		youkuVideo.PageUrl = pageUrl
		youkuVideo.PublishDate = publishDate
		youkuVideo.ContentType = constant.MediaTypeEpisode
		youkuVideo.Provider = constant.ContentProviderYouKu
		count++
	}else{
		logger.Debug(err)
	}

	if count == 0{
		err,title,description,thumb,score,actors,directors,pageUrl := GetTencentVideoInfoByTitle(title)
		if err == nil{
			youkuVideo.Title = title
			youkuVideo.Description = description
			youkuVideo.Actors = actors
			youkuVideo.Directors = directors
			youkuVideo.ThumbY = thumb
			youkuVideo.PageUrl = pageUrl
			youkuVideo.Score = score
			youkuVideo.ContentType = constant.MediaTypeEpisode
			youkuVideo.Provider = constant.ContentProviderTencent
			count++
		}else{
			logger.Debug(err)
		}
	}

	if count == 0{
		err,title,score,area,description,actors,directors,thumb,pageUrl,publishDate := GetIqiyiVideoInfoByTitle(title)
		if err == nil{
			youkuVideo.Title = title
			youkuVideo.Description = description
			youkuVideo.Actors = actors
			youkuVideo.Directors = directors
			youkuVideo.ThumbY = thumb
			youkuVideo.PageUrl = pageUrl
			youkuVideo.Score = score
			youkuVideo.Country = area
			youkuVideo.ContentType = constant.MediaTypeEpisode
			youkuVideo.PublishDate = publishDate
			youkuVideo.Provider = constant.ContentProviderIqiyi
			count++
		}else{
			logger.Debug(err)
		}
	}


	return &youkuVideo
}