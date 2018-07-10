package apimodel

import "time"

type Video struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ContentType uint8     `json:"content_type"`
	Score       string    `json:"score"`
	ThumbX      string    `json:"thumb_x"`
	ThumbY      string    `json:"thumb_y"`
	PublishDate string    `json:"publish_date"`
	Year        string    `json:"year"`
	Language    string    `json:"language"`
	Country     string    `json:"country"`
	Directors   string    `json:"directors"`
	Actors      string    `json:"actors"`
	Tags        string    `json:"tags"`
	PageUrl     string    `json:"page_url"`
	Provider    uint32    `json:"provider"`
}

func VideoToNews(toUserName,fromUserName string,src *Video) *NewsResMessage {
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
	/*
	<Articles>
		<item>
			<Title>< ![CDATA[title1] ]></Title>
			<Description>< ![CDATA[description1] ]></Description>
			<PicUrl>< ![CDATA[picurl] ]></PicUrl>
			<Url>< ![CDATA[url] ]></Url>
		</item>
		<item>
			<Title>< ![CDATA[title] ]></Title>
			<Description>< ![CDATA[description] ]></Description>
			<PicUrl>< ![CDATA[picurl] ]></PicUrl>
			<Url>< ![CDATA[url] ]></Url>
		</item>
	</Articles>
	*/
	var news NewsResMessage
	var article Article
	article.Title = src.Title
	article.Description = src.Description
	article.PicUrl = src.ThumbY //"http://image14.m1905.cn/uploadfile/2017/0731/thumb_1_147_100_20170731013311381963.jpg" //src.ThumbY
	article.Url = "http://www.ezhantao.com/play?url=" + src.PageUrl
	var item Item
	item.Art = article
	news.CreateTime = time.Now().Unix()
	news.ArticleCount = 1
	news.ToUserName = toUserName
	news.FromUserName = fromUserName
	news.MsgType = "news"
	news.Articles = append(news.Articles,item)
	return &news
}
