package common

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "errors"
  "strings"
  "time"

  "background/stock/model"
  "background/common/logger"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jinzhu/gorm"
  "github.com/tidwall/gjson"
  "github.com/axgle/mahonia"
)

//获取同花顺上升通道股票
func Sstd(db *gorm.DB)error{
  url := fmt.Sprintf("http://www.iwencai.com/stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=上升通道&queryarea=")

  logger.Debug(url)

  logger.Debug(url)
  requ, err := http.NewRequest("GET", url, nil)
  requ.Header.Add("Host","www.iwencai.com")
  requ.Header.Add("Cookie","cid=ff58ed3612fdf50d2ecdbbe84e7473b21560331514; ComputerID=ff58ed3612fdf50d2ecdbbe84e7473b21560331514; guideState=1; PHPSESSID=18cf7f95c14b1723b57136b28107c9ef; v=AgPSKiyWWmjMNxacNG8zAx7LksyueJe60Qzb7jXgX2LZ9C26vUgnCuHcaz9G")

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


  response := string(recv)
  logger.Debug(response)
  success := gjson.Get(response, "success").Bool()

  if !success {
    logger.Error("数据抓取失败")
    return errors.New("数据抓取失败")
  }

  result := gjson.Get(response, "data.result.result")

  if result.Exists() {
    items := result.Array()
    p := time.Now()
    today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
    for _, item := range items {
      var stock model.TonghuashunSstd
      stock.Date = today
      ss := item.Array()
      for j ,s := range ss{
        enc := mahonia.NewEncoder("utf-8")
        s1 := enc.ConvertString(s.String())
        logger.Debug(s1)
        if j == 0{
          stock.Code = strings.Split(s1,".")[0]
          stock.Jys = strings.ToLower(strings.Split(s1,".")[1])
        }else if j == 1{
          stock.Name = s1
        }else if j == 2{
          stock.Price = s1
        }else if j == 3{
          stock.Zdf = s1
        }else if j == 4{
          stock.Mrxh = s1
        }else if j == 5{
          stock.Jsxt = s1
        }
      }
      var count int
      if err := db.Table("tonghuashun_sstd").Where("date = ? and code = ?",stock.Date,stock.Code).Count(&count).Error ; err != nil{
        logger.Error(err)
        return err
      }

      logger.Debug(count)
      if count != 0{
        if err := db.Exec("delete from stock.tonghuashun_sstd where date = ? and code = ?",today,stock.Code).Error ; err != nil{
          logger.Error(err)
          return err
        }
        if err := db.Save(&stock).Error ; err != nil{
          logger.Error(err)
          return err
        }
      }else{
        if err := db.Create(&stock).Error ; err != nil{
          logger.Error(err)
          return err
        }
      }

    }
  }

  return nil
}
