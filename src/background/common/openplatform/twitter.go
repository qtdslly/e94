package openplatform

import (
	"common/logger"

	"github.com/ChimeraCoder/anaconda"
)

/*
const (
	twitterConsumerKet    = "iUbiLhqrLGxEjbO42FAtvyfos"
	twitterConsumerSecret = "QII7hima61czqvMcTzysGwd2drILePLyTCWMzd7H4SvhYEPYGM"
)
*/

type TwitterClient struct {
	TwitterConsumerKet    string
	TwitterConsumerSecret string
}

func InitTwitter(twConsumerKet, twConsumerSecret string) *TwitterClient {
	client := new(TwitterClient)
	client.TwitterConsumerKet = twConsumerKet
	client.TwitterConsumerSecret = twConsumerSecret
	return client
}

/*
	Error message:

	{
	  "errors": [
	    {
	      "code": 215,
	      "message": "Bad Authentication data."
	    }
	  ]
	}
*/
type TwitterUserInfo struct {
	Errors []*struct {
		ErrCode int    `json:"code"`    // err_code
		ErrMsg  string `json:"message"` // err_msg
	} `json:"errors"` // error
	Id         string // id
	Name       string // 名称
	ScreenName string // 用户唯一标识
	Location   string // 地点
	Avatar     string // 头像
}

func (c *TwitterClient) GetTwitterUserInfo(token, secret string) *TwitterUserInfo {
	anaconda.SetConsumerKey(c.TwitterConsumerKet)
	anaconda.SetConsumerSecret(c.TwitterConsumerSecret)
	api := anaconda.NewTwitterApi(token, secret)

	twitterUser, err := api.GetSelf(nil)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if twitterUser.Id == 0 {
		logger.Error("twitter user info is empty !")
		return nil
	}

	var user TwitterUserInfo
	user.Id = twitterUser.IdStr
	user.Name = twitterUser.Name
	user.ScreenName = twitterUser.ScreenName
	user.Location = twitterUser.Location
	user.Avatar = twitterUser.ProfileImageURL

	logger.Debug("twitter user info: ", user)

	return &user
}
