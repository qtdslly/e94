package script

import (
	"background/common/logger"
	"github.com/robertkrimen/otto"
	"encoding/json"
	"io/ioutil"
	"net/http"
)
func GetSohuRealPlayUrl(url string)(string){
	jsCode := GetSohuJsCode()

	type CallBackDatas struct{
		FetchUrl   string      `gorm:"fetch_url" json:"fetch_url"`
	}

	type ReqParam struct{
		Times             uint32    `gorm:"times" json:"times"`
		Quality           uint8     `gorm:"quality" json:"quality"`
		ContentType       uint32    `gorm:"content_type" json:"content_type"`
		HtmlData          string    `gorm:"html_data" json:"html_data"`
		Url               string    `gorm:"url" json:"url"`
		CallBackData      CallBackDatas  `gorm:"call_back_data" json:"call_back_data"`
	}

	type Headers struct{
		UserAgent           string    `gorm:"user_agent" json:"user_agent"`
		Referer             string    `gorm:"referer" json:"referer"`
		ContentType         string    `gorm:"content_type" json:"content_type"`
		UserId              string    `gorm:"userId" json:"userId"`
		UserToken           string    `gorm:"userToken" json:"userToken"`
		SdkceId             string    `gorm:"SDKCEId" json:"SDKCEId"`
		ClientId            string    `gorm:"clientId" json:"clientId"`
	}

	type Fetch struct{
		Url           string    `gorm:"url" json:"url"`
		Method        string    `gorm:"method" json:"method"`
		Header        Headers   `gorm:"header" json:"header"`
		Body          string    `gorm:"body" json:"body"`
	}

	type ResParam struct{
		Done              bool      `gorm:"done" json:"done"`
		FetchUrl          Fetch     `gorm:"fetch_url" json:"fetch_url"`
		Urls              []string  `gorm:"urls" json:"urls"`
		CallBackData      CallBackDatas  `gorm:"call_back_data" json:"call_back_data"`
	}

	first := true
	var reqParam ReqParam
	var resParam ResParam
	var times uint32 = 0

	var result string
	for first || !resParam.Done {
		times = times + 1
		reqParam.Times = times
		reqParam.Quality = 3
		reqParam.ContentType = 2
		reqParam.Url = url

		if len(resParam.FetchUrl.Url) > 0{
			requ, err := http.NewRequest("GET", resParam.FetchUrl.Url,nil)

			if len(resParam.FetchUrl.Header.UserAgent) > 0{
				requ.Header.Add("User-Agent", resParam.FetchUrl.Header.UserAgent)
			}
			if len(resParam.FetchUrl.Header.Referer) > 0{
				requ.Header.Add("Referer", resParam.FetchUrl.Header.Referer)
			}
			if len(resParam.FetchUrl.Header.ContentType) > 0{
				requ.Header.Add("Content-Type", resParam.FetchUrl.Header.ContentType)
			}
			if len(resParam.FetchUrl.Header.UserId) > 0{
				requ.Header.Add("userId", resParam.FetchUrl.Header.UserId)
			}
			if len(resParam.FetchUrl.Header.ContentType) > 0{
				requ.Header.Add("userToken", resParam.FetchUrl.Header.UserToken)
			}
			if len(resParam.FetchUrl.Header.ContentType) > 0{
				requ.Header.Add("SDKCEId", resParam.FetchUrl.Header.SdkceId)
			}
			if len(resParam.FetchUrl.Header.ClientId) > 0{
				requ.Header.Add("clientId", resParam.FetchUrl.Header.ClientId)
			}

			resp, err := http.DefaultClient.Do(requ)
			if err != nil {
				logger.Debug("err")
				return ""
			}

			recv,err := ioutil.ReadAll(resp.Body)
			if err != nil{
				logger.Error(err)
			}
			//logger.Debug(string(recv))
			reqParam.HtmlData = string(recv)

			//logger.Debug("html_data:",reqParam.HtmlData)

		}

		b, _ := json.Marshal(reqParam)
		//logger.Debug(string(b))

		vm := otto.New()
		vm.Run(jsCode)
		data ,err := vm.Call("GetSohuRealPlayUrl",nil,string(b))
		if err != nil{
			logger.Error(err)
			return ""
		}
		result = data.String()
		//logger.Debug(result)

		if err := json.Unmarshal([]byte(result), &resParam); err != nil {
			logger.Error(err)
			return ""
		}

		reqParam.CallBackData = resParam.CallBackData
		first = false
	}

	if len(resParam.Urls) > 0{
		return resParam.Urls[0]
	}
	return ""
}

