package setting

import (
	"common/constant"
	"common/logger"
	"common/openplatform"
	"component/kv"
	"encoding/json"
	"sync"

	"github.com/jinzhu/gorm"
)

/*
	conventioned cert file names, recreated under ConfigPath with file content stored in DB.
*/
const (
	filenameAlipayRSAPrivateKey = "alipay_public_key_%d.pem"
	filenameAlipayRSAPublicKey  = "alipay_wap_pkcs8_rsa_private_key_%d.pem"
)

var (
	clientMap map[uint32]*platformObjectSetting // key为app_id value为第三方平台的参数集合
	rwClient  *sync.RWMutex
)

type platformObjectSetting struct {
	Wx      *openplatform.WeixinClient
	Qq      *openplatform.QqClient
	Google  *openplatform.GoogleClient
	Twitter *openplatform.TwitterClient
}

type OpenplatformSetting struct {
	FacebookAppId         string `json:"facebook_app_id"`
	FacebookAppSecret     string `json:"facebook_app_secret"`
	GoogleClientId        string `json:"google_client_id"`
	GoogleClientSecret    string `json:"google_client_secret"`
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	WeixinAppId           string `json:"weixin_app_id"`
	WeixinAppSecret       string `json:"weixin_app_secret"`
	QqAppId               string `json:"qq_app_id"`
	AlipayPartnerId       string `json:"alipay_partner_id"`
	AlipaySellerId        string `json:"alipay_seller_id"`
	PaypalSellerEmail     string `json:"paypal_seller_email"`
	WxpayAppId            string `json:"wxpay_app_id"`
	WxpayMerchantId       string `json:"wxpay_merchant_id"`
	WxpayAppKey           string `json:"wxpay_app_key"`
	AccessKey             string `json:"access_key"`
	AlipayRSAPrivateKey   string `json:"alipay_rsa_private_key"`
	AlipayRSAPublicKey    string `json:"alipay_rsa_public_key"`
}

func init() {
	clientMap = make(map[uint32]*platformObjectSetting)
	rwClient = new(sync.RWMutex)
}

