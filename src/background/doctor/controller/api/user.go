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
		Gender   string `form:"gender" json:"gender"`
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
		user.OpenId = p.OpenId
		user.Avtar = p.Avtar
		user.Nick = p.Nick
		user.Country = p.Country
		user.Province = p.Province
		user.City = p.City
		user.Language = p.Language
		if p.Gender == 1{
			p.Gender = "男"
		}else{
			p.Gender = "女"
		}
		user.Gender = p.Gender
		if err = db.Create(&user).Error; err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": user})
}