package task

import (
  "fmt"
  "time"
  "io/ioutil"

  "background/common/util"
  "background/common/logger"
  cc "background/stock/common"

  "github.com/imroc/req"
  "github.com/tidwall/gjson"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
  "background/stock/model"
)

/*
获取所有股票市场行情
*/
func GetBasicInfo(db *gorm.DB){
  var err error
  p := time.Now()
  today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
  var stocks []model.WangYiStock
  if err = db.Order("code asc").Where("date = ?",today).Find(&stocks).Error; err != nil {
    logger.Error(err)
    return
  }

  for _,stock := range stocks{
    cc.GetWangyiStockBasic(stock.Jys,stock.Code,db)
  }
}

func GetWangyiStockList(db *gorm.DB)(error){

  logger.Debug("开始抓取网易股票数据")
  p := time.Now()
  date := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())

  dd, _ := time.ParseDuration("24h")
  to := p.Add(dd)
  tomorry := fmt.Sprintf("%04d-%02d-%02d",to.Year(),to.Month(),to.Day())

  var task model.StockTask
  if err := db.Where("`key` = 'wangyi_stock'").First(&task).Error ; err != nil{
    logger.Error(err)
    return err
  }

  today := fmt.Sprintf("%04d-%02d-%02d",p.Year(),p.Month(),p.Day())
  if task.Date > today{
    logger.Debug("今日数据已抓取")
    return nil
  }

  var page int64 = 0
  for{
    url := "http://quotes.money.163.com/hs/service/diyrank.php?host=http://quotes.money.163.com/hs/service/diyrank.php&page=" + fmt.Sprintf("%d",page) + "&query=STYPE:EQA&fields=NO,SYMBOL,NAME,PRICE,PERCENT,UPDOWN,FIVE_MINUTE,OPEN,YESTCLOSE,HIGH,LOW,VOLUME,TURNOVER,HS,LB,WB,ZF,PE,MCAP,TCAP,MFSUM,MFRATIO.MFRATIO2,MFRATIO.MFRATIO10,SNAME,CODE,ANNOUNMT,UVSNEWS&sort=PERCENT&order=desc&count=24&type=query"
    resp, err := req.Get(url)
    if err != nil {
      logger.Error(err)
      return err
    }

    recv,err := ioutil.ReadAll(resp.Response().Body)
    if err != nil{
      logger.Error(err)
      return err
    }

    data,_ := util.DecodeToGBK(string(recv))

    logger.Debug(data)

    pageCount := gjson.Get(data, "pagecount").Int()

    list := gjson.Get(data, "list")

    if list.Exists() {
      items := list.Array()

      for _, item := range items {
        var stock model.WangYiStock
        stock.Code = item.Get("SYMBOL").String()
        stock.Date = date

        if err := db.Where("date = ? and code = ?",stock.Date,stock.Code).First(&stock).Error ; err == nil{
          continue
        }
        stock.Jys = item.Get("CODE").String()[0:1]
        stock.FiveMinute = item.Get("FIVE_MINUTE").String()
        stock.High = item.Get("HIGH").String()
        stock.Hs = item.Get("HS").String()
        stock.Lb = item.Get("LB").String()
        stock.Low = item.Get("LOW").String()
        stock.Mcap = item.Get("MCAP").String()
        stock.Eps = item.Get("MFSUM").String()
        stock.Name = item.Get("NAME").String()
        stock.Open = item.Get("OPEN").String()
        stock.Percent = item.Get("PERCENT").String()
        stock.Price = item.Get("PRICE").String()
        stock.Pe = item.Get("PE").String()
        stock.Code = item.Get("SYMBOL").String()
        stock.Tcap = item.Get("TCAP").String()
        stock.Turnover = item.Get("TURNOVER").String()
        stock.UpDown = item.Get("UPDOWN").String()
        stock.Volume = item.Get("VOLUME").String()
        stock.Wb = item.Get("WB").String()
        stock.YestClose = item.Get("YESTCLOSE").String()
        stock.Zf = item.Get("ZF").String()
        stock.NetProfit = item.Get("MFRATIO.MFRATIO2").String()
        stock.TotalRevenue = item.Get("MFRATIO.MFRATIO10").String()

        if err := db.Create(&stock).Error ; err != nil{
          logger.Error(err)
          return err
        }
      }
    }

    page++
    if page >= pageCount{
      break
    }

  }

  if err := db.Model(model.StockTask{}).Where("`key` = 'wangyi_stock'").Update("date", tomorry).Error; err != nil {
    logger.Error(err)
    return err
  }

  return nil
}
