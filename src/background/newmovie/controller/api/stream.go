package api

import (
	"fmt"
	"net/http"
	"background/newmovie/service"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	apimodel "background/newmovie/controller/api/model"
	"background/common/util"
)

func StreamListHandler(c *gin.Context) {

	type param struct {
		ResourceGroupId uint32    `form:"id"`
		Limit           int       `form:"limit" binding:"required"`
		Offset          int       `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var streams []*model.Stream
	if err = db.Order("stream.sort asc").Offset(p.Offset).Limit(p.Limit).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ? and stream.on_line = 1", p.ResourceGroupId).Find(&streams).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var count uint32
	if err = db.Model(&model.Stream{}).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ?", p.ResourceGroupId).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiStream struct {
		Id    uint32 `json:"id"`
		Title string `json:"title"`
		Thumb string `json:"thumb"`
	}

	var apiStreams []*ApiStream
	for _, stream := range streams {
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Title = stream.Title
		apiStream.Thumb = stream.Thumb
		apiStreams = append(apiStreams, &apiStream)
	}

	var hasMore bool = true
	if len(apiStreams) < p.Limit {
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiStreams, "count":count, "has_more":hasMore})
}

func StreamDetailHandler(c *gin.Context) {

	type param struct {
		Id             uint32 `form:"id" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	installationId := c.MustGet(constant.ContextInstallationId).(uint64)

	var stream model.Stream
	if err = db.Where("id = ? and on_line = ?", p.Id, constant.MediaStatusOnLine).First(&stream).Error; err != nil {
		logger.Error("query video err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiStream struct {
		Id      uint32                `json:"id"`
		IsDigg  bool                  `json:"is_digg"`
		Title   string                `json:"title"`
		Thumb   string                `json:"thumb"`
		PlayUrl []*apimodel.PlayUrl   `json:"play_url"`
	}

	var apiStream ApiStream
	apiStream.Id = stream.Id
	apiStream.Thumb = stream.Thumb
	apiStream.Title = stream.Title
	var playUrls []model.PlayUrl
	if err := db.Order("ready asc").Where("content_type = 4 and content_id = ?", stream.Id).Find(&playUrls).Error; err != nil {
		logger.Error("query play_url err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	apiStream.IsDigg = false
	if installationId != 0{
		var contentAction model.ContentAction
		err  = db.Where("installation_id = ? and content_type = ? and content_id = ? and action = 1",installationId,constant.MediaTypeStream,p.Id).First(&contentAction).Error ;
		if err == nil{
			apiStream.IsDigg = true
		}
	}

	for _, playUrl := range playUrls {
		var pUrl *apimodel.PlayUrl
		pUrl = apimodel.PlayUrlFromDb(playUrl)
		if pUrl.IsPlay {
			apiStream.PlayUrl = append(apiStream.PlayUrl, pUrl)
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiStream})
}

func SearchHandler(c *gin.Context) {

	type param struct {
		Title  string    `form:"title" binding:"required"`
		Limit  int       `form:"limit" binding:"required"`
		Offset int       `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var streams []model.Stream
	if err = db.Offset(p.Offset).Limit(p.Limit).Order("sort asc").Where("title like ? and on_line = ?", "%" + p.Title + "%", constant.MediaStatusOnLine).Find(&streams).Error; err != nil {
		logger.Error("query video err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	//title,score,area,description,actors,directors,thumb,pageUrl,publishDate
	type ApiStream struct {
		Id          uint32                `json:"id"`
		Title       string                `json:"title"`
		Thumb       string                `json:"thumb"`
		ContentType uint8                 `json:"content_type"`
		Actors      string                `json:"actors"`
		Country     string                `json:"country"`
		Directors   string                `json:"directors"`
		PageUrl     string                `json:"pageUrl"`
		PublishDate string                `json:"publish_date"`
		Score       string                `json:"score"`
		Description string                `json:"description"`
		Provider    uint32                `json:"provider"`
	}

	var apiModels []*ApiStream
	for _, stream := range streams {
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Thumb = stream.Thumb
		apiStream.Title = stream.Title
		apiStream.ContentType = constant.MediaTypeStream
		apiModels = append(apiModels, &apiStream)
	}

	var count int
	if err = db.Model(&model.Stream{}).Where("title like ? and on_line = ?", "%" + p.Title + "%", constant.MediaStatusOnLine).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count == 0{
		if p.Offset == 0{
			var youkuVideo ApiStream

			err,title,description,actors,directors,thumb,pageUrl,publishDate := service.GetYoukuVideoInfoByTitle(p.Title)
			if err == nil{
				youkuVideo.Title = title
				youkuVideo.Description = description
				youkuVideo.Actors = actors
				youkuVideo.Directors = directors
				youkuVideo.Thumb = thumb
				youkuVideo.PageUrl = pageUrl
				youkuVideo.PublishDate = publishDate
				youkuVideo.ContentType = constant.MediaTypeEpisode
				youkuVideo.Provider = constant.ContentProviderYouKu
				count++
			}

			if count == 0{
				err,title,description,thumb,score,actors,directors,pageUrl := service.GetTencentVideoInfoByTitle(p.Title)
				if err == nil{
					youkuVideo.Title = title
					youkuVideo.Description = description
					youkuVideo.Actors = actors
					youkuVideo.Directors = directors
					youkuVideo.Thumb = thumb
					youkuVideo.PageUrl = pageUrl
					youkuVideo.Score = score
					youkuVideo.ContentType = constant.MediaTypeEpisode
					youkuVideo.Provider = constant.ContentProviderTencent
					count++
				}
			}

			if count == 0{
				err,title,score,area,description,actors,directors,thumb,pageUrl,publishDate := service.GetIqiyiVideoInfoByTitle(p.Title)
				if err == nil{
					youkuVideo.Title = title
					youkuVideo.Description = description
					youkuVideo.Actors = actors
					youkuVideo.Directors = directors
					youkuVideo.Thumb = thumb
					youkuVideo.PageUrl = pageUrl
					youkuVideo.Score = score
					youkuVideo.Country = area
					youkuVideo.ContentType = constant.MediaTypeEpisode
					youkuVideo.PublishDate = publishDate
					youkuVideo.Provider = constant.ContentProviderIqiyi
					count++
				}
			}
			if count > 0{
				var video model.Video
				video.Title = youkuVideo.Title
				if err = db.Where("title = ?",video.Title).First(&video).Error ; err == gorm.ErrRecordNotFound{
					video.Description = youkuVideo.Description
					video.Actors = youkuVideo.Actors
					video.Directors = youkuVideo.Directors
					video.ThumbY = youkuVideo.Thumb
					video.PublishDate = youkuVideo.PublishDate
					video.Pinyin = util.TitleToPinyin(video.Title)
					if err = db.Save(&video).Error ; err != nil{
						logger.Error("query video err!!!,", err)
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}

					var episode model.Episode
					episode.Score = video.Score
					episode.PublishDate = video.PublishDate
					episode.Description = video.Description
					episode.ThumbY = video.ThumbY
					episode.VideoId = video.Id
					episode.Pinyin = util.TitleToPinyin(episode.Title)
					if err = db.Save(&video).Error ; err != nil{
						logger.Error("query video err!!!,", err)
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}

					var playUrl model.PlayUrl
					playUrl.Title = episode.Title
					playUrl.ContentType = constant.MediaTypeEpisode
					playUrl.ContentId = episode.Id
					playUrl.OnLine = true
					playUrl.Provider = youkuVideo.Provider
					playUrl.PageUrl = youkuVideo.PageUrl
					episode.Pinyin = util.TitleToPinyin(episode.Title)
					if err = db.Save(&video).Error ; err != nil{
						logger.Error("query video err!!!,", err)
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}
				}else{

				}
				youkuVideo.Id = video.Id
				apiModels = append(apiModels, &youkuVideo)
			}
		}

		var videos []model.Video
		if err = db.Offset(p.Offset).Limit(p.Limit - count).Order("sort asc").Where("title like ? and on_line = ?", "%" + p.Title + "%", constant.MediaStatusOnLine).Find(&videos).Error; err != nil {
			logger.Error("query video err!!!,", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, video := range videos {
			var apiStream ApiStream
			apiStream.Id = video.Id
			apiStream.Thumb = video.ThumbY
			apiStream.Title = video.Title
			apiStream.ContentType = constant.MediaTypeEpisode
			apiStream.Actors = video.Actors
			apiStream.Country = video.Country
			apiStream.Directors = video.Directors
			apiStream.PublishDate = video.PublishDate
			apiStream.Score = fmt.Sprint(video.Score)
			apiStream.Description = video.Description
			apiStream.Provider = constant.ContentProviderSystem
			apiModels = append(apiModels, &apiStream)
		}
	}


	var hasMore bool = true
	if len(apiModels) < p.Limit {
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiModels, "has_more":hasMore, "count":count})
}

func TopSearchHandler(c *gin.Context) {

	type param struct {
		Limit uint32 `form:"limit" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var tops []model.TopSearch
	if err = db.Order("sort asc").Limit(p.Limit).Find(&tops).Error; err != nil {
		logger.Error("query top_search err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": tops})
}