func GetSohuJsCode()(string){
	jsCode := "\n" +
		"/**\n" +
		" * @return {string}\n" +
		" */\n" +
		"function GetSohuRealPlayUrl(request) {\n" +
		"    var json_data;\n" +
		"    var fetch_url;\n" +
		"    var result;\n" +
		"\n" +
		"    json_data = JSON.parse(request);\n" +
		"    var referer = json_data.url;\n" +
		"\n" +
		"    var agent = \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36\";\n" +
		"    fetch_url = {method: \"get\", header: {user_agent: agent,referer: referer}, body: \"\", url: \"\"};\n" +
		"    result = {done: false, fetch_url: fetch_url, urls: [], call_back_data: {}};\n" +
		"    var quality;  //quality  : 1 流畅 2 标清 3 高清 4 720P 5 1080P\n" +
		"    if(typeof(json_data.quality) === \"undefined\" ){\n" +
		"        quality = 3;\n" +
		"    }else{\n" +
		"        quality = json_data.quality;\n" +
		"    }\n" +
		"\n" +
		"    if(json_data.times === 1){\n" +
		"        fetch_url.url = json_data.url;\n" +
		"        result.done = false;\n" +
		"    }else if(json_data.times === 2){\n" +
		"        try{\n" +
		"            re = \"var vid=\\\"(.+?)\\\";\";\n" +
		"            video_id = json_data.html_data.match(re)[1];\n" +
		"            fetch_url.url = \"https://m.tv.sohu.com/phone_playinfo?callback=jsonpx1517539948922_56_11&vid=\" + video_id + \"&site=1&appid=tv&api_key=f351515304020cad28c92f70f002261c&plat=17&sver=1.0&partner=1&uid=1711141102362091&muid=1517539772399202&_c=1&pt=5&qd=680&src=11060001&ssl=1&_=1517539948922\";\n" +
		"            result.done = false;\n" +
		"        }catch (e){\n" +
		"            result.error = e.toString();\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            result.urls.splice(0,result.urls.length);\n" +
		"            return JSON.stringify(result);\n" +
		"        }\n" +
		"    }else if(json_data.times === 3){\n" +
		"        try{\n" +
		"            //https://m.tv.sohu.com/phone_playinfo?callback=jsonpx1517539948922_56_11&vid=3055369&site=1&appid=tv&api_key=f351515304020cad28c92f70f002261c&plat=17&sver=1.0&partner=1&uid=1711141102362091&muid=1517539772399202&_c=1&pt=5&qd=680&src=11060001&ssl=1&_=1517539948922\n" +
		"            data = json_data.html_data;\n" +
		"            data = data.replace(\"jsonpx1517539948922_56_11(\",\"\");\n" +
		"            data = data.replace(\"})\",\"}\");\n" +
		"            response = JSON.parse(data);\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            if(quality === 2){ // nor 640 * 360\n" +
		"                result.urls.push(response.data.urls.m3u8.nor[0]);\n" +
		"            }else if(quality === 4){ //sup 1280 * 720\n" +
		"                result.urls.push(response.data.urls.m3u8.sup[0]);\n" +
		"            }else if(quality === 3){ // high 852 * 480\n" +
		"                result.urls.push(response.data.urls.m3u8.hig[0]);\n" +
		"            }else{\n" +
		"                fetch_url.url = \"\";\n" +
		"                result.done = true;\n" +
		"                result.urls.splice(0,result.urls.length);\n" +
		"                return JSON.stringify(result);\n" +
		"            }\n" +
		"\n" +
		"            if(result.urls[0] === \"\" || typeof(result.urls[0]) === \"undefined\"){\n" +
		"                result.urls.splice(0,result.urls.length);\n" +
		"                if(response.data.urls.m3u8.sup[0] !== \"\" && typeof(response.data.urls.m3u8.sup[0]) !== \"undefined\"){\n" +
		"                    result.urls.push(response.data.urls.m3u8.sup[0]);\n" +
		"                }else if(response.data.urls.m3u8.hig[0] !== \"\" && typeof(response.data.urls.m3u8.hig[0]) !== \"undefined\"){\n" +
		"                    result.urls.push(response.data.urls.m3u8.hig[0]);\n" +
		"                }else if(response.data.urls.m3u8.nor[0] !== \"\" && typeof(response.data.urls.m3u8.nor[0]) !== \"undefined\" ){\n" +
		"                    result.urls.push(response.data.urls.m3u8.nor[0]);\n" +
		"                }\n" +
		"            }\n" +
		"        }catch (e){\n" +
		"            result.error = e.toString();\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            result.urls.splice(0,result.urls.length);\n" +
		"            return JSON.stringify(result);\n" +
		"        }\n" +
		"    }else{\n" +
		"        fetch_url.url = \"\";\n" +
		"        result.error = e.toString();\n" +
		"        result.done = true;\n" +
		"        result.urls.splice(0,result.urls.length);\n" +
		"        return JSON.stringify(result);\n" +
		"    }\n" +
		"    return JSON.stringify(result);\n" +
		"}\n" +
		"\n"
	return jsCode
}