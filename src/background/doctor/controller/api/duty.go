package api

import (
	"net/http"
	"background/doctor/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

func DutyList(c *gin.Context) {

	type param struct {
		DoctorId   uint32 `form:"doctor_id" json:"doctor_id"`
		Limit      int    `form:"limit" json:"limit"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if p.DoctorId == 0{
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	now := time.Now()

	date := fmt.Sprintf("%04d-%02d-%02d",now.Year(),now.Month(),now.Day())

	var dutys []model.Duty

	if err = db.Order("date asc").Limit(p.Limit).Where("doctor_id = ? and date >= ?",p.DoctorId,date).Find(&dutys).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": dutys})
}