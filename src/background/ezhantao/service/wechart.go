package service

import (
  "background/common/logger"
  "github.com/tidwall/gjson"
  "background/common/util"
  "encoding/xml"
  "time"
)

type TextResMessage struct {
  XMLName       xml.Name    `xml:"xml"`
  ToUserName    string      `xml:"ToUserName"`
  FromUserName  string      `xml:"FromUserName"`
  CreateTime    int64       `xml:"CreateTime"`
  MsgType       string      `xml:"MsgType"`
  Content       string      `xml:"Content"`
}

func SearchGoods(toUserName,fromUserName,title string)*TextResMessage{
  recv := GetHaoQuanList(1,100,title)
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

      response := GetTuiGuangQuan("",item_id,"")
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

  response := GetTPwd(url,title,logo)
  pwd := gjson.Get(response, "tbk_tpwd_create_response.data.model").String()

  var news TextResMessage
  news.CreateTime = time.Now().Unix()
  news.ToUserName = toUserName
  news.FromUserName = fromUserName
  news.MsgType = "text"
  if pwd == ""{
    news.Content = "抱歉，未找到相关的购物优惠券,快快加入QQ群825744743，每天发布大额优惠券，最高还可0元购买宝贝"
  }else{
    news.Content = "快快复制本段内容到淘宝领取购物券吧，淘口令 " + pwd
  }

  return &news
}


type WxMessage struct {
  XMLName       xml.Name    `xml:"xml"`
  ToUserName    string      `xml:"ToUserName"`
  FromUserName  string      `xml:"FromUserName"`
  CreateTime    int64       `xml:"CreateTime"`
  MsgType       string      `xml:"MsgType"`
  Content       string      `xml:"Content"`
}



func EventContent(eventType ,toUserName,fromUserName string) *WxMessage {
  var wm WxMessage
  wm.ToUserName = toUserName
  wm.FromUserName = fromUserName
  wm.CreateTime = time.Now().Unix()
  wm.MsgType = "text"
  if eventType == "subscribe"{
    wm.Content = "欢迎关注e站淘，购物复制宝贝名称发送到此订阅号，领取隐藏的购物券，你也可以加入QQ群825744743，每天领取大额优惠券，最高还可0元购买宝贝"
  }else {
    wm.Content = "谢谢你曾经的关注与支持，我们会继续努力做到更好，期待您的再次关注！！！"
  }
  return &wm
}
