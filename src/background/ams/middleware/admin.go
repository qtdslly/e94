package middleware

import (
	"ams/model"
	"common/constant"
	"common/logger"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"common/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

/*
	middleware for user login
*/
func AdminVerifyHandler(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			//c.AbortWithStatus(http.StatusForbidden)
			//c.Abort()
			c.JSON(http.StatusOK, gin.H{"err_code": constant.Failure, "err_msg": "invalid token"})

			return
		}
		c.Next()
	}()

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	db.LogMode(true)

	params, err := middleware.ParseParam(c)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(len(params))
	logger.Debug(params)

	bearerTokenStr, ok := params["Authorization"]
	logger.Debug(bearerTokenStr.(string))
	if !ok || !strings.HasPrefix(bearerTokenStr.(string), "Bearer ") {
		err = errors.New("token not starting with [Bearer ]")
		logger.Error(err)
		return
	}

	tokenStr := strings.TrimPrefix(bearerTokenStr.(string), "Bearer ")
	if tokenStr == "" {
		err = errors.New(fmt.Sprintf("invalid token [%s]", tokenStr))
		logger.Error(err)
		return
	}

	// tokenStr format "token_id|random_string"
	tokenSplit := strings.Split(tokenStr, "|")
	if len(tokenSplit) < 2 {
		err = errors.New(fmt.Sprintf("invalid token [%s]", tokenStr))
		logger.Error(err)
		return
	}

	var token *model.AdminToken
	tokenId, err := strconv.Atoi(tokenSplit[0])
	if err != nil {
		logger.Error(err)
		return
	}

	token = &model.AdminToken{}
	db.Where("id=?", tokenId).First(token)
	if token.Id == 0 {
		err = errors.New(fmt.Sprintf("invalid token [%s]", tokenStr))
		logger.Error(err)
		return
	}

	if token.Token != tokenStr {
		err = errors.New(fmt.Sprintf("invalid token [%s]", tokenStr))
		logger.Error(err)
		return
	}

	// check if token is disabled
	if token.Disabled {
		err = errors.New(fmt.Sprintf("disabled token [%s]", tokenStr))
		logger.Error(err)
		return
	}

	admin := &model.Admin{}
	db.Where("id=?", token.AdminId).First(admin)
	if admin.Id == 0 {
		err = errors.New(fmt.Sprintf("invalid token [%s]", tokenStr))
		logger.Error(err)
		return
	}

	c.Set(constant.ContextAdmin, admin)
	c.Set(constant.ContextToken, tokenStr)
}
