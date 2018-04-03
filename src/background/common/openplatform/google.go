package openplatform

import (
	"common/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GoogleClient struct {
	ClientId     string
	ClientSecret string
}

func InitGoogle(googleId, googleSecret string) *GoogleClient {
	client := new(GoogleClient)
	client.ClientId = googleId
	client.ClientSecret = googleSecret
	return client
}

/*
	{
	  "access_token": "ya29.Ci_DA2yjlM160APEQ7ZM-xXhW6tD_eZmO4gBpU6gSwakHsUgiEewcozC3mGzrcUaDw",
	  "token_type": "Bearer",
	  "expires_in": 3600,
	  "refresh_token": "1/7vQvpJhjGpKXfKMeebNjSuS01n6ToPD3U8PinJWoEos",
	  "id_token": "**-cA"
	}
*/
type googleAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

/*
	Error message:
	"error": {
	    "errors": [
	      {
		"domain": "global",
		"reason": "authError",
		"message": "Invalid Credentials",
		"locationType": "header",
		"location": "Authorization"
	      }
	    ],
	    "code": 401,
	    "message": "Invalid Credentials"
	  }


	{
	  "id": "101093901634076996428",
	  "name": "wei yu",
	  "given_name": "wei",
	  "family_name": "yu",
	  "link": "https://plus.google.com/101093901634076996428",
	  "picture": "https://lh5.googleusercontent.com/-440WIbSGIuc/AAAAAAAAAAI/AAAAAAAAABQ/QmuvQYmKH0g/photo.jpg",
	  "gender": "other",
	  "locale": "zh-CN"
	}
*/
type GoogleUserInfo struct {
	Error *struct {
		ErrCode int    `json:"code"`    // err_code
		ErrMsg  string `json:"message"` // err_msg
	} `json:"error"` // error
	Id         string `json:"id"`          // 昵称
	Name       string `json:"name"`        // 姓名
	GivenName  string `json:"given_name"`  // 姓
	FamilyName string `json:"family_name"` // 名
	Picture    string `json:"picture"`     // 头像
	Gender     string `json:"gender"`      // 性别, "男", "女"
	Locale     string `json:"locale"`      // 地区
}

/*
	retrieve the access_token from a code.
*/
func (c *GoogleClient) getGoogleAccessToken(code string, scope, grantType string) *googleAccessToken {
	values := url.Values{}
	values.Add("code", code)
	values.Add("redirect_uri", "")
	values.Add("client_id", c.ClientId)
	values.Add("client_secret", c.ClientSecret)
	if scope != "" {
		values.Add("scope", scope)
	}
	values.Add("grant_type", grantType)
	url := "https://www.googleapis.com/oauth2/v4/token?" + values.Encode()

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("google access token body::", string(body))

	var token googleAccessToken
	if err = json.Unmarshal(body, &token); err != nil {
		logger.Error(err)
		return nil
	}
	return &token
}

func (c *GoogleClient) GetGoogleUserInfo(code string) *GoogleUserInfo {
	// step 1, get access token
	scope := ""
	grantType := "authorization_code"
	token := c.getGoogleAccessToken(code, scope, grantType)

	if token == nil || token.AccessToken == "" {
		return nil
	}

	// step 2, retrive user information
	return c.GetGoogleUserInfoByAccessToken(token.AccessToken)
}

func (c *GoogleClient) GetGoogleUserInfoByAccessToken(accessToken string) *GoogleUserInfo {
	// retrive user information
	url := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err)
		return nil
	}

	req.Header.Add("authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}

	logger.Debug("google user info: ", string(body))

	var user GoogleUserInfo
	if err = json.Unmarshal(body, &user); err != nil {
		logger.Error(err)
		return nil
	}

	if user.Error == nil && user.Id != "" {
		return &user
	}

	return nil
}
