package common

import (
	"github.com/jinzhu/gorm"
	"background/stock/model"
	"background/common/logger"
)
func GetJysCodeByStockCode(code string)(string){
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
	var head string
	head = code[0:3]
	if head == "200" || head == "300" || head == "184"{
		return "sz"
	}else if (head == "900" || head == "500"){
		return "sh"
	}
	head = code[0:2]
	if head == "00" || head == "15" || head == "16" || head == "18"{
		return "sz"
	}else if head == "60" || head == "50" || head == "51" || head == "52"{
		return "sh"
	}
	return ""
}

func GetNameByCode(code string,db *gorm.DB)(string){
	var stock model.StockList
	if err := db.Where("code = ?",code).First(&stock).Error; err != nil{
		logger.Error(err)
		return ""
	}
	return stock.Name
}
