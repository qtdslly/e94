package main

import (
  "background/ezhantao/service"
  "background/wechart/config"
  "background/common/logger"
  "background/common/util"
  "github.com/tidwall/gjson"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())

  title := "古尚古iPhone6数据线6s苹果5加长5s手机6Plus充电线器7P快充8X短se六7P冲iphonex正品平板电脑ipad适用XS Max"
  recv := service.GetHaoQuanList(1,100,title)
  logger.Debug(recv)

  list := gjson.Get(recv, "tbk_dg_item_coupon_get_response.results.tbk_coupon")

  minPrice := 1000000.00
  url := ""
  logo := ""
  if list.Exists() {
    items := list.Array()
    for _, item := range items {
      item_id := item.Get("num_iid").String()                //商品ID
      //reserve_price := item.Get("reserve_price").String()    //一口价
      zk_final_price := item.Get("zk_final_price").String()  //折后价
      title1 := item.Get("title").String()  //折后价
      if title != title1{
        continue
      }

      response := service.GetTuiGuangQuan("",item_id,"")
      coupon_click_url := item.Get("coupon_click_url").String()  //折后价
      pict_url := item.Get("pict_url").String()  //折后价

      //优惠券金额
      coupon_amount := gjson.Get(response, "tbk_coupon_get_response.data.coupon_amount").String()
      //券类型，1 表示全网公开券，4 表示妈妈渠道券


      fp := util.GetAmt(zk_final_price)
      ca := util.GetAmt(coupon_amount)

      price := fp - ca;
      if price < minPrice{
        minPrice = price
        url = coupon_click_url
        logo = pict_url
      }
    }
  }

  response := service.GetTPwd(url,title,logo)
  pwd := gjson.Get(response, "tbk_tpwd_create_response.data.model").String()

  logger.Debug(title,"|",pwd);

}

