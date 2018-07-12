package controller

import (
	"background/qrcode/logger"
	"background/qrcode/util"
	"image/jpeg"
	"net/http"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

// 根据运营商ID跟机顶盒ID返回二维码图片，其二维码内容为一段JSON串
func IptvQrcodeHandler(c *gin.Context) {
	type param struct {
		Message           string `form:"msg" json:"msg"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error("Invalid device register param ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	logger.Debug(c.Request.URL)

	var err error
	var q *qrcode.QRCode
	if q, err = qrcode.New(p.Message, qrcode.Highest); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	buff := util.BuffReadWriter{}
	if err = jpeg.Encode(&buff, q.Image(430), &jpeg.Options{jpeg.DefaultQuality}); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	logger.Info("[IptvQrcode] Resposned qrcode image with ", p.Message)

	c.Data(http.StatusOK, "image/jpeg", buff.Buff)
}
