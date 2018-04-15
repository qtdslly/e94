package cms

import (
	"net/http"
	"background/stock/model"
	"background/common/constant"
	"background/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"fmt"
)

func StockListHandler(c *gin.Context) {
	type param struct {
		StockCode       uint32 `form:"stock_code" binding:"required"`
	}
	var p param
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)

	var realStockInfo model.RealTimeStock
	if err = db.Where("stock_code = ?" , p.StockCode).First(&realStockInfo).Error ; err != nil{
		logger.Error("query realtime_stock err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}


	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data":fmt.Sprint(realStockInfo.NowPrice) })
}


func StockPriceHandler(c *gin.Context) {
	type param struct {
		StockCode       uint32 `form:"stock_code" binding:"required"`
		Start           string `form:"start" binding:"required"`
		End             string `form:"end" binding:"required"`
	}
	var p param
	var err error

	if err = c.Bind(&p); err != nil {
		logger.Error(err)
		logger.Error(err)

		return
	}

	db := c.MustGet(constant.ContextDb).(*gorm.DB)
	var realStockInfos []model.RealTimeStock
	if err = db.Where("stock_code = ? and deal_date between ? an ?" , p.StockCode,p.Start,p.End).Find(&realStockInfos).Error ; err != nil{
		logger.Error("query realtime_stock err!!!,",err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"err_code": constant.Success, "data":realStockInfos })
}
