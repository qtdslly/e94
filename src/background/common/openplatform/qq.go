package openplatform

import (
	"common/logger"
	"common/util"
	"net/url"
)

/*
const (
	qqAppKey = "100735919"
)
*/

type QqClient struct {
	AppId string
}

func InitQq(appId string) *QqClient {
	client := new(QqClient)
	client.AppId = appId
	return client
}

/*
	Erroe message:
	{
	  "ret": 100008,
	  "msg": "client request's app is not existed"
	}
*/
type QqUserInfo struct {
	ErrCode    int    `json:"ret"`         // return
	ErrMsg     string `json:"msg"`         // message
	Nickname   string `json:"nickname"`    // 昵称
	Gender     string `json:"gender"`      // 性别, "男", "女"
	Province   string `json:"province"`    // 省份
	City       string `json:"city"`        // 城市
	Year       string `json:"year"`        // 生日年份
	Figureurl1 string `json:"figureurl_2"` // 100*100头像
}

func (c *QqClient) GetQqUserInfo(openid, accessToken string) *QqUserInfo {
	values := url.Values{}
	values.Add("openid", openid)
	values.Add("access_token", accessToken)
	values.Add("oauth_consumer_key", c.AppId)
	values.Add("format", "json")
	url := "https://graph.qq.com/user/get_user_info?" + values.Encode()

	var user QqUserInfo
	var err error
	if err = util.Get(url, &user, nil); err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("qq user info: ", user)

	if user.ErrCode == 0 && user.Nickname != "" {
		return &user
	}

	return nil
}
