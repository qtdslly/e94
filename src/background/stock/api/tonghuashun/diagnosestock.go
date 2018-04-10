package tonghuashun

import (
	"background/common/logger"
	"background/stock/model"

	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"

	"github.com/imroc/req"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/axgle/mahonia"
)

type TrackingGuidance struct {
	Date	    string             `gorm:"date" json:"date"`                   //日期
	Suggestion  int32              `gorm:"suggestion" json:"suggestion"`       //建议类型
	Price	    string             `gorm:"price" json:"price"`                 //价格
}

type CompreData struct {
	TotalScore	   float32	      `gorm:"total_score" json:"totalScore"`                //综合得分
	TotalAnalyseInfo   string	      `gorm:"total_analyse_info" json:"totalAnalyseInfo"`   //买卖信号
	ClassNumber        float32	      `gorm:"class_number" json:"classnumber"`              //
	Suggestion	   string             `gorm:"suggestion" json:"suggestion"`                 //建议
	TotalAnalyse	   string             `gorm:"total_analyse" json:"totalAnalyse"`            //综合分析
	StockName	   string             `gorm:"stock_name" json:"stockname"`                  //股票名称
	CompanyChartInfo   []TrackingGuidance `gorm:"company_chart_info" json:"company_chart_info"`   //跟踪指导信息
}
type Comprehensive struct {
	ErrorCode     string     `gorm:"error_code" json:"errorcode"`        //错误码
	Data          CompreData
	Message       string     `gorm:"message" json:"message"`              //消息

}

/*
获取综合信息
*/
func GetComprehensive(code string,db *gorm.DB){
	cur := time.Now()
	timestamp := cur.UnixNano() / 1000000
	url := "https://vaserviece.10jqka.com.cn/diagnosestock/index.php?op=getComboData&dataType=index&stockcode=" + code + "&_=" + fmt.Sprint(timestamp)
	logger.Print(url)
	resp, err := req.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}

	recv,err := ioutil.ReadAll(resp.Response().Body)
	if err != nil{
		logger.Error(err)
	}

	logger.Print(string(recv))
	var compre Comprehensive
	if err = json.Unmarshal(recv,&compre) ; err != nil{
		logger.Error(err)
		return
	}

	var suggestion model.TonghuashunSuggestion
	suggestion.Code = code
	p := time.Now()
	suggestion.Date = fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
	if err := db.Where("code = ? and date = ?",suggestion.Code,suggestion.Date).First(&suggestion).Error ; err != nil{
		if err != gorm.ErrRecordNotFound{
			logger.Error(err)
			return
		}
	}
	enc := mahonia.NewEncoder("utf-8")

	suggestion.TotalAnalyse = enc.ConvertString(compre.Data.TotalAnalyse)
	suggestion.Suggestion = enc.ConvertString(compre.Data.Suggestion)
	suggestion.ClassNumber = compre.Data.ClassNumber
	suggestion.TotalScore = compre.Data.TotalScore
	suggestion.TotalAnalyseInfo = enc.ConvertString(compre.Data.TotalAnalyseInfo)
	if err := db.Save(&suggestion).Error ; err != nil{
		logger.Error(err)
		return
	}
}