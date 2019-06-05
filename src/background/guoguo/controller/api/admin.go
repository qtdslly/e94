package api
import (
	"background/guoguo/model"
	"common/constant"
	"common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"errors"

	"background/common/util"

	"time"
	"fmt"
)

/*
	POST /admin/login
	管理员登录
	@Author:HYK
	http://localhost:2000/#!./ams/ams-admin.md
*/
func AdminLoginHandler(c *gin.Context) {
	type param struct {
		UserName  string `form:"user_name"  json:"user_name"`  //username 或 mobile 或 email
		Password  string `form:"password" json:"password"` //登录密码, password, smscode至少需要一项有值
	}
	type adminInfo struct {
		Id       uint32 `json:"id"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}

	var p param
	var err error
	if err = c.Bind(&p); err != nil {
		logger.Error("Invalid request param :",c.Params, err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var admin model.Admin

	// add operation log when handler return
	defer func() {
		// not log the password
		c.Set(constant.ContextRequestBody, &param{UserName: p.UserName})
		c.Set(constant.ContextTableName, model.Admin{}.TableName())
		c.Set(constant.ContextAdmin, &admin)
		if err != nil {
			c.Set(constant.ContextError, err.Error())
		}
	}()

	// Find the user
	if err = db.Where("(username = ? OR email = ? OR mobile = ?)", p.UserName, p.UserName, p.UserName).First(&admin).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.AdminNotExists, "err_msg": constant.TranslateErrCode(constant.AdminNotExists)})
		return
	}

	logger.Error(util.SHA512(p.Password))

	if admin.Password != util.SHA512(p.Password) {
		err = errors.New("wrong password")
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{"err_code": constant.WrongUsernamePassword, "err_msg": constant.TranslateErrCode(constant.WrongUsernamePassword)})
		return
	}

	dev, err := createAdminToken(c, admin.Id, db)
	if err != nil {
		logger.Debug("创建token失败")
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := db.Save(&admin).Error; err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	logger.Debug(dev.Token)
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "token": "Bearer " + dev.Token, "account": admin.Username})
}


// create or update user_token
func createAdminToken(c *gin.Context, userId uint32, db *gorm.DB) (*model.AdminToken, error) {
	var token model.AdminToken
	// just check if record found, no need to log
	db.Where("admin_id = ?", userId).First(&token)

	// diable the other user_tokens belonging to this user
	var tokensToDisable []*model.AdminToken
	if err := db.Where("admin_id = ?", userId).Find(&tokensToDisable).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	for _, tokenToDisable := range tokensToDisable {
		tokenToDisable.Disabled = true
		db.Save(tokenToDisable)
	}
	token.AdminId = userId
	token.Disabled = false // set disabled to false
	token.UserAgent = c.Request.UserAgent()
	token.LoginIp = c.ClientIP()
	now := time.Now()
	token.LoginAt = &now
	if token.Id == 0 {
		token.CreatedIp = c.ClientIP()
		if err := db.Create(&token).Error; err != nil {
			logger.Error(err)
			return nil, err
		}
	}
	token.Token = fmt.Sprint(token.Id) + "|" + util.GenerateSessionId(time.Now().String()+c.Request.UserAgent())
	if err := db.Save(&token).Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	return &token, nil
}

