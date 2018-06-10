package cms

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"background/common/logger"
	"background/common/constant"
	"background/newmovie/model"
	"background/common/util"
	"errors"
	"net/http"
)

/*
	POST /admin/login
	管理员登录
	@Author:HYK
	http://localhost:2000/#!./ams/ams-admin.md
*/
func AdminLoginHandler(c *gin.Context) {
	type param struct {
		Account  string `form:"account"  json:"account" binding:"required"`  //username 或 mobile 或 email
		Password string `form:"password" json:"password" binding:"required"` //登录密码, password, smscode至少需要一项有值
	}
	type adminInfo struct {
		Id       uint32 `json:"id"`
		Username string `json:"username"`
	}
	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Debug("Invalid request param ", err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var dbAdmin model.Admin

	// add operation log when handler return
	defer func() {
		// not log the password
		c.Set(constant.ContextRequestBody, &param{Account: p.Account})
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextAdmin, &dbAdmin)
		if err != nil {
			c.Set(constant.ContextError, err.Error())
		}
	}()

	// Find the user
	if err = db.Where("(username = ? OR email = ? OR mobile = ?)", p.Account, p.Account, p.Account).First(&dbAdmin).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.AdminNotExists, "err_msg": constant.TranslateErrCode(constant.AdminNotExists)})
		return
	}

	if dbAdmin.Password != util.SHA512(p.Password) {
		err = errors.New("wrong password")
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.WrongUsernamePassword, "err_msg": constant.TranslateErrCode(constant.WrongUsernamePassword)})
		return
	}

	// if admin login successfully, then create an admin token
	var info adminInfo
	info.Id = dbAdmin.Id
	info.Username = dbAdmin.Username
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data": info})
}
