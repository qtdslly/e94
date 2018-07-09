package api

import (
	"net/http"
	"background/newmovie/service"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	apimodel "background/newmovie/controller/api/model"
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
	if err = db.Order("stream.sort asc").Offset(p.Offset).Limit(p.Limit).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ?", p.ResourceGroupId).Find(&streams).Error; err != nil {
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
		apiStream.Thumb = "http:/www.ezhantao.com" + stream.Thumb
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
		Id uint32 `form:"id" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var stream model.Stream
	if err = db.Where("id = ? and on_line = ?", p.Id, constant.MediaStatusOnLine).First(&stream).Error; err != nil {
		logger.Error("query video err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiStream struct {
		Id      uint32                `json:"id"`
		Title   string                `json:"title"`
		Thumb   string                `json:"thumb"`
		PlayUrl []*apimodel.PlayUrl   `json:"play_url"`
	}

	var apiStream ApiStream
	apiStream.Id = stream.Id
	apiStream.Thumb = "http://www.ezhantao.com" + stream.Thumb
	apiStream.Title = stream.Title
	var playUrls []model.PlayUrl
	if err := db.Order("ready asc").Where("content_type = 4 and content_id = ?", stream.Id).Find(&playUrls).Error; err != nil {
		logger.Error("query play_url err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		Area        string                `json:"area"`
		Directors   string                `json:"directors"`
		PageUrl     string                `json:"pageUrl"`
		PublishDate string                `json:"publish_date"`
		Score       string                `json:"score"`
		Description string                `json:"description"`
		Provider    uint8                 `json:"provider"`
	}

	var apiModels []*ApiStream
	for _, stream := range streams {
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Thumb = "http://www.ezhantao.com" + stream.Thumb
		apiStream.Title = stream.Title
		apiStream.ContentType = 4
		apiModels = append(apiModels, &apiStream)
	}

	var count uint32
	if err = db.Model(&model.Stream{}).Where("title like ? and on_line = ?", "%" + p.Title + "%", constant.MediaStatusOnLine).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count == 0{
		if p.Offset == 0{
			err,title,description,actors,directors,thumb,pageUrl,publishDate := service.GetYoukuVideoInfoByTitle(p.Title)
			if err == nil{
				var youkuVideo ApiStream
				youkuVideo.Title = title
				youkuVideo.Description = description
				youkuVideo.Actors = actors
				youkuVideo.Directors = directors
				youkuVideo.Thumb = thumb
				youkuVideo.PageUrl = pageUrl
				youkuVideo.PublishDate = publishDate
				youkuVideo.Provider = constant.ContentProviderYouKu
				apiModels = append(apiModels, &youkuVideo)
				count++
			}

			err,title,description,thumb,score,actors,directors,pageUrl := service.GetTencentVideoInfoByTitle(p.Title)
			if err == nil{
				var tenVideo ApiStream
				tenVideo.Title = title
				tenVideo.Description = description
				tenVideo.Actors = actors
				tenVideo.Directors = directors
				tenVideo.Thumb = thumb
				tenVideo.PageUrl = pageUrl
				tenVideo.Score = score
				tenVideo.Provider = constant.ContentProviderTencent
				apiModels = append(apiModels, &tenVideo)
				count++
			}

			err,title,score,area,description,actors,directors,thumb,pageUrl,publishDate := service.GetIqiyiVideoInfoByTitle(p.Title)
			if err == nil{
				var iqiyiVideo ApiStream
				iqiyiVideo.Title = title
				iqiyiVideo.Description = description
				iqiyiVideo.Actors = actors
				iqiyiVideo.Directors = directors
				iqiyiVideo.Thumb = thumb
				iqiyiVideo.PageUrl = pageUrl
				iqiyiVideo.Score = score
				iqiyiVideo.Area = area
				iqiyiVideo.Provider = constant.ContentProviderIqiyi
				apiModels = append(apiModels, &iqiyiVideo)
				count++
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
			apiStream.ContentType = 2
			apiStream.Actors = video.Actors
			apiStream.Area = video.Country
			apiStream.Directors = video.Directors
			apiStream.PublishDate = video.PublishDate
			apiStream.Score = video.Score
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


