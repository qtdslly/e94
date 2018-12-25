package api

import (
	"net/http"
	"background/doctor/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"fmt"
)

func CommentList(c *gin.Context) {
	type param struct {
		DoctorId   uint32 `form:"doctor_id" json:"doctor_id"`
		Offset     int    `form:"offset" json:"offset"`
		Limit      int    `form:"limit" json:"limit"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	if p.DoctorId == 0{
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var doctor model.Doctor
	if err = db.Where("id = ?", p.DoctorId).First(&doctor).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type ApiComment struct {
		UserId   uint32 `json:"user_id"`
		Nick     string `json:"nick"`
		Avtar    string `json:"avtar"`
		Content   string `json:"content"`
		DoctorNick   string `json:"doctor_nick"`
		Reply    string `json:"reply"`
		CreatedAt   string `json:"created_at"`
	}

	var comments []model.Comment
	if err = db.Order("created_at desc").Limit(20).Where(" state = ? and doctor_id = ?",model.CommentStatePublish,p.DoctorId).Find(&comments).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var apiComments []*ApiComment
	for _ , comment := range comments{
		var user model.User
		if err = db.Where("id = ?",comment.UserId).First(&user).Error ; err != nil{
			logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var apiComment ApiComment
		apiComment.UserId = comment.UserId
		apiComment.Content = comment.Content
		apiComment.DoctorNick = doctor.Nick
		apiComment.Avtar = user.Avtar
		apiComment.Nick = user.Nick
		apiComment.Reply = comment.Reply

		var n = comment.CreatedAt
		apiComment.CreatedAt = fmt.Sprintf("%04d-%02d-%02d %02d:%02d%02d",n.Year(),n.Month(),n.Day(),n.Hour(),n.Minute(),n.Second())

		apiComments = append(apiComments,&apiComment)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": apiComments})
}



func CommentAdd(c *gin.Context) {
	type param struct {
		UserId     uint32 `form:"user_id" json:"user_id"`
		DoctorId   uint32 `form:"doctor_id" json:"doctor_id"`
		Content    string `form:"content" json:"content"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var err error

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var comment model.Comment
	comment.UserId = p.UserId
	comment.DoctorId = p.DoctorId
	comment.Content = p.Content
	comment.State = model.CommentStateUnknow

	if err = db.Create(&comment).Error ; err != nil{
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": comment})
}