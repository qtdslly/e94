package api

import (
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	cmod "background/cms/model"
	"background/newmovie/model"
)

func VideoListHandler(c *gin.Context) {

	//type param struct {
	//	Limit       int `form:"limit" binding:"required"`
	//	Offset      int `form:"offset" binding:"exists"`
	//}
	//
	//var p param
	//if err := c.Bind(&p); err != nil {
	//	logger.Error(err)
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var videos []cmod.Video
	if err = db.Order("publish_date desc").Where("on_line = ?", constant.MediaStatusOnLine).Find(&videos).Error; err != nil {
		//if err = db.Where("on_line = ?",constant.MediaStatusOnLine).Find(&videos).Error ; err != nil{
		logger.Error("query movie err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//var hasMore bool = true
	//if len(videos) != p.Limit{
	//	hasMore = false
	//}

	var count uint32
	if err = db.Model(&cmod.Video{}).Where("on_line = ?", constant.MediaStatusOnLine).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": videos, "count":count})
}

func VideoAddHandler(c *gin.Context) {

	type param struct {
		Title       string `form:"title" json:"title"`
		Category    string `form:"category" json:"category"`
		Score       float64 `form:"score" json:"score"`
		Tags        string `form:"tags" json:"tags"`
		Description string `form:"description" json:"description"`
		Directors   string `form:"directors" json:"directors"`
		Actors      string `form:"actors" json:"actors"`
		PublishDate string `form:"publish_date" json:"publish_date"`
		Language    string `form:"language" json:"language"`
		Country     string `form:"country" json:"country"`
		ThumbX      string `form:"thumb_x" json:"thumb_x"`
		ThumbY      string `form:"thumb_y" json:"thumb_y"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var video cmod.Video
	if err = db.Where("title = ?", p.Title).First(&video).Error; err == nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	video.Title = p.Title
	video.Category = p.Category
	video.Score = p.Score
	video.Tags = p.Tags
	video.Description = p.Description
	video.Directors = p.Directors
	video.Actors = p.Actors
	video.PublishDate = p.PublishDate
	video.Country = p.Country
	video.ThumbX = p.ThumbX
	video.ThumbY = p.ThumbY

	if err := db.Create(&video).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}



func VideoDeleteHandler(c *gin.Context) {
	type param struct {
		Id          uint32         `form:"id" json:"id"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var episodes []model.Episode
	if err := db.Where("video_id = ?",p.Id).Find(&episodes).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, episode := range episodes{
		if err := db.Delete(cmod.PlayUrl{},"content_type = ? and content_id = ?",constant.MediaTypeEpisode,episode.Id).Error ; err != nil{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := db.Delete(cmod.Episode{},"id = ?",episode.Id).Error ; err != nil{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if err := db.Delete(cmod.Video{},"id = ?",p.Id).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

