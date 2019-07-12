package service

import (
  "time"
  "fmt"
  "crypto/md5"
  "bytes"
  "strconv"
  "strings"
  "net/url"
  "net/http"
  "io/ioutil"
  "sort"

  "background/common/logger"
  "background/ezhantao/common/constant"
)

//生成淘口令
func GetTPwd(click_url ,text,logo string)(string){
  body := GetPublicParm()
  body["user_id"] = "53171692"   //页数
  body["text"] = "【天天特价】" + text  //口令弹框内容
  body["url"] = click_url   //口令跳转目标页
  body["logo"] = logo   //页数
  body["ext"] = "{}"   //扩展字段JSON格式
  logger.Debug(text)
  logger.Debug(click_url)
  logger.Debug(logo)

  body["method"] = constant.TAOBAO_API_GET_TPWD
  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}

func ShareQuery(password_content string)(string){
  body := GetPublicParm()
  body["password_content"] = password_content   //淘口令  【天猫品牌号】，复制这条信息￥sMCl0Yra3Ae￥后打开手机淘宝
  body["method"] = constant.TAOBAO_API_EXTRACT
  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}


func GetExtract(click_url string)(string){
  body := GetPublicParm()
  body["click_url"] = click_url   //页数
  body["method"] = constant.TAOBAO_API_EXTRACT
  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}

func GetHaoQuanList(page,size int,title string)(string){
  body := GetPublicParm()
  body["page_no"] = fmt.Sprintf("%d",page)   //页数
  body["page_size"] = fmt.Sprintf("%d",size) //每页数量
  body["adzone_id"] = "109081500310"  //推广位id
  body["q"] = title  //查询词
  body["platform"] = "2"  //1：PC，2：无线，默认：1
  body["cat"] = ""  //后台类目ID，用,分割，最大10个，该ID可以通过taobao.itemcats.get接口获取到

  body["method"] = constant.TAOBAO_API_HAOQUAN_LIST_GET
  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}

func GetXuanPinKuGoodsInfo1(item_id string)(string){
  body := GetPublicParm()
  body["num_iids"] = item_id    //商品ID串，用,分割，最大40个
  body["platform"] = "" //链接形式：1：PC，2：无线，默认：１
  body["ip"] = "221.235.86.242" //ip地址，影响邮费获取，如果不传或者传入不准确，邮费无法精准提供
  body["method"] = constant.TAOBAO_API_GOODS_GET1

  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}

func GetXuanPinKuGoodsInfo()(string){
  body := GetPublicParm()
  body["adzone_id"] = "109081500310"
  body["favorites_id"] = "19560515" //9.9包邮
  body["method"] = constant.TAOBAO_API_XUANPINKU_GOODS_GET
  body["fields"] = "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url,seller_id,volume,nick"

  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}

//根据nid批量查询优惠券
//pid 三方pid，满足mm_xxx_xxx_xxx格式
//num_iids 商品ID串，用逗号分割，从taobao.tbk.item.coupon.get接口获取num_iid字段，最大40个。（如果传入了没有优惠券的宝贝num_iid，则优惠券相关的字段返回为空，请做好容错）
func GetCoupon(item_id string)(string){
  body := GetPublicParm()
  //logger.Debug(me)
  body["pid"] = "mm_53171692_584850166_109081500310"
  body["item_id"] = item_id
  body["platform"] = "2"

  body["method"] = constant.TAOBAO_API_TUIGUANGQUAN_GET

  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)
}



//阿里妈妈推广券信息查询
//me 带券ID与商品ID的加密串
//item_id 商品ID
//activity_id 券ID
func GetTuiGuangQuan(me,item_id,activity_id string)(string){
  body := GetPublicParm()
  //logger.Debug(me)
  body["me"] = me
  body["item_id"] = item_id
  body["activity_id"] = activity_id

  body["method"] = constant.TAOBAO_API_TUIGUANGQUAN_GET

  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)
}


func GetXuanPinKuList()(string){
  body := GetPublicParm()
  body["method"] = constant.TAOBAO_API_XUANPINKU_LIST_GET
  body["fields"] = "favorites_title,favorites_id,type"

  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)

}


//淘宝客商品查询
func GetGoodsInfo(title string)(string){

  body := GetPublicParm()
  body["q"] = title
  body["method"] = constant.TAOBAO_API_GOODS_GET
  body["fields"] = "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url,seller_id,volume,nick"

  sign := GetSign(body)
  //logger.Debug(sign)
  return ToTaobao(constant.TAOBAO_API_URL,body,sign)
}


func GetPublicParm()map[string]string{
  p := time.Now()
  timeStamp := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",p.Year(),p.Month(),p.Day(),p.Hour(),p.Minute(),p.Second())
  //logger.Debug(timeStamp)
  body := map[string]string{"sign_method":"md5","app_key":constant.EZHANTAO_APP_KEY,"format":"json","timestamp":timeStamp,"v":"2.0"}
  return body
}


func ToTaobao(apiUrl string, body map[string]string,sign string) string {
  data := url.Values{}
  for key, val := range body {
    data.Add(key, val)
  }
  data.Add("sign", sign)

  requ, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
  requ.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

  resp, err := http.DefaultClient.Do(requ)
  if err != nil {
    logger.Debug(err)
    return ""
  }

  recv, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    logger.Error(err)
    return ""
  }
  //logger.Debug(string(recv))
  return string(recv)
}


func GetSign(body map[string]string)string{
  /**
  1、第一步：把字典按Key的字母顺序排序
  2、第二步：把所有参数名和参数值串在一起
  3、第三步：使用MD5/HMAC加密
  4、第四步：把二进制转化为大写的十六进制
  */

  //第一步：把字典按Key的字母顺序排序
  var names []string
  for name,_ := range body {
    if name == "sign"{
      continue
    }
    names = append(names, name)
  }
  sort.Strings(names)
  //第二步：把所有参数名和参数值串在一起
  var bb bytes.Buffer
  bb.WriteString(constant.EZHANTAO_APP_SECRET)
  for _, v := range names {
    val := body[v]
    //if len(val) > 0 {
      bb.WriteString(v)
      bb.WriteString(val)
    //}
  }
  //第三步：使用MD5/HMAC加密
  b := make([]byte, 0)

  bb.WriteString(constant.EZHANTAO_APP_SECRET)
  //fmt.Println(bb.String())

  md5instence := md5.New()
  md5instence.Write(bb.Bytes())
  b = md5instence.Sum(nil)
  //第四步：把二进制转化为大写的十六进制
  var result bytes.Buffer
  for i := 0; i < len(b); i++ {
    s := strconv.FormatInt(int64(b[i]&0xff), 16)
    if len(s) == 1 {
      result.WriteString("0")
    }
    result.WriteString(s)
  }

  //返回签名完成的字符串
  return strings.ToUpper(result.String())
}












