package api

import (
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)


func OpinionHandler(c *gin.Context) {

	type param struct {
		Category         uint8  `json:"category" binding:"required"`
		Description      string `json:"description"`
		Thumb            string `json:"thumb"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	installationId := c.MustGet(constant.ContextInstallationId).(uint64)

	var opinion model.UserOpinion
	opinion.InstallationId = installationId
	opinion.Category = p.Category
	opinion.Description = p.Description
	opinion.Thumb = strings.Replace(p.Thumb,"http://www.ezhantao.com","",-1)
	if err := db.Create(&opinion).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
