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
	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var movies model.Movie
	if err = db.Limit(3).Find(&movies).Error ; err != nil{
		logger.Error("query movie err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": movies})
}

