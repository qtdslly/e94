package api


import (
	"net/http"
	"background/photo/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func PhotoListHandler(c *gin.Context) {

	type param struct {
		Limit             int      `form:"limit" binding:"required"`
		Offset            int      `form:"offset" binding:"exists"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var photoes []model.Photo
	if err := db.Offset(p.Offset).Limit(p.Limit).Where("state = ? and url like '%big%'", model.PhotoStateOnLine).Find(&photoes).Error; err != nil {
		if err != gorm.ErrRecordNotFound{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	var count uint32
	if err := db.Model(&model.Photo{}).Where("state = ? and url like '%big%'", model.PhotoStateOnLine).Count(&count).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var hasMore bool = true
	if len(photoes) < p.Limit{
		hasMore = false
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success,"data":photoes,"count":count,"has_more":hasMore})
}


