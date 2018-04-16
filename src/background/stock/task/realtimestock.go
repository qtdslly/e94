package task

import (
	"background/common/logger"
	"background/stock/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"background/stock/service"
	"time"
)


func SyncAllRealTimeStockInfo(db *gorm.DB){
	var p = time.Now()
	var task model.StockTask
	if err := db.Where("`key` = 'realtimestock'").First(&task).Error ; err != nil{
		logger.Error(err)
		return
	}

	today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
	if task.Date > today{
		logger.Debug("今日数据已抓取")
		return
	}

	//六位代码00、200、300开头的都是深圳的股票，00开头的是深圳A股，200开头的是深圳B股，
	//300开头的是创业板(创业板都是在深市交易的)。六位代码60、900开头的是上海的股票，
	//60开头的都是上海A股，900开头的是上海B股。六位代码184开头的都是深圳封闭式基金，
	//六位代码500开头的都是上海封闭式基金。
	//
	//证券交易所挂牌的证券投资基金产品编码为6位数字，
	//前两位为15、16、18的是深圳证券交易所基金，
	//50、51、52的是上海证券交易所基金。
	//不在证券交易所挂牌的证券投资基金编码规则为：基金编码为6位数字，
	//前两位为基金管理公司的注册登记机构编码(TA编码)，后四位为产品流水号
	headList := []string{"600","601","603","9009","5000","5010","5013","5020","5058","5100","5101","5102","5103","5104",
	"5105","5106","5107","5108","5109","5110","5112","5116","5117","5118","5119","5120","5121","5122",
	"5123","5125","5126","5128","5129","5130","5131","5135","5136","5188","5190","5191","5192","5193",
	"5195","5196","5197","5198","5199","5211","5212","5213","5216","5219","5220","5221","5222","5223",
	"5225","5226","5227","5229","5230","5231","5232","5233","5235","5236","5237","5239","000",
	"0016","0018","0019","002","3000","3001","3002","3003","3004","3005","3006","200","1847",
	"1848","1500","1501","1502","1503","1590","1599","1601","1602","1603","1604","1605","1606",
	"1607","1608","1609","1610","1611","1612","1615","1616","1617","1618","1619","1620","1621",
	"1622","1623","1624","1625","1626","1627","1630","1631","1632","1633","1634","1635","1638",
	"1639","1641","1642","1643","1644","1645","1646","1647","1648","1649","1653","1655","1657",
	"1658","1660","1661","1664","1669","1670","1675","1681","1683","1691","1692"}
	jysCode := "sh"
        for _ ,head := range headList{
		if head == "000" {
			jysCode = "sz"
		}
		endi := 0
		if len(head) == 4{
			endi = 99
		}else if len(head) == 3{
			endi = 999
		}else if len(head) == 2{
			endi = 9999
		}
		i := 0
		for i <= endi{
			stockCode := ""
			if endi < 100{
				stockCode = head + fmt.Sprintf("%02d",i)
			}else if endi < 1000{
				stockCode = head + fmt.Sprintf("%03d",i)
			}else if endi < 10000{
				stockCode = head + fmt.Sprintf("%04d",i)
			}
			var realTimeStock *model.RealTimeStock
			var err error
			if err,realTimeStock = service.GetRealTimeStockInfoByStockCode(jysCode, stockCode) ; err != nil{
				if err.Error() == "股票代码不存在"{
					i++
					continue
				}
			}

			var count int
			if err = db.Where("stock_code = ? and deal_date = ?",realTimeStock.StockCode,realTimeStock.DealDate).Table("real_time_stock").Count(&count).Error ; err != nil{
				logger.Error(err)
				return
			}
			if count == 0{
				if err := db.Create(&realTimeStock).Error ; err != nil{
					logger.Error("保存股票实时信息失败,",err)
					return
				}
			}

			var stockHistoryDataQ model.StockHistoryDataQ
			if err = db.Where("code = ? and date = ?",realTimeStock.StockCode,realTimeStock.DealDate).First(&stockHistoryDataQ).Error ; err != nil{
				if err != gorm.ErrRecordNotFound{
					logger.Error(err)
					return
				}else{
					stockHistoryDataQ.Code = realTimeStock.StockCode
					stockHistoryDataQ.Date = realTimeStock.DealDate
					stockHistoryDataQ.Close = realTimeStock.NowPrice
					stockHistoryDataQ.High = realTimeStock.TodayHighPrice
					stockHistoryDataQ.Low = realTimeStock.TodayLowPrice
					stockHistoryDataQ.Open = realTimeStock.TodayOpenPrice
					stockHistoryDataQ.Volume = float64(realTimeStock.DealCount) / 100.00

					if err = db.Create(&stockHistoryDataQ).Error ; err != nil{
						logger.Error(err)
						return
					}
				}
			}

			var stockHistoryDataQNew model.StockHistoryDataQNew
			if err = db.Where("code = ? and date = ?",realTimeStock.StockCode,realTimeStock.DealDate).First(&stockHistoryDataQNew).Error ; err != nil{
				if err != gorm.ErrRecordNotFound{
					logger.Error(err)
					return
				}else{
					stockHistoryDataQNew.Code = realTimeStock.StockCode
					stockHistoryDataQNew.Date = realTimeStock.DealDate
					stockHistoryDataQNew.Close = realTimeStock.NowPrice
					stockHistoryDataQNew.High = realTimeStock.TodayHighPrice
					stockHistoryDataQNew.Low = realTimeStock.TodayLowPrice
					stockHistoryDataQNew.Open = realTimeStock.TodayOpenPrice
					stockHistoryDataQNew.Volume = float64(realTimeStock.DealCount) / 100.00
					stockHistoryDataQNew.Amount = realTimeStock.DealMoney

					if err = db.Create(&stockHistoryDataQNew).Error ; err != nil{
						logger.Error(err)
						return
					}
				}
			}


			i++
		}

	}
	now := time.Now()
	dd, _ := time.ParseDuration("24h")
	to := now.Add(dd)
	tomorry := fmt.Sprintf("%04d-%02d-%02d",to.Year(),to.Month(),to.Day())
	if err := db.Model(model.StockTask{}).Where("`key` = 'realtimestock'").Update("date", tomorry).Error; err != nil {
		logger.Error(err)
		return
	}
}
