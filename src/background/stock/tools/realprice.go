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

	if *stockName == "all"{
		var transPrompts []model.TransPrompt
		if err = db.Find(&transPrompts).Error; err != nil{
			logger.Error(err)
			return
		}
		for _ , transPrompt := range transPrompts{
			GetNowStockInfo(transPrompt.StockCode,db)
		}
		return
	}

	var stock model.StockList
	if err = db.Where("name = ?",*stockName).First(&stock).Error ; err != nil{
		logger.Error("stock name:",*stockName , " error:",err)
		return
	}
	GetNowStockInfo(stock.Code,db)
}

func GetNowStockInfo(code string,db *gorm.DB){
	var err error

	jysCode := util.GetJysCodeByStockCode(code)
	var realTimeStock *model.RealTimeStock
	if err, realTimeStock = service.GetRealTimeStockInfoByStockCode(jysCode,code) ; err != nil{
		logger.Error(err)
		return
	}

	name := util.GetNameByCode(code,db)

	logger.Printf(name + " Now price : " + fmt.Sprint(realTimeStock.NowPrice) + " Rose : " + fmt.Sprintf("%.2f",(realTimeStock.NowPrice - realTimeStock.YestdayClosePrice) / realTimeStock.YestdayClosePrice * 100.00))

}
