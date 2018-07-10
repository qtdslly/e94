package apimodel

import (
	"encoding/xml"
)

type TextResMessage struct {
	XMLName       xml.Name    `xml:"xml"`
	ToUserName    string      `xml:"ToUserName"`
	FromUserName  string      `xml:"FromUserName"`
	CreateTime    int64       `xml:"CreateTime"`
	MsgType       string      `xml:"MsgType"`
	Content       string      `xml:"Content"`
}

type Article struct {
	Title         string      `xml:"Title"`
	Description   string      `xml:"Description"`
	PicUrl        string      `xml:"PicUrl"`
	Url           string      `xml:"Url"`
}

type Item struct {
	Art     Article     `xml:"item"`
}
type NewsResMessage struct {
	/*
	参数	是否必须	说明
	ToUserName	是	接收方帐号（收到的OpenID）
	FromUserName	是	开发者微信号
	CreateTime	是	消息创建时间 （整型）
	MsgType	是	news
	ArticleCount	是	图文消息个数，限制为8条以内
	Articles	是	多条图文消息信息，默认第一个item为大图,注意，如果图文数超过8，则将会无响应
	Title	是	图文消息标题
	Description	是	图文消息描述
	PicUrl	是	图片链接，支持JPG、PNG格式，较好的效果为大图360*200，小图200*200
	Url	是	点击图文消息跳转链接
	*/
	XMLName       xml.Name    `xml:"xml"`
	ToUserName    string      `xml:"ToUserName"`
	FromUserName  string      `xml:"FromUserName"`
	CreateTime    int64       `xml:"CreateTime"`
	MsgType       string      `xml:"MsgType"`
	ArticleCount  int         `xml:"ArticleCount"`
	Articles      []Item      `xml:"Articles"`
}