package cms

import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func NewMovieListHandler(c *gin.Context) {

	type param struct {

		Limit       uint32 `form:"limit" binding:"required"`
		Offset      uint32 `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var movies []model.Movie
	if err = db.Order("publish_date desc").Offset(p.Offset).Limit(p.Limit).Find(&movies).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var count uint32
	if err = db.Model(&model.Movie{}).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": movies,"count":count})
}

