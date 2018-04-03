package openplatform

import (
	"common/logger"
	"common/util"
	"net/url"
)

/*
const (
	weixinAppId     = "wx903a16f5ed5f9510"
	weixinAppSecret = "df598b6ba30f6d40e3cbaaaf00f681d2"
)
*/

type WeixinClient struct {
	AppId     string
	AppSecret string
}

func InitWeixin(wxAppId, wxAppSecret string) *WeixinClient {
	client := new(WeixinClient)
	client.AppId = wxAppId
	client.AppSecret = wxAppSecret
	return client
}

type weixinAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

/*
	Error message:
	{
	  "errcode": 40029,
	  "errmsg": "invalid code, hints: [ req_id: Oq1BNa0722ns86 ]"
	}
*/
type WeixinUserInfo struct {
	ErrCode    int      `json:"errcode"`    // errcode
	ErrMsg     string   `json:"errmsg"`     // errmsg
	OpenId     string   `json:"openid"`     // 普通用户的标识，对当前开发者帐号唯一
	Nickname   string   `json:"nickname"`   // 普通用户昵称
	Sex        int      `json:"sex"`        // 普通用户性别，1为男性，2为女性
	Province   string   `json:"province"`   // 普通用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	Headimgurl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0为默认值，代表640*640正方形头像），用户没有头像时该项为空
	Privilege  []string `json:"privilege"`  // 用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	UnionId    string   `json:"unionid"`    // 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
}

func getWeixinAccessToken(code string, c *WeixinClient) *weixinAccessToken {
	values := url.Values{}
	values.Add("appid", c.AppId)
	values.Add("secret", c.AppSecret)
	values.Add("code", code)
	values.Add("grant_type", "authorization_code")
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?" + values.Encode()

	var response weixinAccessToken
	var err error
	if err = util.Get(url, &response, nil); err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("weixin access token: ", &response, url)

	return &response
}

func (c *WeixinClient) GetWeixinUserInfoByAccessToken(openid, accessToken string) *WeixinUserInfo {
	values := url.Values{}
	values.Add("openid", openid)
	values.Add("access_token", accessToken)
	values.Add("lang", "zh_CN")
	url := "https://api.weixin.qq.com/sns/userinfo?" + values.Encode()

	var user WeixinUserInfo
	var err error
	if err = util.Get(url, &user, nil); err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("weixin user info: ", user, url)

	if user.ErrCode == 0 && user.OpenId != "" {
		return &user
	}

	return nil
}

func (c *WeixinClient) GetWeixinUserInfo(code string) *WeixinUserInfo {
	accessToken := getWeixinAccessToken(code, c)
	if accessToken == nil {
		return nil
	}

	return c.GetWeixinUserInfoByAccessToken(accessToken.Openid, accessToken.AccessToken)
}
