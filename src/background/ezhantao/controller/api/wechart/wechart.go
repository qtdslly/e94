package api

import (
	"background/ezhantao/service"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/xml"
)

func WeChartHandler(c *gin.Context) {
	/*
	参数	描述
	ToUserName	开发者微信号
	FromUserName	发送方帐号（一个OpenID）
	CreateTime	消息创建时间 （整型）
	MsgType	text
	Content	文本消息内容
	MsgId	消息id，64位整型
	*/
	type param struct {
		ToUserName   string `xml:"ToUserName"`
		FromUserName string `xml:"FromUserName"`
		CreateTime   int64 `xml:"CreateTime"`
		MsgType      string `xml:"MsgType"`
		Content      string `xml:"Content"`
		MsgId        int64 `xml:"MsgId"`
		Event        string `xml:"Event"`

	}
	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(p.ToUserName)
	logger.Debug(p.FromUserName)
	logger.Debug(p.CreateTime)
	logger.Debug(p.MsgType)
	logger.Debug(p.Content)
	logger.Debug(p.MsgId)
	logger.Debug(p.Event)


	var data []byte
	if p.MsgType == "text"{
		result := service.SearchGoods(p.FromUserName,p.ToUserName,p.Content)
    data, _ = xml.MarshalIndent(result, "", "  ")
  }else if p.MsgType == "event"{
		wx := service.EventContent(p.Event,p.FromUserName,p.ToUserName)
		data, _ = xml.MarshalIndent(wx, "", "  ")
	}

	c.String(http.StatusOK,string(data))

}
