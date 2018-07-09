package api

import (
	"background/common/constant"
	"background/common/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/PuerkitoBio/goquery"
)


func WebVideoHandler(c *gin.Context) {
	type param struct {
		Url   string `form:"url"`
	}

	var p param
	if err := c.Bind(&p); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(p.Url) == 0{
		logger.Error("Url为空")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	apiUrl := "http://jx.618g.com/?url=" + p.Url
	//apiUrl := "http://youkus.baiyug.cn:3670/vip_all/baiyug.php?xml=https://tv.sohu.com/v/dXMvMzg4MDcyOS8yOTY2MDMzOS5zaHRtbA==.html&md5=ab59a6c4e08ab5d681f959c50cfaloij&type=sohu&hd=cq&wap=0&siteuser=&lg=&sohuuid=98E442AE12748F007B6D39A6AA43B7AF"
	//https://tv.sohu.com/v/dXMvMzg4MDcyOS8yOTY2MDMzOS5zaHRtbA==.html
	query, err := goquery.NewDocument(apiUrl)
	if err != nil {
		logger.Debug(apiUrl)
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	base,exist := query.Find("iframe").Eq(0).Attr("src")
	if !exist{
		logger.Debug(apiUrl)
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
	url := base[strings.Index(base,"url=") + 4:]
	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success,"url":url})
}