// func initOpenplatformSetting (db *gorm.DB) error {
// 	thirdPartyCache, err := getOpenplatformSetting (db)
// 	if err != nil && err.Error() == kv.KeyNotFoundError {
// 		// TODO: value change to array
// 		var info OpenplatformSetting
// 		info.AlipayRSAPrivateKey = `-----BEGIN PRIVATE KEY-----
// MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKaVqGaS+Q74fqUk
// 2jHeHz6TJNX2eH+oGDyk3IS3biIAYkalZlXwfmlcvNiV9ajz2IRtfWInzexnjz0l
// tBw0atMSB7ywj3wELP6wR7dF09//ZRRHfq7liLqIg6ACX7wQirLJsS4yCVeb0RQE
// HS7SvqtAidYXrak9sUD8Ig2vVUhxAgMBAAECgYAdGf3+VWSvKdguD38SwSQiMWB5
// BILOglYhmgdwI/9Yb0q73r8++jGLRIk1KRVue5Lyp5IE6ME/sGPEoeaSGtIiX0O0
// lErKuFQM5UQmg5AA23aLHWDR6hh+DqpSE6pNOW38CtKIu8Vq8Ujpv16mBi07d9Ew
// 0t7Ltjpz8UDA+mvsQQJBANhdwEVyg+bmwZSkDu1oBok10nQrFFN7+kl123RSAJAk
// Qu5Gsvey9eameMRpHWefbIH7+cF7xZio5qACdVIx9DkCQQDFGXRKZGMJjtYzNKrV
// cm+wMy34tslOCzogl0mqG2WZr0Z+z2bwxQDOUx9BkqGIYraXz9SwLh6m+WmjaftZ
// BqX5AkEAx5CxZ1zYjIEKzC8GFbN9U9Fw6/VQKjCQCnHKmN+J5WmM0nojWJSYesXR
// XlgV5x4E11+yXePrvYMMNUoPLGEnAQJAbH2oX1KGmTpAeYsiBb+p7skdIuwqPoU7
// h7j+2V2fPUsXeLHdLeainO9wIv39YD9F1qaVoiygvrRHC6ZIriZUsQJBAJnQ5co3
// iFotP3QQdin1tT1PHm4IE3dsk4nrrdlSJ9B1WQgm79bA5tMg9Tl22p4ZYO6wohUc
// m4w2M2ubI5RoPd0=
// -----END PRIVATE KEY-----`
// 		info.AlipayRSAPublicKey = `-----BEGIN PUBLIC KEY-----
// MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnxj/9qwVfgoUh/y2W89L6BkRA
// FljhNhgPdyPuBV64bfQNN1PjbCzkIM6qRdKBoLPXmKKMiFYnkd6rAoprih3/PrQE
// B/VsW8OoM8fxn67UDYuyBTqA23MML9q1+ilIZwBC2AQ2UBVOrFXfFl75p6/B5Ksi
// NG9zpgmLCUYuLkxpLQIDAQAB
// -----END PUBLIC KEY-----`
// 		info.FacebookAppId = "1070700103040023"
// 		info.FacebookAppSecret = "91105fec4f040e4d1bd73261454dfe05"
// 		info.GoogleClientId = "158949607299-ohu7bnt1m7ud95c76ct6tcd7aklbmuou.apps.googleusercontent.com"
// 		info.GoogleClientSecret = "bUgT8aXQGleQw_CVbEF7bSxb"
// 		info.TwitterConsumerKey = "qJH8t2qon70hfQWwgOgzOlkcW"
// 		info.TwitterConsumerSecret = "vVXJqnaHcsLfSPljMEItlufY8pr0Tr5JBlBHt04CgEGXBmhkCO"
// 		info.WeixinAppId = "wx903a16f5ed5f9510"
// 		info.WeixinAppSecret = "df598b6ba30f6d40e3cbaaaf00f681d2"
// 		info.QqAppId = "100735919"
// 		info.AlipayPartnerId = "2088111854386315"
// 		info.AlipaySellerId = "guoguangxingkong@163.com"
// 		info.PaypalSellerEmail = "info-facilitator@starschina.com"
// 		info.WxpayAppId = "wx903a16f5ed5f9510"
// 		info.WxpayAppKey = "2557A8332CEB8C3438491D73A65F7573"
// 		info.WxpayMerchantId = "1239565402"
// 		err = UpdateOpenplatformSetting (&info, db)
// 	} else {
// 		saveKeyFiles(thirdPartyCache)
// 		updateThirdParty(thirdPartyCache, db)
// 	}
// 	return err
// }

func getOpenplatformSetting(appId uint32, db *gorm.DB) (*OpenplatformSetting, error) {
	value, err := kv.GetValueForKey(appId, 0, constant.ThirdPartyInfoKey, false, db)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return &OpenplatformSetting{}, nil
	}

	setting := &OpenplatformSetting{}
	err = json.Unmarshal([]byte(value), setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func UpdateOpenplatformSetting(appId uint32, info *OpenplatformSetting, db *gorm.DB) error {
	value, err := json.Marshal(*info)
	if err != nil {
		return err
	}

	err = kv.SetValueForKey(appId, 0, constant.ThirdPartyInfoKey, string(value), db)
	if err != nil {
		return err
	}

	s, err := newOpSetting(info, db)
	if err != nil {
		logger.Error(err)
		return err
	}

	rwClient.Lock()
	clientMap[appId] = s
	rwClient.Unlock()

	return nil
}

// 初始化配置
func newOpSetting(info *OpenplatformSetting, db *gorm.DB) (*platformObjectSetting, error) {
	s := new(platformObjectSetting)
	// initialize the third party login setting
	s.Wx = openplatform.InitWeixin(info.WeixinAppId, info.WeixinAppSecret)
	s.Qq = openplatform.InitQq(info.QqAppId)
	s.Google = openplatform.InitGoogle(info.GoogleClientId, info.GoogleClientSecret)
	s.Twitter = openplatform.InitTwitter(info.TwitterConsumerKey, info.TwitterConsumerSecret)
	return s, nil
}

// 根据app_id获取设置
func GetOpSetting(appId uint32, db *gorm.DB) (*platformObjectSetting, error) {
	v, ok := clientMap[appId]
	if !ok {
		info, err := getOpenplatformSetting(appId, db)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		v, err = newOpSetting(info, db)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		rwClient.Lock()
		clientMap[appId] = v
		rwClient.Unlock()
	}
	return v, nil
}
