package main

import (
  "background/ezhantao/service"
  "background/wechart/config"
  "background/common/logger"
  "background/common/util"

  "github.com/tidwall/gjson"

  "fmt"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())

  title := "数据线usb小风扇迷你静音小型小电扇宿舍大风力学生寝室办公室桌上床上台式笔记本电脑充电宝家用4/6/8铁艺"
  recv := service.GetHaoQuanList(1,20,title)
  logger.Debug(recv)

  list := gjson.Get(recv, "tbk_dg_item_coupon_get_response.results.tbk_coupon")

  if list.Exists() {
    items := list.Array()
    for _, item := range items {
      item_id := item.Get("num_iid").String()                //商品ID
      reserve_price := item.Get("reserve_price").String()    //一口价
      zk_final_price := item.Get("zk_final_price").String()  //折后价
      coupon_click_url := item.Get("coupon_click_url").String()  //折后价

      logger.Debug(item_id)
      logger.Debug(reserve_price)
      logger.Debug(zk_final_price)

      response := service.GetTuiGuangQuan("",item_id,"")
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
        logger.Debug(fp-ca,"|",zk_final_price,"|",coupon_start_fee,"|",coupon_amount,"|",scene,"|",couponType,"|",title,"|",coupon_click_url)
      }



      break
    }
  }

}

