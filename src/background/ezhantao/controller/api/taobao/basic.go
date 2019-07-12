package taobao

import(
  "github.com/gin-gonic/gin"
  "background/common/logger"

  "background/ezhantao/service"
  "net/http"
)

func GoodsListHandler(c *gin.Context) {
  type param struct {
    Page        int    `form:"page" json:"page"`
    Size        int    `form:"size" json:"size"`
    Category    string `form:"category" json:"category"`
  }

  var p param
  if err := c.Bind(&p); err != nil {
    logger.Error(err)
    return
  }

  recv := service.GetHaoQuanList(p.Page,p.Size,"爆款"+p.Category)
  logger.Debug(recv)
  c.String(http.StatusOK,recv)
}

func TPwdHandler(c *gin.Context) {
  type param struct {
    Title        string    `form:"title" json:"title"`
    Logo         string    `form:"logo" json:"logo"`
    Url          string    `form:"url" json:"url"`
  }

  var p param
  if err := c.Bind(&p); err != nil {
    logger.Error(err)
    return
  }

  recv := service.GetTPwd(p.Url,p.Title,p.Logo)
  logger.Debug(recv)
  c.String(http.StatusOK,recv)
}

func SearchHandler(c *gin.Context) {
  type param struct {
    Page        int    `form:"page" json:"page"`
    Size        int    `form:"size" json:"size"`
    Category    string `form:"category" json:"category"`
  }

  var p param
  if err := c.Bind(&p); err != nil {
    logger.Error(err)
    return
  }

  recv := service.GetHaoQuanList(p.Page,p.Size,p.Category)
  logger.Debug(recv)
  c.String(http.StatusOK,recv)
}
