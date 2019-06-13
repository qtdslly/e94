package tonghuashun

import (

  "fmt"
  "time"
  "strings"
  "errors"
  "io/ioutil"
  "encoding/json"

  "background/common/logger"
  "background/stock/model"

  "github.com/imroc/req"
  "github.com/axgle/mahonia"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
)

type TrackingGuidance struct {
  Date       string             `gorm:"date" json:"date"`             //日期
  Suggestion int32              `gorm:"suggestion" json:"suggestion"` //建议类型
  Price      string             `gorm:"price" json:"price"`           //价格
}

type CompreData struct {
  TotalScore       float32        `gorm:"total_score" json:"totalScore"`                    //综合得分
  TotalAnalyseInfo string        `gorm:"total_analyse_info" json:"totalAnalyseInfo"`        //买卖信号
  ClassNumber      float32        `gorm:"class_number" json:"classnumber"`                  //
  Suggestion       string             `gorm:"suggestion" json:"suggestion"`                 //建议
  TotalAnalyse     string             `gorm:"total_analyse" json:"totalAnalyse"`            //综合分析
  StockName        string             `gorm:"stock_name" json:"stockname"`                  //股票名称
  CompanyChartInfo []TrackingGuidance `gorm:"company_chart_info" json:"company_chart_info"` //跟踪指导信息
}
type Comprehensive struct {
  ErrorCode string     `gorm:"error_code" json:"errorcode"` //错误码
  Data      CompreData
  Message   string     `gorm:"message" json:"message"`      //消息
}

/*
获取综合信息
*/
func GetComprehensive(code string, db *gorm.DB) error {
  var err error

  cur := time.Now()

  var suggestion model.TonghuashunSuggestion
  suggestion.Code = code
  suggestion.Date = fmt.Sprintf("%04d-%02d-%02d", cur.Year(), cur.Month(), cur.Day())
  if err = db.Where("code = ? and date = ?", suggestion.Code, suggestion.Date).First(&suggestion).Error; err == nil {
    return errors.New(suggestion.Code + suggestion.Date + "日数据以抓取")
  }

  timestamp := cur.UnixNano() / 1000000
  url := "https://vaserviece.10jqka.com.cn/diagnosestock/index.php?op=getComboData&dataType=index&stockcode=" + code + "&_=" + fmt.Sprint(timestamp)
  logger.Print(url)
  resp, err := req.Get(url)
  if err != nil {
    logger.Error(err)
    return err
  }

  recv, err := ioutil.ReadAll(resp.Response().Body)
  if err != nil {
    logger.Error(err)
  }

  logger.Print(string(recv))
  var compre Comprehensive
  if err = json.Unmarshal(recv, &compre); err != nil {
    logger.Error(err)
    return err
  }

  enc := mahonia.NewEncoder("utf-8")
  suggestion.TotalAnalyse = enc.ConvertString(compre.Data.TotalAnalyse)
  suggestion.Suggestion = enc.ConvertString(compre.Data.Suggestion)
  suggestion.ClassNumber = compre.Data.ClassNumber
  suggestion.TotalScore = compre.Data.TotalScore

  suggestion.TotalAnalyseInfo = enc.ConvertString(compre.Data.TotalAnalyseInfo)
  if err := db.Create(&suggestion).Error; err != nil {
    logger.Error(err)
    return err
  }
  return nil
}

type Capital struct {
  date   string        `gorm:"date" json:"date"`      //日期
  state  int32        `gorm:"state" json:"state"`     //状态
  amount float32        `gorm:"amount" json:"amount"` //金额
}

type ControlData struct {
  FundAnalyse  string        `gorm:"fund_analyse" json:"fundAnalyse"`         //主力迹象
  CurrentFund  string        `gorm:"current_fund" json:"currentFund"`         //当前资金净流入
  State        string        `gorm:"state" json:"state"`                      //资金流入流出状态
  Amount       float32        `gorm:"amount" json:"amount"`                   //汇总金额
  FundDataJson []Capital        `gorm:"fund_data_json" json:"fund_data_json"` //资金数据
  ControlValue string        `gorm:"control_value" json:"controlvalue"`       //控盘度
}

type ControlInfo struct {
  ErrorCode string      `gorm:"error_code" json:"errorcode"` //错误码
  Data      ControlData
  Message   string      `gorm:"message" json:"message"`      //消息
}
/*获取资金及控盘信息*/
func GetControlInfo(code string, db *gorm.DB) (error) {
  var err error
  cur := time.Now()

  var mainForceControl model.TonghuashunMainForceControl
  mainForceControl.Code = code
  mainForceControl.Date = fmt.Sprintf("%04d-%02d-%02d", cur.Year(), cur.Month(), cur.Day())
  if err = db.Where("code = ? and date = ?", mainForceControl.Code, mainForceControl.Date).First(&mainForceControl).Error; err == nil {
    return errors.New(mainForceControl.Code + mainForceControl.Date + "日数据已抓取")
  }

  if err == nil {
    return nil
  }

  timestamp := cur.UnixNano() / 1000000
  url := "https://vaserviece.10jqka.com.cn/diagnosestock/index.php?op=getComboData&&dataType=currentFunds&&stockcode=" + code + "&_=" + fmt.Sprint(timestamp)
  logger.Print(url)
  resp, err := req.Get(url)
  if err != nil {
    logger.Error(err)
    return err
  }

  recv, err := ioutil.ReadAll(resp.Response().Body)
  if err != nil {
    logger.Error(err)
    return err
  }

  logger.Print(string(recv))
  recvTmp := string(recv)
  recvTmp = strings.Replace(recvTmp, "\"<span style=\\\"color:green\\\">", "", -1)
  recvTmp = strings.Replace(recvTmp, "<\\/span>\"", "", -1)
  errCode := strings.Split(recvTmp, "\"")[3]
  if errCode != "0" {
    logger.Debug("通花顺返回错误!!!")
    return nil
  }
  recv = []byte(recvTmp)
  var control ControlInfo
  if err = json.Unmarshal(recv, &control); err != nil {
    logger.Error(err)
    return err
  }

  enc := mahonia.NewEncoder("utf-8")

  mainForceControl.FundAnalyse = enc.ConvertString(control.Data.FundAnalyse)
  mainForceControl.CurrentFund = enc.ConvertString(control.Data.CurrentFund)
  mainForceControl.ControlValue = enc.ConvertString(control.Data.ControlValue)
  mainForceControl.State = enc.ConvertString(control.Data.State)
  mainForceControl.Amount = control.Data.Amount
  if err = db.Create(&mainForceControl).Error; err != nil {
    logger.Error(err)
    return err
  }
  return nil
}
