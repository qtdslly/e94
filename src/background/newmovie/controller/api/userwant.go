package api


import (
	"net/http"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UserWantHandler(c *gin.Context) {

	type param struct {
		InstallationId    uint64   `json:"installation_id"`
		Title             string   `json:"title"`
		Description       string   `json:"description"`
		Email             string   `json:"email"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var want model.UserWant
	want.Title = p.Title
	want.Description = p.Description
	want.Email = p.Email
	want.InstallationId = p.InstallationId
	if err := db.Save(&want).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
