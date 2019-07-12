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

//获取底部反转的股票
func GetDbfz(db *gorm.DB) error {
  return WenCai("底部反转",db)
}

//获取cci买入信号的股票
func GetCCImrxh(db *gorm.DB) error {
  return WenCai("cci买入信号",db)
}

//获取同花顺上升通道股票
func GetSstd(db *gorm.DB) error {
  return WenCai("上升通道",db)
}
func WenCai(title string,db *gorm.DB)error{
  url := "http://www.iwencai.com/stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=" + title + "&queryarea="

  logger.Debug(url)

  logger.Debug(url)
  requ, err := http.NewRequest("GET", url, nil)
  requ.Header.Add("Host","www.iwencai.com")
  requ.Header.Add("Cookie","PHPSESSID=4640e971fb54745bbf76c30974365961; cid=4640e971fb54745bbf76c309743659611560849630; ComputerID=4640e971fb54745bbf76c309743659611560849630; guideState=1; v=AlL4Iq4WW-UM4KcZ-CHSZPdwoxM3Y1b9iGdKIRyrfoXwL_yFBPOmDVj3mjTv")

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


  //https://eq.10jqka.com.cn/newOpinion/app/main.php?question=%E8%BF%9E%E7%BB%AD%E4%B8%8A%E6%B6%A83%E5%A4%A9&source=ths_mobile_yuyinzhushou&multistage=%7B%22is_multistage%22%3Afalse%7D&channel=2770&entity_info=%7B%22device_type%22%3A%22iphone%22%2C%22stock_code%22%3A%22%22%2C%22stock_name%22%3A%22%22%2C%22comefrom%22%3A%2260%22%2C%22entry_type%22%3A%220%22%7D&kefu_user_id=&kefu_record_id=&user_id=226874239&log_info=%7B%22other_info%22%3A%22%7B%5C%22eventId%5C%22%3A%5C%22iwencai_app_send_click%5C%22%2C%5C%22ct%5C%22%3A1560910617472%7D%22%2C%22input_type%22%3A%22typewrite%22%2C%22edition%22%3A%22normal%22%2C%22entry_type%22%3A%220%22%2C%22udid%22%3A%22854A0ECB-D1C1-48A8-A37B-6A93FA83D7E9%22%7D&user_name=mx_226874239&pagename=cn.com.10jqka.IHexin&innerversion=I037.08.357&statid=yuyinzhushou&version=2.0&add_info=%7B%22zhengu%22%3A%7B%7D%2C%22command%22%3A%7B%22needCommand%22%3Atrue%7D%2C%22talk%22%3A%7B%22needIdentity%22%3Atrue%2C%22needImg%22%3Atrue%7D%2C%22finance%22%3A%7B%22ans_type%22%3A%22%22%2C%22needFurther%22%3Atrue%2C%22needFurtherKnow%22%3Atrue%7D%7D&control=ControlCenterV2&action=getAnswer
  //https://search.10jqka.com.cn/unified-wap/cache?page=2&perpage=30&token=5df4383dbef1508d5f377fe08566fd4a&_=1560910254768&callback=Zepto1560910249044
  //https://search.10jqka.com.cn/unified-wap/get-parser-data?w=连续上涨天数=3日 &domain=abs_股票领域
  //https://search.10jqka.com.cn/unified-wap/cache?page=2&perpage=30&token=800d127f828d224f4b71a0b64fd2e112&_=1560911329824&callback=Zepto1560911199860
  //https://search.10jqka.com.cn/unified-wap/cache?page=2&perpage=30&token=6b6c76e59223786a630a6f059cd46be4&_=1560911928161&callback=Zepto1560911837722
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
      var stock model.TonghuashunJsxt
      stock.Date = today
      stock.Title = title
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
      if err := db.Table("tonghuashun_jsxt").Where("date = ? and code = ? and title = ?",stock.Date,stock.Code,title).Count(&count).Error ; err != nil{
        logger.Error(err)
        return err
      }

      logger.Debug(count)
      if count != 0{
        if err := db.Exec("delete from stock.tonghuashun_jsxt where date = ? and code = ? and title = ?",today,stock.Code,title).Error ; err != nil{
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


//撑压分析
func Cyfx(code string) string{
  url := "https://eq.10jqka.com.cn/wencai/interface.php?marketid=17&op=getConfigByRealTime&requesttype=2&stockcode=" + code + "&tabIndex=2&userid=226874236"
  logger.Debug(url)
  requ, err := http.NewRequest("GET", url, nil)

  client := &http.Client{}
  resp, err := client.Do(requ)
  if err != nil {
    logger.Error(err)
    return ""
  }

  defer resp.Body.Close()

  recv, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    logger.Error(err)
    return ""
  }


  response := string(recv)
  logger.Debug(response)
  errCode := gjson.Get(response, "code").Int()

  if errCode != 0 {
    logger.Error("数据抓取失败")
    return ""
  }
  content := gjson.Get(response, "data.content")
  if !content.Exists(){
    return ""
  }
  return content.Array()[0].String()
}



//上升通道
//https://search.10jqka.com.cn/unified-wap/get-parser-data?w=%E4%B8%8A%E5%8D%87%E9%80%9A%E9%81%93%20&domain=abs_%E8%82%A1%E7%A5%A8%E9%A2%86%E5%9F%
//condition [{"indexName":"上升通道","indexProperties":["nodate 1","交易日期 20190617"],"indexPropertiesMap":{"交易日期":"20190617","nodate":"1"},"type":"tech","sonSize":0,"reportType":"TRADE_DAILY","valueType":"","domain":"abs_股票领域","source":"new_parser","dateType":"交易日期","tag":"上升通道","uiText":"上升通道","chunkedResult":"上升通道","queryText":"上升通道","relatedSize":0}]
