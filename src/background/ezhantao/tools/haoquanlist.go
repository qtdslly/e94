package main

import (
  "background/ezhantao/service"
  "background/wechart/config"
  "background/common/util"
  "background/common/logger"


  "github.com/tidwall/gjson"
  "fmt"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())

  for i := 0 ; i < 100 ; i++{
    logger.Debug("第",i+1 ,"页================================================")
    GetCouponInfo(i + 1,100)
  }
}


func GetCouponInfo(page,size int){
  recv := service.GetHaoQuanList(page,size,"勺子")
  //logger.Debug(recv)
  list := gjson.Get(recv, "tbk_dg_item_coupon_get_response.results.tbk_coupon")
  if list.Exists() {
    items := list.Array()
    for _, item := range items {

      title := item.Get("title").String()
      item_url := item.Get("item_url").String()

      //商品价格
      finalPrice := item.Get("zk_final_price").String()
      //商品id
      item_id := item.Get("num_iid").String()
      //领券地址
      coupon_click_url := item.Get("coupon_click_url").String()

      //logger.Debug(finalPrice)
      response := service.GetTuiGuangQuan("",item_id,"")

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

      fp := util.GetAmt(finalPrice)
      ca := util.GetAmt(coupon_amount)
      if fp - ca == 0 {
        logger.Debug(fp-ca,"|",finalPrice,"|",coupon_start_fee,"|",coupon_amount,"|",scene,"|",couponType,"|",item_url,"|",title,"|",coupon_click_url)
      }

      //res := service.GetXuanPinKuGoodsInfo1(item_id)
      //logger.Debug(res)
      //coupon_type := gjson.Get(response, "tbk_coupon_get_response.data.coupon_type").Int()


      //title := res.Get("title").String()

    }
  }
}
