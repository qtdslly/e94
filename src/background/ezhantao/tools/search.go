package main

import (
  "background/ezhantao/service"
  "background/wechart/config"
  "background/common/logger"
)

func main(){
  logger.SetLevel(config.GetLoggerLevel())

  //recv := service.GetGoodsInfo()

  //recv := service.GetXuanPinKuList()
  //logger.Debug(recv)
  //
  //recv = service.GetXuanPinKuGoodsInfo()
  //logger.Debug(recv)
  //
  //recv = service.GetTuiGuangQuan()
  //logger.Debug(recv)

  //recv := service.GetHaoQuanList()
  //logger.Debug(recv)
  //
  //recv := service.GetExtract("https://uland.taobao.com/coupon/edetail?e=bjhgxXUhiywGQASttHIRqRrmJquQ%2FjbFIXUwMB3XPv3l43M3mIB1t7o46G7PUKik03VEUTctUDvqED9mu%2B62cKdDpfCY3RElmtj%2FjdqaYAdiMp1hCSxiEd%2F%2FIrM%2BSQPg0Li%2B7vIEjtYK6zASXfQ%2BWAMFTfGLA%2FdyF9wuEQDlqL5diwTmD3eVNiZ6Y%2FpkHtT5QS0Flu%2FfbSp4QsdWMikAauib2QyKSVw%2FCF5GnkyyjY0%3D&traceId=0bb69b7d15616056334838787e")
  //logger.Debug(recv)

  //recv := service.ShareQuery("【现货 如假白送！Converse 1970s White 白色低帮 149448C 162065C】https://m.tb.cn/h.egiLqoD?sm=cebc3c 点击链接，再选择浏览器咑閞；或復·制这段描述￥uuKZY5smpno￥后到淘♂寳♀")
  //logger.Debug(recv)

  //recv := service.GetCoupon("37016930337")
  //logger.Debug(recv)

}

