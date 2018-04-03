package middleware

import (
	"background/common/constant"
	"background/common/logger"
	"background/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func AdminVerifyHandler(amsBindAddr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errCode int
		var scope string

		defer func() {
			if errCode != constant.Success {
				c.JSON(http.StatusForbidden, gin.H{"err_code": errCode, "err_msg": constant.TranslateErrCode(errCode)})
				c.Abort()
				return
			}
			c.Set(constant.ContextScopeBody, scope)
			c.Next()
		}()

		url := c.Request.URL.Path
		accessKey := c.Request.Header.Get(constant.AccessKey)
		if accessKey != "" {
			c.Set(constant.ContextModuleAccess, true)

			// check if module access garentted for this access_key
			errCode, scope = verifyAccessKey(accessKey, url, amsBindAddr)
		} else {
			c.Set(constant.ContextModuleAccess, false)

			var token string

			if t := c.Request.Header.Get("Authorization"); t != "" {
				token = t
			} else {
				token, _ = c.GetQuery("token")
			}

			if token == "" {
				errCode = constant.AdminNotLogin
				return
			}

			var adminId uint32
			var username string
			// check if module access garentted for admin
			adminId, username, scope, errCode = verifyAdmin(token, url, amsBindAddr)
			if errCode != constant.Success {
				return
			}

			c.Set(constant.ContextAdminId, adminId)
			c.Set(constant.ContextAdminName, username)
			c.Set(constant.ContextScopeBody, scope)
		}
	}
}

/*
	This method allows to call ams server for admin module access validation.
*/
func verifyAdmin(token, url, amsBindAddr string) (uint32, string, string, int) {
	apiUrl := amsBindAddr + "/ams/admin/verify"

	type param struct {
		Token     string `json:"token"`
		Url       string `json:"url"`
		Timestamp int64  `json:"timestamp"`
	}
	request := &param{}
	request.Token = token
	request.Url = url
	request.Timestamp = time.Now().Unix()

	type response struct {
		ErrCode int    `json:"err_code"`
		ErrMsg  string `json:"err_msg"`
		Data    *struct {
			Id       uint32 `json:"id"`
			Username string `json:"username"`
			Scope    string `json:"scope"`
		} `json:"data"` //admin_id
	}
	var resp response

	paramJson, _ := json.Marshal(request)
	paramMap := make(map[string]interface{})
	json.Unmarshal(paramJson, &paramMap)

	moduleSign := MakeModuleSignature(paramMap, constant.ModuleSalt, request.Timestamp)

	err := util.Post(apiUrl, request, &resp, map[string]string{constant.ModuleSignature: moduleSign})
	if err != nil {
		logger.Error(err)
		return 0, "", "", constant.StatusInternalServerError
	}

	if resp.ErrCode != 0 {
		err = errors.New(fmt.Sprintf("invalid err_code [%d] - [%s]", resp.ErrCode, resp.ErrMsg))
		logger.Error(err)
		return 0, "", "", resp.ErrCode
	}
	return resp.Data.Id, resp.Data.Username, resp.Data.Scope, constant.Success
}

/*
	This method allows to call ams server for admin module access validation.
*/
func verifyAccessKey(accessKey, url, amsBindAddr string) (int, string) {
	apiUrl := amsBindAddr + "/ams/accesskey/verify"
	type param struct {
		AccessKey string `json:"access_key"`
		Url       string `json:"url"`
		Timestamp int64  `json:"timestamp"`
	}
	request := &param{}
	request.AccessKey = accessKey
	request.Url = url
	request.Timestamp = time.Now().Unix()

	type response struct {
		ErrCode int    `json:"err_code"`
		ErrMsg  string `json:"err_msg"`
		Data    string `json:"data"`
	}
	var resp response

	paramJson, _ := json.Marshal(request)
	paramMap := make(map[string]interface{})
	json.Unmarshal(paramJson, &paramMap)

	moduleSign := MakeModuleSignature(paramMap, constant.ModuleSalt, request.Timestamp)

	err := util.Post(apiUrl, request, &resp, map[string]string{constant.ModuleSignature: moduleSign})
	if err != nil {
		logger.Error(err)
		return constant.StatusInternalServerError, ""
	}

	if resp.ErrCode != 0 {
		err = errors.New(fmt.Sprintf("invalid err_code [%d] - [%s]", resp.ErrCode, resp.ErrMsg))
		logger.Error(err)
		return resp.ErrCode, ""
	}
	return constant.Success, resp.Data
}

func PermissionVerifyHandler(c *gin.Context) {
	c.Next()
	//url := c.Request.URL.Path
	//var module model.Module
	//if err := db.Where("url = ?", url).First(&module).Error; err != nil {
	//	logger.Error(err)
	//	c.JSON(http.StatusOK, gin.H{"err_code": constant.ApiNotRegisted, "err_msg": constant.TranslateErrCode(constant.ApiNotRegisted)})
	//	return
	//}
	//
	//admin := c.MustGet(constant.ContextAdmin).(model.Admin)
	//allowed, err := service.IsModuleAllowed(&module, &admin, db)
	//if err != nil {
	//	logger.Error(err)
	//	c.AbortWithStatus(http.StatusBadRequest)
	//	return
	//}
	//
	//if !allowed {
	//	c.JSON(http.StatusOK, gin.H{"err_code": constant.ModuleAccessDenied, "err_msg": constant.TranslateErrCode(constant.ModuleAccessDenied)})
	//	return
	//}
}
