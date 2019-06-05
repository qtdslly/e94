package api

import (
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	cmod "background/cms/model"
)

func EpisodeListHandler(c *gin.Context) {

	type param struct {
		VideoId     uint32 `form:"video_id" json:"video_id"`
		//Limit  int `form:"limit" binding:"required"`
		//Offset int `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var epiosdes []cmod.Episode
	if err = db.Order("sort asc").Where("video_id = ?", p.VideoId).Find(&epiosdes).Error; err != nil {
		logger.Error("query movie err!!!,", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//var hasMore bool = true
	//if len(epiosdes) != p.Limit {
	//	hasMore = false
	//}

	var count uint32
	if err = db.Model(&cmod.Episode{}).Where("video_id = ?", p.VideoId).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": epiosdes, "count":count})
}

func EpisodeAddHandler(c *gin.Context) {
	type param struct {
		VideoId          uint32 `form:"video_id" json:"video_id"`
		Title       string `form:"title" json:"title"`
		Score       float64 `form:"score" json:"score"`
		Number      string  `form:"number" json:"number"`
		Description string `form:"description" json:"description"`
		PublishDate string `form:"publish_date" json:"publish_date"`
		Duration    uint32 `form:"duration" json:"duration"`
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

	var episode cmod.Episode
	if err = db.Where("video_id = ? and number=?", p.VideoId, p.Number).First(&episode).Error; err == nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	episode.Title = p.Title
	episode.Score = p.Score
	episode.Description = p.Description
	episode.Number = p.Number
	episode.PublishDate = p.PublishDate
	episode.Duration = p.Duration * 60
	episode.ThumbX = p.ThumbX
	episode.ThumbY = p.ThumbY
	episode.VideoId = p.VideoId

	if err := db.Create(&episode).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}



func EpisodeDeleteHandler(c *gin.Context) {
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

	if err := db.Delete(cmod.PlayUrl{},"content_type = ? and content_id = ?",constant.MediaTypeEpisode,p.Id).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := db.Delete(cmod.Episode{},"id = ?",p.Id).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

