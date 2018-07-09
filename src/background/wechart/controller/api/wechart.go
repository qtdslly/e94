package api

import (
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	/*
	参数	是否必须	描述
	ToUserName	是	接收方帐号（收到的OpenID）
	FromUserName	是	开发者微信号
	CreateTime	是	消息创建时间 （整型）
	MsgType	是	text
	Content	是	回复的消息内容（换行：在content中能够换行，微信客户端就支持换行显示）
	*/


	c.XML(http.StatusOK,gin.H{
		"ToUserName": p.FromUserName,
		"FromUserName": "wx508e9e50a737c414",
		"CreateTime": time.Now().Unix(),
		"MsgType": "text",
		"Content": "SUCCESS:" + p.Content,
	})

}
