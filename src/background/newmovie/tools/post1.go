package main

import (
	"net/http"
	"encoding/json"
	"background/common/logger"
	"io/ioutil"
	"background/wechart/config"
	"background/wechart/service"
	"bytes"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	accessToken := service.AccessToken()

	logger.Debug(accessToken)
	/*
	{
	     "button":[
	     {
		  "type":"click",
		  "name":"今日歌曲",
		  "key":"V1001_TODAY_MUSIC"
	      },
	      {
		   "name":"菜单",
		   "sub_button":[
		   {
		       "type":"view",
		       "name":"搜索",
		       "url":"http://www.soso.com/"
		    },
		    {
			 "type":"miniprogram",
			 "name":"wxa",
			 "url":"http://mp.weixin.qq.com",
			 "appid":"wx286b93c14bbf93aa",
			 "pagepath":"pages/lunar/index"
		     },
		    {
		       "type":"click",
		       "name":"赞一下我们",
		       "key":"V1001_GOOD"
		    }]
	       }]
	 }

	*/
	type SubMenu struct {
		MenuType      string     `json:"type"`
		Name          string     `json:"name"`
		Url           string     `json:"url"`
		AppId         string     `json:"appid"`
		PagePath      string     `json:"pagepath"`
		Key           string     `json:"key"`
	}
	type FirstMenu struct {
		MenuType      string     `json:"type"`
		Name          string     `json:"name"`
		Key           string     `json:"key"`
		SubMenus      []SubMenu  `json:"sub_button"`
	}
	type Menu struct {
		FirstMenus       []FirstMenu   `json:"button"`
	}

	var menu Menu
	var firstMenu1,firstMenu2 FirstMenu
	firstMenu1.Key = "V1001_TODAY_MUSIC"
	firstMenu1.Name = "今日歌曲"
	firstMenu1.MenuType = "click"
	menu.FirstMenus = append(menu.FirstMenus,firstMenu1)

	firstMenu2.Name = "菜单"

	var subMenu1,subMenu2,subMenu3 SubMenu
	subMenu1.Name = "搜索"
	subMenu1.MenuType = "view"
	subMenu1.Url = "http://www.soso.com/"

	subMenu2.Name = "wxa"
	subMenu2.MenuType = "miniprogram"
	subMenu2.Url = "http://mp.weixin.qq.com"
	subMenu2.AppId = "wx508e9e50a737c414"
	subMenu2.PagePath = "pages/lunar/index"

	subMenu3.MenuType = "click"
	subMenu3.Name = "赞一下我们"
	subMenu3.Key = "V1001_GOOD"

	firstMenu2.SubMenus = append(firstMenu2.SubMenus,subMenu1)
	firstMenu2.SubMenus = append(firstMenu2.SubMenus,subMenu2)
	firstMenu2.SubMenus = append(firstMenu2.SubMenus,subMenu3)

	menu.FirstMenus = append(menu.FirstMenus,firstMenu1)
	menu.FirstMenus = append(menu.FirstMenus,firstMenu2)

	bytesData, err := json.Marshal(menu)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(bytesData))

	reader := bytes.NewReader(bytesData)

	apiUrl := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + accessToken

	requ, err := http.NewRequest("POST", apiUrl, reader)

	requ.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3278.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(requ)
	if err != nil {
		logger.Error(err)
		return
	}

	defer resp.Body.Close()

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Debug(string(recv))
}
