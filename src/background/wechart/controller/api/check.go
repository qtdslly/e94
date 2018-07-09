package api

import (
	"background/common/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
	"sort"
	"crypto/sha1"
	"encoding/hex"
)

func CheckHandler(c *gin.Context) {
	/*
	参数	        描述
	signature	微信加密签名，signature结合了开发者填写的token参数和请求中的timestamp参数、nonce参数。
	timestamp	时间戳
	nonce	        随机数
	echostr	        随机字符串
	*/
	type param struct {
		Signature   string `form:"signature"`
		Timestamp   int64  `form:"timestamp"`
		Nonce       string `form:"nonce"`
		Echostr     string `form:"echostr"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	token := "5a61efdc52411a670b9f7c9db0a5275b"
	tmpStr := []string{fmt.Sprint(p.Timestamp),p.Nonce,token}
	sort.Strings(tmpStr)

	sortStr := tmpStr[0] + tmpStr[1] + tmpStr[2]

	hash := sha1.New()
	result := hex.EncodeToString(hash.Sum([]byte(sortStr)))

	if result != p.Signature{
		logger.Debug("验证失败")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK,p.Echostr)
}
