package main

import (
  "background/ezhantao/service"
  "background/wechart/config"
  "background/common/logger"
  "background/common/util"

  "github.com/tidwall/gjson"

  "fmt"
  "strings"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())

  url := "https://detail.tmall.com/item.htm?id=37016930337"
  itemId := GetItemId(url)
  logger.Debug(itemId)
  recv := service.GetXuanPinKuGoodsInfo1(itemId)
  logger.Debug(recv)

  list := gjson.Get(recv, "tbk_item_info_get_response.results.n_tbk_item")

  if list.Exists() {
    items := list.Array()
    for _, item := range items {
      title := item.Get("title").String()    //一口价

      reserve_price := item.Get("reserve_price").String()    //一口价
      zk_final_price := item.Get("zk_final_price").String()  //折后价
      coupon_click_url := item.Get("coupon_click_url").String()  //折后价

      logger.Debug(title)
      logger.Debug(reserve_price)
      logger.Debug(zk_final_price)

      response := service.GetTuiGuangQuan("",itemId,"")
      logger.Debug(response)
      //优惠券门槛金额
      coupon_start_fee := gjson.Get(response, "tbk_coupon_get_response.data.coupon_start_fee").String()
      //优惠券金额
      coupon_amount := gjson.Get(response, "tbk_coupon_get_response.data.coupon_amount").String()
      //券类型，1 表示全网公开券，4 表示妈妈渠道券
      coupon_src_scene := gjson.Get(response, "tbk_coupon_get_response.data.coupon_src_scene").Int()
      var scene string
      if coupon_src_scene == 1{
        scene = "全网公开券"
      }else if coupon_src_scene == 4{
        scene = "妈妈渠道券"
      }else{
        scene = fmt.Sprintf("%d",coupon_src_scene)
      }
      //券属性，0表示店铺券，1表示单品券
      coupon_type := gjson.Get(response, "tbk_coupon_get_response.data.coupon_type").Int()
      var couponType string
      if coupon_type == 0{
        couponType = "店铺券"
      }else if coupon_type == 1{
        couponType = "单品券"
      }else{
        couponType = fmt.Sprintf("%d",coupon_type)
      }
      //logger.Debug(response)

      fp := util.GetAmt(zk_final_price)
      ca := util.GetAmt(coupon_amount)
      if fp - ca < 5{
        logger.Debug(fp-ca,"|",zk_final_price,"|",coupon_start_fee,"|",coupon_amount,"|",scene,"|",couponType,"|","|",title)
      }

      logger.Debug("宝贝名称:",title)
      logger.Debug("宝贝地址:",url)
      logger.Debug("领券地址:",coupon_click_url)
      logger.Debug("一口价:",zk_final_price)
      logger.Debug("折后价:",coupon_start_fee)
      logger.Debug("优惠券:",coupon_amount)
      logger.Debug("到手价:",fp - ca)

      break
    }
  }

}


func GetItemId(url string)string{
  id := strings.LastIndex(url,"=")
  return url[id+1:]
}
