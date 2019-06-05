package api
import (
	"background/lafter/model"
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
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
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "content": content})
}
