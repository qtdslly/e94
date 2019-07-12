package api
import (
	"background/lafter/model"
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

/*
	POST /admin/login
	管理员登录
	@Author:HYK
	http://localhost:2000/#!./ams/ams-admin.md
*/
func LafterHandler(c *gin.Context) {
	type param struct {
		Offset       string `form:"offset" json:"offset"`
		Direct       string `form:"direct" json:"direct"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var content model.Content
	if err := db.Order("id " + p.Direct).Limit(1).Where("id > ?",p.Offset).First(&content).Error ; err != nil{
		if err == gorm.ErrRecordNotFound{
			if err := db.Order("id " + p.Direct).First(&content).Error ; err != nil{
				if err == gorm.ErrRecordNotFound{

				}
				logger.Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}else{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	}

	content.Content = strings.Replace(content.Content,"\r\n","<br/>",-1)

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "content": content})
}

func ZanHandler(c *gin.Context) {
	type param struct {
		Id      uint32 `form:"id" json:"id"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var content model.Content
	if err := db.Where("id = ?",p.Id).First(&content).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	content.Zan += 1
	if err := db.Save(&content).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
