package middleware

import (
	"background/common/constant"
	"background/common/logger"
	"background/common/util"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SignatureVerifyHandler(c *gin.Context) {
	errCode := constant.Success
	var timestamp int64
	var sign string
	params := map[string]interface{}{}

	defer func() {
		if errCode != constant.Success {
			logger.Warn("[SignatureVerifyFailed] URL: ", c.Request.URL.String(), " | ", timestamp, " | ", sign, " | ", params)
			// c.JSON(http.StatusBadRequest, gin.H{"err_code": errCode})
			// c.Abort()
			// return
		}
		c.Next()
	}()

	// 读取json参数, 获取timestamp以及表单
	params, err := ParseParam(c)
	if err != nil {
		logger.Error(err)
		return
	}

	if t, ok := params["timestamp"]; ok {
		if ts, k := t.(float64); k {
			timestamp = int64(ts)
		} else if ts, k := t.(string); k {
			timestamp, _ = strconv.ParseInt(ts, 10, 64)
		}
	}
	if s, ok := params["sign"]; ok {
		sign, _ = s.(string)
	}

	appKey, _ := params["app_key"]
	appVersion, _ := params["app_version"]
	installationId, _ := params["installation_id"]

	c.Set(constant.ContextAppKey, appKey)
	c.Set(constant.ContextAppVersion, appVersion)
	if c.Request.Method == "POST" {
		c.Set(constant.ContextInstallationId, uint64(getInt64(installationId, "POST")))
	}else {
		c.Set(constant.ContextInstallationId, uint64(getInt64(installationId, "")))
	}


	// logger.Debug(fmt.Sprintf("-------------------signature param: app_key:[%+v] os_type:[%+v] app_version:[%+v] installation_id:[%+v]", appKey, osType, appVersion, installationId))

	//比较timestamp是否过期
	if time.Now().Unix()-timestamp > 10*60 {
		logger.Error("Timestamp expired: ", time.Now().Unix()-timestamp)
		errCode = constant.Failure
		return
	}

	// 根据参数（除了sign参数）生成签名字符串
	signCalc := MakeSignature(params, timestamp)

	// 对比sign参数
	if sign == "" || sign != signCalc {
		errCode = constant.Failure
		logger.Error("Signature verify failed! ", signCalc, " | ", c.Request.URL.String())
		return
	}
	// logger.Debug("Signature verify successed! ", signCalc, " | ", c.Request.URL.String())
}

func ModuleSignatureVerifyHandler(c *gin.Context) {
	errCode := constant.Success
	var timestamp int64
	var sign string
	params := map[string]interface{}{}

	defer func() {
		if errCode != constant.Success {
			logger.Warn("[ModuleSignatureVerifyFailed] URL: ", c.Request.URL.String(), " | ", timestamp, " | ", sign, " | ", params)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Next()
	}()

	// 读取json参数, 获取timestamp以及表单
	params, err := ParseParam(c)
	if err != nil {
		logger.Error(err)
		return
	}

	//logger.Debug("verify module signature timestamp：", params["timestamp"])

	if t, ok := params["timestamp"]; ok {
		if ts, k := t.(float64); k {
			timestamp = int64(ts)
		} else if ts, k := t.(string); k {
			timestamp, _ = strconv.ParseInt(ts, 10, 64)
		}
	}
	sign = c.Request.Header.Get(constant.ModuleSignature)

	//比较timestamp是否过期
	if time.Now().Unix()-timestamp > 10*60 {
		logger.Error("Signature verify failed! Timestamp is expired: ", time.Now().Unix()-timestamp)
		errCode = constant.Failure
		return
	}

	// 根据参数（除了sign参数）生成签名字符串
	signCalc := MakeModuleSignature(params, constant.ModuleSalt, timestamp)

	// 对比sign参数
	if sign == "" || sign != signCalc {
		errCode = constant.Failure
		logger.Error("Module Signature verify failed! ", signCalc, " | ", c.Request.URL.String())
		return
	}
	//logger.Debug("Module Signature verify successed! ", signCalc, " | ", c.Request.URL.String())
}

// timestamp + params + salt 进行sha1加密后hex编码
func MakeModuleSignature(params map[string]interface{}, salt string, timestamp int64) string {
	p := strconv.FormatInt(timestamp, 10)
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v, ok := params[k]; ok {
			switch t := v.(type) {
			case int, int32, int64, uint, uint32, uint64:
				p += url.QueryEscape(fmt.Sprintf("%d", t))
			case float32, float64:
				p += url.QueryEscape(fmt.Sprintf("%.0f", t))
			case string:
				p += url.QueryEscape(fmt.Sprintf("%s", t))
			case bool:
				p += url.QueryEscape(fmt.Sprintf("%t", t))
			}
		}
	}
	p += salt
	//logger.Debug("make module sign str：", p)
	hash := sha1.New()
	return hex.EncodeToString(hash.Sum([]byte(p)))
}

// 将输入参数组装成map, 只保留值为字符串和数字类型的键值。按key进行排序后组装成字符串进行md5加密
func MakeSignature(params map[string]interface{}, timestamp int64) string {
	var p []string
	for k, v := range params {
		if k == "sign" {
			continue
		}
		if _, ok := v.(float64); ok {
			p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%.0f", v)))
		} else if _, ok := v.(string); ok {
			p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%v", v)))
		} else if arr, ok := v.([]string); ok {
			sort.Strings(arr)
			for _, arrv := range arr {
				p = append(p, k+"="+url.QueryEscape(fmt.Sprintf("%v", arrv)))
			}
		}
	}
	par, _ := url.ParseQuery(strings.Join(p, "&"))
	source := par.Encode()
	source, _ = url.QueryUnescape(source)

	// 计算出密钥并加密
	//t := uint64(timestamp)
	//n := t % 64
	//secret := fmt.Sprintf("%d", ((t<<n)&(0x7fffffffffffffff))|(t>>(64-n)))
	secret := "5a61efdc52411a670b9f7c9db0a5275b"
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(source))
	// logger.Debug("signature calculated: " + source + " | " + secret)

	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

func ParseParam(c *gin.Context) (map[string]interface{}, error) {
	var err error
	params := map[string]interface{}{}
	params["app_key"] = c.GetHeader("app_key")
	params["app_version"] = c.GetHeader("app_version")
	params["installation_id"] = c.GetHeader("installation_id")

	logger.Debug("app_key:",params["app_key"])
	logger.Debug("app_version:",params["app_version"])
	logger.Debug("installation_id:",params["installation_id"])

	if ctype := c.Request.Header.Get("Content-Type"); len(ctype) > 0 && strings.Contains(ctype, "application/json") {
		resp, _ := ioutil.ReadAll(c.Request.Body)
		buff := util.BuffReadWriter{}
		buff.Write(resp)
		c.Request.Body = &buff
		if err = json.Unmarshal(resp, &params); err == nil {
			params["app_key"] = c.GetHeader("app_key")
			params["app_version"] = c.GetHeader("app_version")
			params["installation_id"] = c.GetHeader("installation_id")
			return params, err
		} else {
			logger.Error(err)
		}
	}

	postLen := 0
	if c.Request.Method == "POST" {
		err = c.Request.ParseForm()
		if err == nil {
			for k, v := range c.Request.PostForm {
				if len(v) == 1 {
					params[k] = v[0]
				} else {
					params[k] = v
				}
				postLen++
			}
		} else {
			logger.Error(err)
		}
	}
	if c.Request.Method == "GET" || postLen <= 0 {
		for k, v := range c.Request.URL.Query() {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}
	}
	return params, nil
}

func getInt(str interface{}) int {
	val, _ := strconv.Atoi(fmt.Sprint(str))
	return val
}

func getInt64(str interface{}, method string) int64 {
	var old string
	if method == "POST" {
		old = fmt.Sprintf("%g", str)
		var newVal float64
		newVal,_ = strconv.ParseFloat(old, 64)
		s := strconv.FormatFloat(newVal, 'f', -1, 64)
		i, _ := strconv.ParseInt(s, 10, 64)
		return i
	}else {
		old = fmt.Sprintf("%s", str)
		i, _ := strconv.ParseInt(old, 10, 64)
		return i
	}
}