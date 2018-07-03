package api

import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	apimodel "background/newmovie/controller/api/model"
)

func StreamListHandler(c *gin.Context) {

	type param struct {
		ResourceGroupId  uint32    `form:"resource_group_id"`
		Limit            int       `form:"limit" binding:"required"`
		Offset           int       `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var streams []*model.Stream
	if err = db.Order("stream.sort asc").Offset(p.Offset).Limit(p.Limit).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ?",p.ResourceGroupId).Find(&streams).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var count uint32
	if err = db.Model(&model.Stream{}).Joins("inner join stream_group where stream.id = stream_group.stream_id and stream_group.resource_group_id = ?",p.ResourceGroupId).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiStream struct {
		Id	uint32 `json:"id"`
		Title	string `json:"title"`
		Thumb	string `json:"thumb"`
	}

	var apiStreams []*ApiStream
	for _ , stream := range streams{
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Title = stream.Title
		apiStream.Thumb = "http:/www.ezhantao.com" + stream.Thumb
		apiStreams = append(apiStreams,&apiStream)
	}

	var hasMore bool = true
	if len(apiStreams) < p.Limit{
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiStreams,"count":count,"has_more":hasMore})
}



func StreamDetailHandler(c *gin.Context) {

	type param struct {
		Id       uint32 `form:"id" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var stream model.Stream
	if err = db.Where("id = ? and on_line = ?",p.Id,constant.MediaStatusOnLine).First(&stream).Error ; err != nil{
		logger.Error("query video err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiStream struct {
		Id	   uint32                `json:"id"`
		Title	   string                `json:"title"`
		Thumb      string                `json:"thumb"`
		PlayUrl    []*apimodel.PlayUrl   `json:"play_url"`
	}

	var apiStream ApiStream
	apiStream.Id = stream.Id
	apiStream.Thumb = "http://www.ezhantao.com" + stream.Thumb
	apiStream.Title = stream.Title
	var playUrls []model.PlayUrl
	if err := db.Where("content_type = 4 and content_id = ?",stream.Id).Find(&playUrls).Error ; err != nil{
		logger.Error("query play_url err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	for _,playUrl := range playUrls{
		var pUrl *apimodel.PlayUrl
		pUrl = apimodel.PlayUrlFromDb(playUrl)
		if pUrl.IsPlay{
			apiStream.PlayUrl = append(apiStream.PlayUrl,pUrl)
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiStream})
}


func StreamSearchHandler(c *gin.Context) {

	type param struct {
		Title       string `form:"title" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var streams []model.Stream
	if err = db.Order("sort asc").Where("title like ? and on_line = ?","%" + p.Title + "%",constant.MediaStatusOnLine).Find(&streams).Error ; err != nil{
		logger.Error("query video err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	type ApiStream struct {
		Id	   uint32                `json:"id"`
		Title	   string                `json:"title"`
		Thumb      string                `json:"thumb"`
		PlayUrl    []*apimodel.PlayUrl   `json:"play_url"`
	}

	var apiStreams []*ApiStream
	for _ , stream := range streams{
		var apiStream ApiStream
		apiStream.Id = stream.Id
		apiStream.Thumb = "http://www.ezhantao.com" + stream.Thumb
		apiStream.Title = stream.Title
		var playUrls []model.PlayUrl
		if err := db.Where("content_type = 4 and content_id = ?",stream.Id).Find(&playUrls).Error ; err != nil{
			logger.Error("query play_url err!!!,",err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _,playUrl := range playUrls{
			var pUrl *apimodel.PlayUrl
			pUrl = apimodel.PlayUrlFromDb(playUrl)
			if pUrl.IsPlay{
				apiStream.PlayUrl = append(apiStream.PlayUrl,pUrl)
			}
		}

		apiStreams = append(apiStreams,&apiStream)
	}


	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiStreams})
}


func StreamTopSearchHandler(c *gin.Context) {

	type param struct {
		Limit       uint32 `form:"limit" binding:"required"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var tops []model.TopSearch
	if err = db.Where("content_type = 4").Limit(p.Limit).Find(&tops).Error ; err != nil{
		logger.Error("query top_search err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": tops})
}