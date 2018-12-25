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

	type ApiDuty struct{
		Id           uint32         `gorm:"primary_key" json:"id"`
		DoctorId     uint32         `json:"doctor_id"`
		Date         string         `json:"date"`
		Morning      bool           `json:"morning"`
		Afternoon    bool           `json:"afternoon"`
		Night        bool           `json:"night"`
		Week         string         `json:"week"`
	}

	flag := ""
	var apiDutys []*ApiDuty
	for _, duty := range dutys{
		var apiDuty ApiDuty
		apiDuty.Id = duty.Id
		apiDuty.DoctorId = duty.DoctorId
		apiDuty.Date = duty.Date
		apiDuty.Morning = duty.Morning
		apiDuty.Afternoon = duty.Afternoon
		apiDuty.Night = duty.Night

		n , _ := time.Parse("2006-01-02 15:04:05", apiDuty.Date + " 00:00:00")
		logger.Debug(fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",n.Year(),n.Month(),n.Day(),n.Hour(),n.Minute(),n.Second()))
		apiDuty.Week = n.Weekday().String()

		if apiDuty.Week == "Monday"{
			apiDuty.Week = flag + "周一"
		}else if apiDuty.Week == "Tuesday"{
			apiDuty.Week = flag + "周二"
		}else if apiDuty.Week == "Wednesday"{
			apiDuty.Week = flag + "周三"
		}else if apiDuty.Week == "Thursday"{
			apiDuty.Week = flag + "周四"
		}else if apiDuty.Week == "Friday"{
			apiDuty.Week = flag + "周五"
		}else if apiDuty.Week == "Saturday"{
			apiDuty.Week = flag + "周六"
		}else if apiDuty.Week == "Sunday"{
			apiDuty.Week = flag + "周日"
			flag = "下"
		}

		apiDutys = append(apiDutys,&apiDuty)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiDutys})
}