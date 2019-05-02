package api

import (
	"background/common/logger"
	"background/newmovie/model"
	"background/common/constant"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func DiggListHandler(c *gin.Context) {
	type param struct {
		ContentType uint8  `form:"content_type" binding:"required"`
		Limit       int `form:"limit" binding:"required"`
		Offset      int `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	installationId := c.MustGet(constant.ContextInstallationId).(uint64)

	var contents []model.ContentAction
	if err := db.Limit(p.Limit).Offset(p.Offset).Where("installation_id = ? AND content_type = ? AND action = 1", installationId, p.ContentType).Find(&contents).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var count uint32
	if err := db.Model(&model.ContentAction{}).Where("installation_id = ? AND content_type = ? AND action = 1", installationId, p.ContentType).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiDiggs struct {
		Id	uint32 `json:"id"`
		StreamId	uint32 `json:"stream_id"`
		Title	string `json:"title"`
		Thumb	string `json:"thumb"`
	}

	var apiDiggs []*ApiDiggs
	for _ , content := range contents{
		var apiDigg ApiDiggs
		apiDigg.Id = content.Id
		apiDigg.StreamId = content.ContentId
		apiDigg.Title = content.Title
		apiDigg.Thumb = content.Thumb
		apiDiggs = append(apiDiggs,&apiDigg)
	}

	var hasMore bool = true
	if len(apiDiggs) < p.Limit{
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiDiggs,"count":count,"has_more":hasMore})
}

func DiggHandler(c *gin.Context) {

	type param struct {
		ContentType      uint8  `json:"content_type"`
		ContentId        uint32 `json:"content_id"`
		Disable          bool   `json:"disable"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	installationId := c.MustGet(constant.ContextInstallationId).(uint64)

	var action model.ContentAction
	action.InstallationId = installationId
	action.ContentType = p.ContentType
	action.ContentId = p.ContentId
	action.Action = uint8(1)
	var stream model.Stream
	if err := db.Where("id = ?",action.ContentId).First(&stream).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	action.Title = stream.Title
	action.Thumb = stream.Thumb
	if !p.Disable {
		if err := db.Where("installation_id = ? AND content_type = ? AND content_id = ? AND action = ?", action.InstallationId, action.ContentType, action.ContentId, action.Action).First(model.ContentAction{}).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			} else {
				if err = db.Create(&action).Error ; err != nil{
					logger.Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			}
		}
	} else {

		if err := db.Where("installation_id = ? AND content_type = ? AND content_id = ? AND action = ?", action.InstallationId, action.ContentType, action.ContentId, action.Action).Delete(model.ContentAction{}).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}

