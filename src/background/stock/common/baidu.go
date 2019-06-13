package common

import (
  "io/ioutil"
  "errors"
  "time"
  "fmt"
  "strings"
  "net/http"

  "background/common/logger"
  "background/stock/model"
  "background/common/util"

  "github.com/axgle/mahonia"
  "github.com/tidwall/gjson"
  "github.com/jinzhu/gorm"
)

//百度智能测评 https://gupiao.baidu.com/
func Zncp(db *gorm.DB,jys string,code string) error {
  p := time.Now()

  today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
  var count int
  if err := db.Table("baidu_zncp").Where("date = ? and code = ?",today,code).Count(&count).Error ; err != nil{
    logger.Error(err)
    return err
  }

  if count != 0{
    return errors.New("数据已抓取")
  }
  timeStamp := p.Unix()
  url := fmt.Sprintf("https://gupiao.baidu.com/api/rails/intelligentevaluation?from=pc&os_ver=1&cuid=xxx&vv=100&format=json&code=%s%s&timestamp=%d",jys,code,timeStamp)

  logger.Debug(url)
  requ, err := http.NewRequest("GET", url, nil)
  client := &http.Client{}
  resp, err := client.Do(requ)
  if err != nil {
    logger.Error(err)
    return err
  }

  defer resp.Body.Close()

  recv, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    logger.Error(err)
    return err
  }

  logger.Debug(string(recv))

  response := string(recv)
  logger.Debug(response)
  status := gjson.Get(response, "errorNo").Int()

  if status != 0 {
    logger.Error("数据抓取失败")
    return errors.New("数据抓取失败")
  }

  response,_ = util.DecodeToGBK(response)
  response = strings.Replace(response,"\"{","{",-1)
  response = strings.Replace(response,"\"}","}",-1)
  response = strings.Replace(response,"\\\"","\"",-1)
  response = strings.Replace(response,"\\\\","\\",-1)

  logger.Debug(response)

  zlw := gjson.Get(response, "data.ub").String()
  zcw := gjson.Get(response, "data.lb").String()
  szgl := gjson.Get(response, "data.szgl").String()
  kpld := gjson.Get(response, "data.kpld").String()

  logger.Debug(zlw)
  logger.Debug(zcw)
  logger.Debug(szgl)
  logger.Debug(kpld)

  var name,jys1 string
  var price ,high,low, netChange,netChangeRatio,amplitudeRatio,turnoverRatio,close,open float64
  var volume,capitalization uint64
  stockbasic := gjson.Get(response, "data.stockbasic")

  if stockbasic.Exists() {
    items := stockbasic.Array()
    for _, item := range items {
      name = item.Get("stockName").String()
      jys1 = item.Get("exchange").String()
      price = item.Get("close").Float()
      low = item.Get("low").Float()
      high = item.Get("high").Float()
      netChange = item.Get("netChange").Float()
      netChangeRatio = item.Get("netChangeRatio").Float()
      volume = item.Get("volume").Uint()
      capitalization = item.Get("capitalization").Uint()

      amplitudeRatio = item.Get("amplitudeRatio").Float()
      turnoverRatio = item.Get("turnoverRatio").Float()
      close = item.Get("preClose").Float()
      open = item.Get("open").Float()
      break
    }
  }

  var stock model.BaiduZncp
  stock.Date = today
  stock.Code = code
  stock.Name = name

  enc := mahonia.NewEncoder("utf-8")
  stock.Name = enc.ConvertString(name)

  fmt.Println(stock.Name)

  logger.Debug(stock.Name)

  if jys1 == "sh"{
    stock.Jys = "0"
  }else {
    stock.Jys = "1"
  }
  stock.Price = price
  stock.Zcw = zcw
  stock.Zlw = zlw
  stock.Kpld = kpld
  stock.Szgl = szgl
  stock.Low = low
  stock.High = high
  stock.Close = close
  stock.Open = open
  stock.Capitalization = capitalization
  stock.NetChange = netChange
  stock.NetChangeRatio = netChangeRatio
  stock.Volume = volume
  stock.AmplitudeRatio = amplitudeRatio
  stock.TurnoverRatio = turnoverRatio

  if err := db.Create(&stock).Error ; err != nil{
    logger.Error(err)
    return err
  }

  return nil

}
