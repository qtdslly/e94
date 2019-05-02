package api

import (
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	cmod "background/cms/model"

)
func VideoListHandler(c *gin.Context) {

	type param struct {
		Limit       int `form:"limit" binding:"required"`
		Offset      int `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var videos []cmod.Video
	if err = db.Order("publish_date desc").Offset(p.Offset).Where("on_line = ?",constant.MediaStatusOnLine).Limit(p.Limit).Find(&videos).Error ; err != nil{
	//if err = db.Where("on_line = ?",constant.MediaStatusOnLine).Find(&videos).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var hasMore bool = true
	if len(videos) != p.Limit{
		hasMore = false
	}

	var count uint32
	if err = db.Model(&cmod.Video{}).Where("on_line = ?",constant.MediaStatusOnLine).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": videos,"count":count,"has_more":hasMore})
}


