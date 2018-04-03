package openplatform

import (
	"common/logger"
	"common/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type FacebookClient struct {
	FacebookAppId     string
	FacebookAppSecret string
}

func InitFacebook(facebookId, facebookSecret string) *FacebookClient {
	client := new(FacebookClient)
	client.FacebookAppId = facebookId
	client.FacebookAppSecret = facebookSecret
	return client
}

type facebookAccessToken struct {
	Token  string
	Expire int64
}

/*
	Error message:

	{
	  "error": {
	    "message": "An active access token must be used to query information about the current user.",
	    "type": "OAuthException",
	    "code": 2500,
	    "fbtrace_id": "ET+boSJbKWT"
	  }
	}

	{
	  "id": "1679078442384055",
	  "name": "Alwin Yu",
	  "first_name": "Alwin",
	  "last_name": "Yu",
	  "age_range": {
	    "min": 21
	  },
	  "link": "https://www.facebook.com/app_scoped_user_id/1679078442384055/",
	  "gender": "male",
	  "locale": "zh_CN",
	  "picture": {
	    "data": {
	      "is_silhouette": false,
	      "url": "https://fb-s-d-a.akamaihd.net/h-ak-xfp1/v/t1.0-1/p50x50/10676182_1375409062750996_814084756399504591_n.jpg?oh=8a25fc80f9abe421a8bbdbcace6d2c18&oe=58DAAD34&__gda__=1491987540_ff5cd75924094765f19c316c4a8849e7"
	    }
	  }
	}

*/
type FacebookUserInfo struct {
	Error *struct {
		Type    string `json:"type"`    // err_code
		ErrCode int    `json:"code"`    // err_code
		ErrMsg  string `json:"message"` // err_msg
	} `json:"error"` // error
	Id         string `json:"id"`         // 昵称
	Name       string `json:"name"`       // 昵称
	First_name string `json:"first_name"` // 昵称
	LastName   string `json:"last_name"`  // 昵称
	AgeRange   *struct {
		Min uint32 `json:"min"`
	} `json:"age_range"`
	Gender  string `json:"gender"` // 性别, "男", "女"
	Locale  string `json:"locale"` // 省份
	Picture *struct {
		Data *struct {
			IsSilhouette bool   `json:"is_silhouette"`
			Url          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

/*
	retrieve the access_token from a code.
*/
func (c *FacebookClient) getFacebookAccessToken(code string) *facebookAccessToken {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", c.FacebookAppId)
	values.Add("client_secret", c.FacebookAppSecret)
	// values.Add("redirect_uri", config.GetApiBindAddr()+facebookRedirectUri)
	url := "https://graph.facebook.com/oauth/access_token?" + values.Encode()

	/*
		var token facebookAccessToken
		var err error
		if err = util.Get(url, &token); err != nil {
			logger.Error(err)
			return nil
		}
	*/

	//https://graph.facebook.com/oauth/access_token?client_id=YOUR_APP_ID&redirect_uri=YOUR_REDIRECT_URI&client_secret=YOUR_APP_SECRET&code=CODE_GENERATED_BY_FACEBOOK
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}

	var token facebookAccessToken
	tokenArr := strings.Split(string(b), "&")
	token.Token = strings.Split(tokenArr[0], "=")[1]
	expireInt, err := strconv.Atoi(strings.Split(tokenArr[1], "=")[1])
	if err == nil {
		token.Expire = int64(expireInt)
	}
	return &token
}

func (c *FacebookClient) getFacebookUserInfo(code string) *FacebookUserInfo {
	accessToken := c.getFacebookAccessToken(code)
	resp, err := http.Get("https://graph.facebook.com/me?access_token=" + accessToken.Token)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("facebook user info: ", string(out))

	var user FacebookUserInfo
	if err = json.Unmarshal(out, &user); err != nil {
		logger.Error(err)
		return nil
	}

	// 	img := "https://graph.facebook.com/" + user.Id + "/picture?width=180&height=180"
	return &user
}

/*
	Only this method is used, client will pass the access_token instad of authorization code.
*/
func GetFacebookUserInfo(accessToken string) *FacebookUserInfo {
	values := url.Values{}
	values.Add("access_token", accessToken)
	values.Add("fields", "id,name,first_name,last_name,age_range,link,gender,locale,picture")
	url := "https://graph.facebook.com/me?" + values.Encode()

	var user FacebookUserInfo
	var err error
	if err = util.Get(url, &user, nil); err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("facebook user info: ", user)

	if user.Error == nil && user.Id != "" {
		return &user
	}
	// img := "https://graph.facebook.com/" + user.Id + "/picture?width=180&height=180"
	return nil
}
