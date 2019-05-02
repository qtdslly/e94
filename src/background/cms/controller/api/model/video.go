package apimodel

import (
	"background/newmovie/model"
	"background/common/logger"
	"github.com/jinzhu/gorm"
	"background/common/constant"
	"background/newmovie/service"
	"background/common/aes1"
)

type Video struct {
	Id       uint32  `json:"id"`
	Title    string  `json:"title"`
	Description    string  `json:"description"`
	Score    float64  `json:"score"`
	ThumbX      string  `json:"thumb_x"`
	ThumbY string  `json:"thumb_y"`
	PublishDate    string  `json:"publish_date"`
	Year uint32    `json:"year"`
	Language string    `json:"language"`
	Country string    `json:"country"`
	Directors string    `json:"directors"`
	Actors string    `json:"actors"`
	Tags string    `json:"tags"`
	Urls     []*PlayUrl `json:"urls"`
}

func VideoFromDb(src model.Video,db *gorm.DB) *Video {
	dst := Video{}
	dst.Id = src.Id
	dst.Title = src.Title
	dst.Description = src.Description
	dst.Score = src.Score
	dst.ThumbX = src.ThumbX
	dst.ThumbY = src.ThumbY
	dst.PublishDate = src.PublishDate
	dst.Year = src.Year
	dst.Language = src.Language
	dst.Country = src.Country
	dst.Directors = src.Directors
	dst.Actors = src.Actors
	dst.Tags = src.Tags

	var episode model.Episode
	if err := db.Where("video_id = ?",src.Id).First(&episode).Error ; err != nil{
		logger.Error(err)
		return nil
	}

	var playUrls []model.PlayUrl
	if err := db.Where("content_type = ? and content_id = ?",constant.MediaTypeEpisode,episode.Id).Find(&playUrls).Error ; err != nil{
		logger.Error(err)
		return nil
	}

	for _,playUrl := range playUrls{
		var pUrl PlayUrl
		pUrl.Id = playUrl.Id
		pUrl.Provider = playUrl.Provider
		pUrl.IsPlay = true
		if playUrl.OnLine{
			pUrl.Url = service.GetRealUrl(playUrl.Provider,playUrl.Url)
		}else{
			pUrl.Url = playUrl.Url
			pUrl.IsPlay = false
		}
		if pUrl.Url == ""{
			pUrl.Url = playUrl.Url
			pUrl.IsPlay = false
		}

		var err error
		pUrl.Url,err = aes1.Encrypt([]byte(pUrl.Url))
		if err != nil{
			logger.Error(err)
			return nil
		}

		dst.Urls = append(dst.Urls,&pUrl)
	}
	return &dst
}
