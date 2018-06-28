package service

import (
	"background/common/logger"
	"background/common/util"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	"fmt"
	"background/common/constant"
	"background/newmovie/service/script"
)
type OtherPlayUrl struct{
	Provider     uint32
	Channel      string
	ContentType  uint32
	Url          string
	TvType       string
	Quality      uint8
	Times        uint32
}

func GetRealUrl(provider uint32, url string)(string){

	if provider == constant.ContentProviderSystem || provider == constant.ContentProviderYouKu || provider == constant.ContentProviderMgtv{
		return url
	}

	var realUrl string
	if provider == constant.ContentProviderIqiyi{
		realUrl = script.GetIqiyiRealPlayUrl(url)
	}else if provider == constant.ContentProviderMigu{
		realUrl = script.GetMiguRealPlayUrl(constant.MediaTypeEpisode,url)
	}else if provider == constant.ContentProviderSohu{
		realUrl = script.GetSohuRealPlayUrl(url)
	}

	logger.Debug(realUrl)
	return realUrl
}


func GetStreamSourceUrl(v OtherPlayUrl,jsCode string)(string){
	if v.Provider == 0 {
		return v.Url
	}
	var err error

	type CallBackDatas struct{
		Value         []string    `gorm:"fetch_stream_ids" json:"fetch_stream_ids"`
		FetchUrlNew   string      `gorm:"fetch_url_new" json:"fetch_url_new"`
		Count         int         `gorm:"count" json:"count"`
		Urls          []string    `gorm:"urls" json:"urls"`
		Host          string      `gorm:"host" json:"host"`
		Url           string      `gorm:"url" json:"url"`
		Vkey          string      `gorm:"vkey" json:"vkey"`
		FnPre         string      `gorm:"fn_pre" json:"fn_pre"`
		FileName      string      `gorm:"file_name" json:"file_name"`
		QualityId     int         `gorm:"quality_id" json:"quality_id"`
		VideoId       string      `gorm:"video_id" json:"video_id"`
	}

	type ReqParam struct{
		/* provider : 1 cntv 2 腾讯 3 爱奇艺  4 优酷 5 芒果 6 搜狐 7 米咕 8 华数TV 9 好趣 10 电视家 11 韩剧TV 12 山寨米咕
		* quality  : 1 流畅 2 标清 3 高清 4 720P 5 1080P
		* tv_type  : 央视 卫视 地方
		*
		* content_type : 4 直播 2 点播
		* content_type为0时以下为必填:
		* channel : cntv频道名称 migu url
		* back_title:回看节目名称
		* start_time:回看
		* end_time:回看
		* content_type 为2时以下为必填:
		* url : 网页地址
		* html_data : 网页内容
		*/
		Times             uint32    `gorm:"times" json:"times"`
		Provider          uint32    `gorm:"provider" json:"provider"`
		Quality           uint8     `gorm:"quality" json:"quality"`
		TvType            string    `gorm:"tv_type" json:"tv_type"`
		ContentType       uint32    `gorm:"content_type" json:"content_type"`
		Channel           string    `gorm:"channel" json:"channel"`
		HtmlData          string    `gorm:"html_data" json:"html_data"`
		Url               string    `gorm:"url" json:"url"`
		CallBackData      CallBackDatas  `gorm:"call_back_data" json:"call_back_data"`
	}

	type Headers struct{
		UserAgent           string    `gorm:"user_agent" json:"user_agent"`
		Referer             string    `gorm:"referer" json:"referer"`
		ContentType         string    `gorm:"content_type" json:"content_type"`
	}

	type Fetch struct{
		Url           string    `gorm:"url" json:"url"`
		Method        string    `gorm:"method" json:"method"`
		Header        Headers   `gorm:"header" json:"header"`
		Body          string    `gorm:"body" json:"body"`
	}

	type ResParam struct{
		/*
		    * done:为true表示调用结束，不需要再次调用,false表示还需再次调用
		    * fetch_url_new:done为false时返回，表示下次调用是需获取网页内容的网址信息，是一个json串
		    * fetch_url_new_new中url表示请求网页地址，method表示请求方法(get/post),headers是请求头,body是请求体请求体
		    * urls:done为true时返回，表示最终的真实播放地址，是一个数组，长度为1时是完整的视频播放地址，为其他时是分段的播放地址，所有播放地址合成一个完整的视频
		    * issupportback:是否支持回看,0不支持，1支持
		*/
		Done              bool      `gorm:"done" json:"done"`
		FetchUrlNew       Fetch     `gorm:"fetch_url_new" json:"fetch_url_new"`
		Urls              []string  `gorm:"urls" json:"urls"`
		IsSupportBack     uint32    `gorm:"is_support_back" json:"is_support_back"`
		FetchUrl          string    `gorm:"fetch_url" json:"fetch_url"`
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
		reqParam.Provider = v.Provider
		reqParam.Quality = v.Quality
		reqParam.TvType = v.TvType
		reqParam.ContentType = v.ContentType
		reqParam.Channel = v.Channel
		reqParam.Url = v.Url

		if len(resParam.FetchUrlNew.Url) > 0{
			if reqParam.Provider == 8{
				var qualityCode string
				if reqParam.Quality == 4{
					qualityCode = "1000996"
				}else{
					qualityCode = "1000995"
				}
				huaShuSrc := "<message module=\"CATALOG_SERVICE\" version=\"1.0\">\n" +
					"\t<header action=\"REQUEST\" command=\"CONTENT_QUERY\" sequence=\"20121030212732_103861\" component-id=\"SYSTEM2\" component-type=\"THIRD_PARTY_SYSTEM\" />\n" +
					"\t<body>\n" +
					"\t\t<contents>\n" +
					"\t\t\t<content>\n" +
					"\t\t\t\t<code>" + reqParam.Channel + "</code>\n" +
					"\t\t\t\t<site-code>1000889</site-code>\n" +
					"\t\t\t\t<items-index>-1</items-index>\n" +
					"\t\t\t\t<folder-code>" + qualityCode + "</folder-code>\n" +
					"\t\t\t\t<format>-1</format>\n" +
					"\t\t\t</content>\n" +
					"\t\t</contents>\n" +
					"\t</body>\n" +
					"</message>"

				requ, err := http.NewRequest("POST", resParam.FetchUrlNew.Url, bytes.NewBuffer([]byte(huaShuSrc)))
				requ.Header.Add("Content-Type", "application/xml")

				resp, err := http.DefaultClient.Do(requ)
				if err != nil {
					logger.Debug("Proxy failed!")
					return ""
				}

				recv,err := ioutil.ReadAll(resp.Body)
				if err != nil{
					logger.Error(err)
				}
				reqParam.HtmlData = string(recv)

			}else{
				//resp, err := req.Get(resParam.FetchUrlNew.Url)
				//if err != nil {
				//	logger.Error(err)
				//	return ""
				//}

				requ, err := http.NewRequest("GET", resParam.FetchUrlNew.Url,nil)

				if len(resParam.FetchUrlNew.Header.UserAgent) > 0{
					requ.Header.Add("User-Agent", resParam.FetchUrlNew.Header.UserAgent)
				}
				if len(resParam.FetchUrlNew.Header.Referer) > 0{
					requ.Header.Add("Referer", resParam.FetchUrlNew.Header.Referer)
				}
				if len(resParam.FetchUrlNew.Header.ContentType) > 0{
					requ.Header.Add("Content-Type", resParam.FetchUrlNew.Header.ContentType)
				}

				resp, err := http.DefaultClient.Do(requ)
				if err != nil {
					logger.Debug("Proxy failed!")
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
		}

		b, _ := json.Marshal(reqParam)

		vm := otto.New()
		vm.Run(jsCode)
		data ,_ := vm.Call("GetRealPlayUrl",nil,string(b))

		result = data.String()
		//logger.Debug(result)
		if err = json.Unmarshal([]byte(result), &resParam); err != nil {
			logger.Error(err)
			return ""
		}

		reqParam.CallBackData = resParam.CallBackData
		first = false
	}

	var realUrl string
	if len(resParam.Urls) > 0{
		for _,v := range resParam.Urls{
			logger.Debug(v)

		}
		realUrl = resParam.Urls[0]
	}

	return realUrl
}



func GetJsCode()(string) {
	url := "http://118.89.132.150:6666/cms/v1.0/script?installation_id=1804141009224245&sign=A0F6DA2FAAC29D43F82A59173AF5F655861CB130&last_update=1524125974&app_version=1.0.14&app_key=1TqCJbPSFIWC&channel=yingyonghui&timestamp=1524192998&os_type=1"
	requ, err := http.NewRequest("GET", url, nil)

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		fmt.Println("get server jscode error!!!",err)
		return ""
	}

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return ""
	}
	//logger.Debug(string(recv))

	type Script struct {
		Data        string    `json:"data"` // s为单位
		ErrorCode   int       `json:"err_code"`
	}

	var s Script
	if err = json.Unmarshal([]byte(recv), &s); err != nil {
		logger.Error(err)
		return ""
	}

	//logger.Debug(s.Data)
	jsCode := util.DecryptAppUrl(s.Data)

	//logger.Debug(jsCode)
	return jsCode
}
