package api

import (
  "background/guoguo/model"
  "common/constant"
  "common/logger"
  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  "net/http"
)

/*
	POST /admin/login
	管理员登录
	@Author:HYK
	http://localhost:2000/#!./ams/ams-admin.md
*/
func NoticeListHandler(c *gin.Context) {

  db := c.MustGet(constant.ContextDb).(*gorm.DB)

  var Stocks []model.Notice
  if err := db.Find(&Stocks).Error; err != nil {
    logger.Error(err)
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }
  c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "stocks": Stocks})
}

func NoticeDetailHandler(c *gin.Context) {

  type param struct {
    Code string `form:"code"  json:"code"`
  }

  var p param
  var err error
  if err = c.Bind(&p); err != nil {
    logger.Error(err)
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  db := c.MustGet(constant.ContextDb).(*gorm.DB)

  var notice model.Notice
  if err := db.Where("code = ?", p.Code).First(&notice).Error; err != nil {
    logger.Error(err)
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "notice": notice})
}

func NoticeUpdateHandler(c *gin.Context) {

    type param struct {
        Code      string `form:"code"  json:"code"`
        Name      string `form:"name"  json:"name"`
        BuyPrice  float64 `form:"buy_price"  json:"buy_price"`
        SellPrice float64 `form:"sell_price"  json:"sell_price"`
        BuyCount  int `form:"buy_count"  json:"buy_count"`
        SellCount int `form:"sell_count"  json:"sell_count"`
        State     string `form:"state"  json:"state"`
        Frequency string `form:"frequency"  json:"frequency"`
    }

    var p param
    var err error
    if err = c.Bind(&p); err != nil {
        logger.Error(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    db := c.MustGet(constant.ContextDb).(*gorm.DB)

    if err := db.Model(model.Notice{}).Where("code = ?", p.Code).Updates(model.Notice{
          Code: p.Code, Name:p.Name, BuyPrice:p.BuyPrice, SellPrice:p.SellPrice,
          BuyCount:p.BuyCount, SellCount:p.SellCount, State:p.State, Frequency:p.Frequency}).Error; err != nil {

        logger.Error(err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    c.JSON(http.StatusOK, gin.H{"err_code": constant.Success})
}
