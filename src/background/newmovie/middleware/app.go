package middleware

import (
	"background/newmovie/controller/api/cache"
	"background/newmovie/model"
	"background/common/constant"
	"background/common/logger"
	"background/common/middleware"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
	middleware for app verify
*/
func AppVerifyHandler(appType uint32) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet(constant.ContextDb).(*gorm.DB)

		var err error
		params := map[string]interface{}{}

		defer func() {
			if err != nil {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}()

		// 读取json参数, 获取timestamp以及表单
		params, err = middleware.ParseParam(c)
		if err != nil {
			logger.Error(err)
			return
		}

		appKey, ok := params["app_key"]
		if !ok || appKey == ""{
			err = errors.New("app_key is not exists")
			logger.Error(err)
			return
		}
		appVersion, ok := params["app_version"]
		if !ok || appVersion == ""{
			err = errors.New("app_version is not exists")
			logger.Error(err)
			return
		}

		// 2017-11-14：新增参数installation_id
		installationId, ok := params["installation_id"]
		if !ok {
			// 兼容旧版，此处不报错，默认为0
			installationId = uint64(0)
		} else {
			tempId, _ := strconv.ParseInt(fmt.Sprint(installationId), 10, 64)
			installationId = uint64(tempId)
		}

		// 获取app
		var version *model.Version
		version, err = cache.GetVersion(appKey.(string), appVersion.(string), db)
		if err != nil {
			logger.Error(err)
			return
		}

		var app *model.App
		app, err = cache.GetApp(version.AppId, version.AppKey, db)
		if err != nil {
			logger.Error(err)
			return
		}

		// 2017-11-07：新增app type校验
		if app.Type != appType {
			err = errors.New(`invalid app type`)
			logger.Error(err)
			return
		}

		c.Set(constant.ContextAppVersion,version)
		c.Set(constant.ContextAppId, version.AppId)
		c.Set(constant.ContextVersionId, version.Id)
	}
}

func getOsType(typeStr interface{}) uint32 {
	typ, _ := strconv.Atoi(fmt.Sprint(typeStr))
	return uint32(typ)
}
