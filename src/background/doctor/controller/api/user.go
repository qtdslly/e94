package api

import (
	"net/http"
	"background/doctor/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddUser(c *gin.Context) {

	type param struct {
		OpenId   string `form:"open_id" json:"open_id"`
		Avtar    string `form:"avtar" json:"avtar"`
		Nick     string `form:"nick" json:"nick"`
		Country  string `form:"country" json:"country"`
		Province string `form:"province" json:"province"`
		City     string `form:"city" json:"city"`
		Language string `form:"language" json:"language"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var user model.User
	if err = db.Where("open_id = ?", p.OpenId).First(&user).Error; err == gorm.ErrRecordNotFound {
		if err = db.Create(&user).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": user})
}