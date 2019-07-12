package constant

const (
  //App key 和App ID
  EZHANTAO_APP_KEY = "27649816"
  EZHANTAO_APP_SECRET = "38eee0c63fe0e44947af40a61d9c8bd2"

  //淘宝API地址，方法
  TAOBAO_API_URL = "http://gw.api.taobao.com/router/rest"
  TAOBAO_API_GOODS_GET = "taobao.tbk.item.get"  //淘宝客商品查询
  TAOBAO_API_XUANPINKU_GOODS_GET = "taobao.tbk.uatm.favorites.item.get"  //获取淘宝联盟选品库的宝贝信息
  TAOBAO_API_XUANPINKU_LIST_GET = "taobao.tbk.uatm.favorites.get"  //获取淘宝联盟选品库列表



  TAOBAO_API_GOODS_GET1 = "taobao.tbk.item.info.get"  //淘宝客商品详情（简版）


  TAOBAO_API_HAOQUAN_LIST_GET = "taobao.tbk.dg.item.coupon.get"  //好券清单API【导购】
  TAOBAO_API_TUIGUANGQUAN_GET = "taobao.tbk.coupon.get"  //阿里妈妈推广券信息查询


  TAOBAO_API_EXTRACT = "taobao.tbk.item.click.extract"  //链接解析api


  TAOBAO_API_GET_COUPON = "taobao.tbk.itemid.coupon.get"  //根据nid批量查询优惠券

  TAOBAO_API_GET_TPWD = "taobao.tbk.tpwd.create"  //淘宝客淘口令



  //===================收费接口=====================
  TAOBAO_API_SHARE_QUERY = "taobao.wireless.share.tpwd.query"  //查询解析淘口令

)
