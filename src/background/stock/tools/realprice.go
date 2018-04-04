package main

import (
	"background/stock/service"
	"background/stock/model"
	"background/common/logger"
	"background/stock/tools/util"

	"fmt"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main(){
	var err error

	stockName := flag.String("n", "", "stock name")
	flag.Parse()

	logger.SetLevel(0)
	db, err := gorm.Open("mysql", "root:hahawap@tcp(localhost:3306)/lyric?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(false)
	model.InitModel(db)

	var stock model.StockList
	if err = db.Where("name = ?",stockName).First(&stock).Error ; err != nil{
		logger.Error(err)
		return
	}
	jysCode := util.GetJysCodeByStockCode(stock.Code)
	var realTimeStock *model.RealTimeStock
	if err, realTimeStock = service.GetRealTimeStockInfoByStockCode(jysCode,stock.Code) ; err != nil{
		logger.Error(err)
		return
	}

	logger.Printf("Now price : " + fmt.Sprintf(realTimeStock.NowPrice))
}
