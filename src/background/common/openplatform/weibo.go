package openplatform

import (
	"common/logger"
	"common/util"
	"net/url"
)

/*
	Error message:

	{
	  "error": "source paramter(appkey) is missing",
	  "error_code": 10006,
	  "request": "/2/users/show.json"
	}
*/
type WeiboUserInfo struct {
	ErrCode         int    `json:"error_code"`        // error_code
	ErrMsg          string `json:"error"`             // error
	Request         string `json:"request"`           // request
	Idstr           string `json:"idstr"`             // 字符串型的用户UID
	ScreenName      string `json:"screen_name"`       // 用户昵称
	Name            string `json:"name"`              // 友好显示名称
	Location        string `json:"location"`          // 用户所在地
	Description     string `json:"description"`       // 用户个人描述
	ProfileImageUrl string `json:"profile_image_url"` // 用户头像地址（中图），50×50像素
	ProfileUrl      string `json:"profile_url"`       // 用户的微博统一URL地址
	Weihao          string `json:"weihao"`            // 用户的微号
	Gender          string `json:"gender"`            // 性别，m：男、f：女、n：未知
	AvatarLarge     string `json:"avatar_large"`      // 用户头像地址（大图），180×180像素
	AvatarHd        string `json:"avatar_hd"`         // 用户头像地址（高清），高清头像原图
}

func GetWeiboUserInfo(uid, accessToken string) *WeiboUserInfo {
	values := url.Values{}
	values.Add("uid", uid)
	values.Add("access_token", accessToken)
	url := "https://api.weibo.com/2/users/show.json?" + values.Encode()

	var user WeiboUserInfo
	var err error
	if err = util.Get(url, &user, nil); err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("weibo user info: ", user)

	if user.ErrCode == 0 && user.Idstr != "" {
		return &user
	}

	return nil
}
