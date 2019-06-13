package common

import (
  "fmt"
  "time"
  "strings"
  "errors"

  "background/common/logger"
  "background/stock/model"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jinzhu/gorm"
  "github.com/PuerkitoBio/goquery"
)

func GetWangyiStockBasic(jys ,code string,db *gorm.DB) error{
  p := time.Now()
  var stock model.StockBasic
  stock.Code = code
  update := false
  if err := db.Where("code = ?",stock.Code).First(&stock).Error ; err == nil{
    update = true
  }

  if update && p.Weekday().String() != "Friday"{
    return "非周五不更新"
  }

  url := fmt.Sprintf("http://quotes.money.163.com/%s%s.html",jys,code)
  logger.Debug(url)
  document, err := goquery.NewDocument(url)
  if err != nil {
    logger.Error(err)

    return err
  }

  ps := document.Find(".corp_info").Eq(0).Find("p")

  title := document.Find("title").Eq(0).Text()
  title = title[:strings.Index(title,"（")]

  stock.Jys = jys
  stock.Name = title
  ps.Each(func(i int, p *goquery.Selection) {
    data := p.Text()
    data = strings.Replace(data,"\n","",-1)
    data = data[strings.Index(data,"："):]
    data = strings.Replace(data, "：","",-1)
    data = strings.Replace(data, " ","",-1)
    data = strings.Replace(data, "	","",-1)

    logger.Debug(i,":",data)
    if i == 0{
      stock.Company = data
    }else if i == 1{
      stock.Zyyw = data
    }else if i == 2{
      stock.Dsz = data
    }else if i == 3{
      stock.Zjl = data
    }else if i == 4{
      stock.Dmdh = data
    }else if i == 5{
      stock.Address = data
    }else if i == 6{
      stock.WebAddress = data
    }else if i == 7{
      stock.Zczb = data
    }else if i == 8{
      stock.ToMarketDate = data
    }else if i == 9{
      stock.Zgb = data
    }else if i == 10{
      stock.Ltgb = data
    }
  })

  if len(stock.Company) == 0 && len(stock.ToMarketDate) == 0{
    return errors.New("数据未抓全")
  }

  if update{
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
  return nil
}
