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

func GetRealUrl(provider uint32, url string,jsCode string)(string){

	if provider == constant.ContentProviderSystem{
		return url
	}

	var playUrl OtherPlayUrl
	playUrl.Url = url
	playUrl.Channel = ""
	playUrl.Quality = 3
	playUrl.TvType = ""
	playUrl.Times = 0
	playUrl.ContentType = 2
	if provider == constant.ContentProviderYouKu{
		playUrl.Provider = 4
	}else if provider == constant.ContentProviderIqiyi{
		playUrl.Provider = 3
	}else if provider == constant.ContentProviderMgtv{
		playUrl.Provider = 5
	}
	realUrl := GetStreamSourceUrl(playUrl,jsCode)

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


//
//func GetJsCode()(string){
//
//	jsCode := "!function(a){\"use strict\";function b(a,b){var c=(65535&a)+(65535&b),d=(a>>16)+(b>>16)+(c>>16);return d<<16|65535&c}function c(a,b){return a<<b|a>>>32-b}function d(a,d,e,f,g,h){return b(c(b(b(d,a),b(f,h)),g),e)}function e(a,b,c,e,f,g,h){return d(b&c|~b&e,a,b,f,g,h)}function f(a,b,c,e,f,g,h){return d(b&e|c&~e,a,b,f,g,h)}function g(a,b,c,e,f,g,h){return d(b^c^e,a,b,f,g,h)}function h(a,b,c,e,f,g,h){return d(c^(b|~e),a,b,f,g,h)}function i(a,c){a[c>>5]|=128<<c%32,a[(c+64>>>9<<4)+14]=c;var d,i,j,k,l,m=1732584193,n=-271733879,o=-1732584194,p=271733878;for(d=0;d<a.length;d+=16)i=m,j=n,k=o,l=p,m=e(m,n,o,p,a[d],7,-680876936),p=e(p,m,n,o,a[d+1],12,-389564586),o=e(o,p,m,n,a[d+2],17,606105819),n=e(n,o,p,m,a[d+3],22,-1044525330),m=e(m,n,o,p,a[d+4],7,-176418897),p=e(p,m,n,o,a[d+5],12,1200080426),o=e(o,p,m,n,a[d+6],17,-1473231341),n=e(n,o,p,m,a[d+7],22,-45705983),m=e(m,n,o,p,a[d+8],7,1770035416),p=e(p,m,n,o,a[d+9],12,-1958414417),o=e(o,p,m,n,a[d+10],17,-42063),n=e(n,o,p,m,a[d+11],22,-1990404162),m=e(m,n,o,p,a[d+12],7,1804603682),p=e(p,m,n,o,a[d+13],12,-40341101),o=e(o,p,m,n,a[d+14],17,-1502002290),n=e(n,o,p,m,a[d+15],22,1236535329),m=f(m,n,o,p,a[d+1],5,-165796510),p=f(p,m,n,o,a[d+6],9,-1069501632),o=f(o,p,m,n,a[d+11],14,643717713),n=f(n,o,p,m,a[d],20,-373897302),m=f(m,n,o,p,a[d+5],5,-701558691),p=f(p,m,n,o,a[d+10],9,38016083),o=f(o,p,m,n,a[d+15],14,-660478335),n=f(n,o,p,m,a[d+4],20,-405537848),m=f(m,n,o,p,a[d+9],5,568446438),p=f(p,m,n,o,a[d+14],9,-1019803690),o=f(o,p,m,n,a[d+3],14,-187363961),n=f(n,o,p,m,a[d+8],20,1163531501),m=f(m,n,o,p,a[d+13],5,-1444681467),p=f(p,m,n,o,a[d+2],9,-51403784),o=f(o,p,m,n,a[d+7],14,1735328473),n=f(n,o,p,m,a[d+12],20,-1926607734),m=g(m,n,o,p,a[d+5],4,-378558),p=g(p,m,n,o,a[d+8],11,-2022574463),o=g(o,p,m,n,a[d+11],16,1839030562),n=g(n,o,p,m,a[d+14],23,-35309556),m=g(m,n,o,p,a[d+1],4,-1530992060),p=g(p,m,n,o,a[d+4],11,1272893353),o=g(o,p,m,n,a[d+7],16,-155497632),n=g(n,o,p,m,a[d+10],23,-1094730640),m=g(m,n,o,p,a[d+13],4,681279174),p=g(p,m,n,o,a[d],11,-358537222),o=g(o,p,m,n,a[d+3],16,-722521979),n=g(n,o,p,m,a[d+6],23,76029189),m=g(m,n,o,p,a[d+9],4,-640364487),p=g(p,m,n,o,a[d+12],11,-421815835),o=g(o,p,m,n,a[d+15],16,530742520),n=g(n,o,p,m,a[d+2],23,-995338651),m=h(m,n,o,p,a[d],6,-198630844),p=h(p,m,n,o,a[d+7],10,1126891415),o=h(o,p,m,n,a[d+14],15,-1416354905),n=h(n,o,p,m,a[d+5],21,-57434055),m=h(m,n,o,p,a[d+12],6,1700485571),p=h(p,m,n,o,a[d+3],10,-1894986606),o=h(o,p,m,n,a[d+10],15,-1051523),n=h(n,o,p,m,a[d+1],21,-2054922799),m=h(m,n,o,p,a[d+8],6,1873313359),p=h(p,m,n,o,a[d+15],10,-30611744),o=h(o,p,m,n,a[d+6],15,-1560198380),n=h(n,o,p,m,a[d+13],21,1309151649),m=h(m,n,o,p,a[d+4],6,-145523070),p=h(p,m,n,o,a[d+11],10,-1120210379),o=h(o,p,m,n,a[d+2],15,718787259),n=h(n,o,p,m,a[d+9],21,-343485551),m=b(m,i),n=b(n,j),o=b(o,k),p=b(p,l);return[m,n,o,p]}function j(a){var b,c=\"\";for(b=0;b<32*a.length;b+=8)c+=String.fromCharCode(a[b>>5]>>>b%32&255);return c}function k(a){var b,c=[];for(c[(a.length>>2)-1]=void 0,b=0;b<c.length;b+=1)c[b]=0;for(b=0;b<8*a.length;b+=8)c[b>>5]|=(255&a.charCodeAt(b/8))<<b%32;return c}function l(a){return j(i(k(a),8*a.length))}function m(a,b){var c,d,e=k(a),f=[],g=[];for(f[15]=g[15]=void 0,e.length>16&&(e=i(e,8*a.length)),c=0;16>c;c+=1)f[c]=909522486^e[c],g[c]=1549556828^e[c];return d=i(f.concat(k(b)),512+8*b.length),j(i(g.concat(d),640))}function n(a){var b,c,d=\"0123456789abcdef\",e=\"\";for(c=0;c<a.length;c+=1)b=a.charCodeAt(c),e+=d.charAt(b>>>4&15)+d.charAt(15&b);return e}function o(a){return unescape(encodeURIComponent(a))}function p(a){return l(o(a))}function q(a){return n(p(a))}function r(a,b){return m(o(a),o(b))}function s(a,b){return n(r(a,b))}function t(a,b,c){return b?c?r(b,a):s(b,a):c?p(a):q(a)}\"function\"==typeof define&&define.amd?define(function(){return t}):a.md5=t}(this);\n" +
//		"\n" +
//		"var _0x9c59=[\"\x33\x64\x70\x6F\x6C\",\"\x70\x61\x72\x73\x65\",\"\x30\",\"\x72\x61\x6E\x64\x6F\x6D\",\"\x66\x6C\x6F\x6F\x72\",\"\x2F\x63\x6E\x74\x76\x6C\x69\x76\x65\x2F\",\n" +
//		"    \"\x6D\x64\x2E\x6D\x33\x75\x38\",\"\x2D\",\"\x32\x65\x66\x61\x65\x39\x72\x33\x32\x72\x6B\x76\x39\x38\x33\x32\x6A\x73\x32\x31\",\"\x6F\",\"\x6E\x33\x64\",\n" +
//		"    \"\x72\x65\x70\x6C\x61\x63\x65\",\"\x68\x74\x74\x70\x3A\x2F\x2F\x68\x6C\x73\x32\x2E\x63\x6E\x74\x76\x2E\x6D\x79\x61\x6C\x69\x63\x64\x6E\x2E\x63\x6F\x6D\",\n" +
//		"    \"\x3F\x61\x75\x74\x68\x5F\x6B\x65\x79\x3D\",\"\x31\",\"\x73\x72\x63\",\"\x5F\x68\x74\x6D\x6C\x35\x50\x6C\x61\x79\x65\x72\",\"\x68\x75\x61\x77\x65\x69\",\n" +
//		"    \"\x69\x6E\x64\x65\x78\x4F\x66\",\"\x74\x6F\x4C\x6F\x77\x65\x72\x43\x61\x73\x65\",\"\x75\x73\x65\x72\x41\x67\x65\x6E\x74\",\"\x6C\x6F\x61\x64\",\n" +
//		"    \"\x64\x69\x61\x6E\x70\x69\x61\x6E\x2E\x6D\x70\x34\",\"\x75\x6E\x64\x65\x66\x69\x6E\x65\x64\",\"\x6C\x65\x6E\x67\x74\x68\",\"\x70\x6F\x73\x74\x65\x72\",\n" +
//		"    \"\x4C\x49\x56\x45\x2D\x48\x44\x53\x2D\x43\x44\x4E\x2D\x41\x4C\x49\"];\n" +
//		"\n" +
//		"function setFlvHtml5AliNewUrl(videoName) {\n" +
//		"    var _0x54c7x3 = _0x9c59[0];\n" +
//		"    var _0x54c7x4 = Date[_0x9c59[1]](new Date()) / 1000;\n" +
//		"    var _0x54c7x5 = _0x9c59[2];\n" +
//		"    var _0x54c7x6 = Math[_0x9c59[4]](Math[_0x9c59[3]]() * 1000);\n" +
//		"    var _0x54c7x7 = _0x9c59[5] + videoName + _0x9c59[6];\n" +
//		"    var _0xa4b0 = [\"\x6A\x70\x39\x76\"];\n" +
//		"    var html5Aauth = _0xa4b0[0];\n" +
//		"    var _0x54c7x8 = md5(_0x54c7x7 + _0x9c59[7] + _0x54c7x4 + _0x9c59[7] + _0x54c7x6 + _0x9c59[7] + _0x54c7x5 + _0x9c59[7] + _0x9c59[8][12] + html5Aauth + _0x54c7x3[_0x9c59[11]](_0x9c59[9], _0x9c59[10]));\n" +
//		"    var _0x54c7x9 = _0x54c7x4 + _0x9c59[7] + _0x54c7x6 + _0x9c59[7] + _0x54c7x5 + _0x9c59[7] + _0x54c7x8;\n" +
//		"    return _0x9c59[12] + _0x54c7x7 + _0x9c59[13] + _0x54c7x9;\n" +
//		"\n" +
//		"}\n" +
//		"\n" +
//		"\n" +
//		"\n" +
//		"/*\n" +
//		"CryptoJS v3.1.2\n" +
//		"code.google.com/p/crypto-js\n" +
//		"(c) 2009-2013 by Jeff Mott. All rights reserved.\n" +
//		"code.google.com/p/crypto-js/wiki/License\n" +
//		"*/\n" +
//		"var CryptoJS=CryptoJS||function(u,l){var d={},n=d.lib={},p=function(){},s=n.Base={extend:function(a){p.prototype=this;var c=new p;a&&c.mixIn(a);c.hasOwnProperty(\"init\")||(c.init=function(){c.$super.init.apply(this,arguments)});c.init.prototype=c;c.$super=this;return c},create:function(){var a=this.extend();a.init.apply(a,arguments);return a},init:function(){},mixIn:function(a){for(var c in a)a.hasOwnProperty(c)&&(this[c]=a[c]);a.hasOwnProperty(\"toString\")&&(this.toString=a.toString)},clone:function(){return this.init.prototype.extend(this)}},\n" +
//		"q=n.WordArray=s.extend({init:function(a,c){a=this.words=a||[];this.sigBytes=c!=l?c:4*a.length},toString:function(a){return(a||v).stringify(this)},concat:function(a){var c=this.words,m=a.words,f=this.sigBytes;a=a.sigBytes;this.clamp();if(f%4)for(var t=0;t<a;t++)c[f+t>>>2]|=(m[t>>>2]>>>24-8*(t%4)&255)<<24-8*((f+t)%4);else if(65535<m.length)for(t=0;t<a;t+=4)c[f+t>>>2]=m[t>>>2];else c.push.apply(c,m);this.sigBytes+=a;return this},clamp:function(){var a=this.words,c=this.sigBytes;a[c>>>2]&=4294967295<<\n" +
//		"32-8*(c%4);a.length=u.ceil(c/4)},clone:function(){var a=s.clone.call(this);a.words=this.words.slice(0);return a},random:function(a){for(var c=[],m=0;m<a;m+=4)c.push(4294967296*u.random()|0);return new q.init(c,a)}}),w=d.enc={},v=w.Hex={stringify:function(a){var c=a.words;a=a.sigBytes;for(var m=[],f=0;f<a;f++){var t=c[f>>>2]>>>24-8*(f%4)&255;m.push((t>>>4).toString(16));m.push((t&15).toString(16))}return m.join(\"\")},parse:function(a){for(var c=a.length,m=[],f=0;f<c;f+=2)m[f>>>3]|=parseInt(a.substr(f,\n" +
//		"2),16)<<24-4*(f%8);return new q.init(m,c/2)}},b=w.Latin1={stringify:function(a){var c=a.words;a=a.sigBytes;for(var m=[],f=0;f<a;f++)m.push(String.fromCharCode(c[f>>>2]>>>24-8*(f%4)&255));return m.join(\"\")},parse:function(a){for(var c=a.length,m=[],f=0;f<c;f++)m[f>>>2]|=(a.charCodeAt(f)&255)<<24-8*(f%4);return new q.init(m,c)}},x=w.Utf8={stringify:function(a){try{return decodeURIComponent(escape(b.stringify(a)))}catch(c){throw Error(\"Malformed UTF-8 data\");}},parse:function(a){return b.parse(unescape(encodeURIComponent(a)))}},\n" +
//		"r=n.BufferedBlockAlgorithm=s.extend({reset:function(){this._data=new q.init;this._nDataBytes=0},_append:function(a){\"string\"==typeof a&&(a=x.parse(a));this._data.concat(a);this._nDataBytes+=a.sigBytes},_process:function(a){var c=this._data,m=c.words,f=c.sigBytes,t=this.blockSize,b=f/(4*t),b=a?u.ceil(b):u.max((b|0)-this._minBufferSize,0);a=b*t;f=u.min(4*a,f);if(a){for(var e=0;e<a;e+=t)this._doProcessBlock(m,e);e=m.splice(0,a);c.sigBytes-=f}return new q.init(e,f)},clone:function(){var a=s.clone.call(this);\n" +
//		"a._data=this._data.clone();return a},_minBufferSize:0});n.Hasher=r.extend({cfg:s.extend(),init:function(a){this.cfg=this.cfg.extend(a);this.reset()},reset:function(){r.reset.call(this);this._doReset()},update:function(a){this._append(a);this._process();return this},finalize:function(a){a&&this._append(a);return this._doFinalize()},blockSize:16,_createHelper:function(a){return function(c,m){return(new a.init(m)).finalize(c)}},_createHmacHelper:function(a){return function(c,m){return(new e.HMAC.init(a,\n" +
//		"m)).finalize(c)}}});var e=d.algo={};return d}(Math);\n" +
//		"(function(){var u=CryptoJS,l=u.lib.WordArray;u.enc.Base64={stringify:function(d){var n=d.words,l=d.sigBytes,s=this._map;d.clamp();d=[];for(var q=0;q<l;q+=3)for(var w=(n[q>>>2]>>>24-8*(q%4)&255)<<16|(n[q+1>>>2]>>>24-8*((q+1)%4)&255)<<8|n[q+2>>>2]>>>24-8*((q+2)%4)&255,v=0;4>v&&q+0.75*v<l;v++)d.push(s.charAt(w>>>6*(3-v)&63));if(n=s.charAt(64))for(;d.length%4;)d.push(n);return d.join(\"\")},parse:function(d){var n=d.length,p=this._map,s=p.charAt(64);s&&(s=d.indexOf(s),-1!=s&&(n=s));for(var s=[],q=0,w=0;w<\n" +
//		"n;w++)if(w%4){var v=p.indexOf(d.charAt(w-1))<<2*(w%4),b=p.indexOf(d.charAt(w))>>>6-2*(w%4);s[q>>>2]|=(v|b)<<24-8*(q%4);q++}return l.create(s,q)},_map:\"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=\"}})();\n" +
//		"(function(u){function l(b,e,a,c,m,f,t){b=b+(e&a|~e&c)+m+t;return(b<<f|b>>>32-f)+e}function d(b,e,a,c,m,f,t){b=b+(e&c|a&~c)+m+t;return(b<<f|b>>>32-f)+e}function n(b,e,a,c,m,f,t){b=b+(e^a^c)+m+t;return(b<<f|b>>>32-f)+e}function p(b,e,a,c,m,f,t){b=b+(a^(e|~c))+m+t;return(b<<f|b>>>32-f)+e}for(var s=CryptoJS,q=s.lib,w=q.WordArray,v=q.Hasher,q=s.algo,b=[],x=0;64>x;x++)b[x]=4294967296*u.abs(u.sin(x+1))|0;q=q.MD5=v.extend({_doReset:function(){this._hash=new w.init([1732584193,4023233417,2562383102,271733878])},\n" +
//		"_doProcessBlock:function(r,e){for(var a=0;16>a;a++){var c=e+a,m=r[c];r[c]=(m<<8|m>>>24)&16711935|(m<<24|m>>>8)&4278255360}var a=this._hash.words,c=r[e+0],m=r[e+1],f=r[e+2],t=r[e+3],y=r[e+4],q=r[e+5],s=r[e+6],w=r[e+7],v=r[e+8],u=r[e+9],x=r[e+10],z=r[e+11],A=r[e+12],B=r[e+13],C=r[e+14],D=r[e+15],g=a[0],h=a[1],j=a[2],k=a[3],g=l(g,h,j,k,c,7,b[0]),k=l(k,g,h,j,m,12,b[1]),j=l(j,k,g,h,f,17,b[2]),h=l(h,j,k,g,t,22,b[3]),g=l(g,h,j,k,y,7,b[4]),k=l(k,g,h,j,q,12,b[5]),j=l(j,k,g,h,s,17,b[6]),h=l(h,j,k,g,w,22,b[7]),\n" +
//		"g=l(g,h,j,k,v,7,b[8]),k=l(k,g,h,j,u,12,b[9]),j=l(j,k,g,h,x,17,b[10]),h=l(h,j,k,g,z,22,b[11]),g=l(g,h,j,k,A,7,b[12]),k=l(k,g,h,j,B,12,b[13]),j=l(j,k,g,h,C,17,b[14]),h=l(h,j,k,g,D,22,b[15]),g=d(g,h,j,k,m,5,b[16]),k=d(k,g,h,j,s,9,b[17]),j=d(j,k,g,h,z,14,b[18]),h=d(h,j,k,g,c,20,b[19]),g=d(g,h,j,k,q,5,b[20]),k=d(k,g,h,j,x,9,b[21]),j=d(j,k,g,h,D,14,b[22]),h=d(h,j,k,g,y,20,b[23]),g=d(g,h,j,k,u,5,b[24]),k=d(k,g,h,j,C,9,b[25]),j=d(j,k,g,h,t,14,b[26]),h=d(h,j,k,g,v,20,b[27]),g=d(g,h,j,k,B,5,b[28]),k=d(k,g,\n" +
//		"h,j,f,9,b[29]),j=d(j,k,g,h,w,14,b[30]),h=d(h,j,k,g,A,20,b[31]),g=n(g,h,j,k,q,4,b[32]),k=n(k,g,h,j,v,11,b[33]),j=n(j,k,g,h,z,16,b[34]),h=n(h,j,k,g,C,23,b[35]),g=n(g,h,j,k,m,4,b[36]),k=n(k,g,h,j,y,11,b[37]),j=n(j,k,g,h,w,16,b[38]),h=n(h,j,k,g,x,23,b[39]),g=n(g,h,j,k,B,4,b[40]),k=n(k,g,h,j,c,11,b[41]),j=n(j,k,g,h,t,16,b[42]),h=n(h,j,k,g,s,23,b[43]),g=n(g,h,j,k,u,4,b[44]),k=n(k,g,h,j,A,11,b[45]),j=n(j,k,g,h,D,16,b[46]),h=n(h,j,k,g,f,23,b[47]),g=p(g,h,j,k,c,6,b[48]),k=p(k,g,h,j,w,10,b[49]),j=p(j,k,g,h,\n" +
//		"C,15,b[50]),h=p(h,j,k,g,q,21,b[51]),g=p(g,h,j,k,A,6,b[52]),k=p(k,g,h,j,t,10,b[53]),j=p(j,k,g,h,x,15,b[54]),h=p(h,j,k,g,m,21,b[55]),g=p(g,h,j,k,v,6,b[56]),k=p(k,g,h,j,D,10,b[57]),j=p(j,k,g,h,s,15,b[58]),h=p(h,j,k,g,B,21,b[59]),g=p(g,h,j,k,y,6,b[60]),k=p(k,g,h,j,z,10,b[61]),j=p(j,k,g,h,f,15,b[62]),h=p(h,j,k,g,u,21,b[63]);a[0]=a[0]+g|0;a[1]=a[1]+h|0;a[2]=a[2]+j|0;a[3]=a[3]+k|0},_doFinalize:function(){var b=this._data,e=b.words,a=8*this._nDataBytes,c=8*b.sigBytes;e[c>>>5]|=128<<24-c%32;var m=u.floor(a/\n" +
//		"4294967296);e[(c+64>>>9<<4)+15]=(m<<8|m>>>24)&16711935|(m<<24|m>>>8)&4278255360;e[(c+64>>>9<<4)+14]=(a<<8|a>>>24)&16711935|(a<<24|a>>>8)&4278255360;b.sigBytes=4*(e.length+1);this._process();b=this._hash;e=b.words;for(a=0;4>a;a++)c=e[a],e[a]=(c<<8|c>>>24)&16711935|(c<<24|c>>>8)&4278255360;return b},clone:function(){var b=v.clone.call(this);b._hash=this._hash.clone();return b}});s.MD5=v._createHelper(q);s.HmacMD5=v._createHmacHelper(q)})(Math);\n" +
//		"(function(){var u=CryptoJS,l=u.lib,d=l.Base,n=l.WordArray,l=u.algo,p=l.EvpKDF=d.extend({cfg:d.extend({keySize:4,hasher:l.MD5,iterations:1}),init:function(d){this.cfg=this.cfg.extend(d)},compute:function(d,l){for(var p=this.cfg,v=p.hasher.create(),b=n.create(),u=b.words,r=p.keySize,p=p.iterations;u.length<r;){e&&v.update(e);var e=v.update(d).finalize(l);v.reset();for(var a=1;a<p;a++)e=v.finalize(e),v.reset();b.concat(e)}b.sigBytes=4*r;return b}});u.EvpKDF=function(d,l,n){return p.create(n).compute(d,\n" +
//		"l)}})();\n" +
//		"CryptoJS.lib.Cipher||function(u){var l=CryptoJS,d=l.lib,n=d.Base,p=d.WordArray,s=d.BufferedBlockAlgorithm,q=l.enc.Base64,w=l.algo.EvpKDF,v=d.Cipher=s.extend({cfg:n.extend(),createEncryptor:function(m,a){return this.create(this._ENC_XFORM_MODE,m,a)},createDecryptor:function(m,a){return this.create(this._DEC_XFORM_MODE,m,a)},init:function(m,a,b){this.cfg=this.cfg.extend(b);this._xformMode=m;this._key=a;this.reset()},reset:function(){s.reset.call(this);this._doReset()},process:function(a){this._append(a);return this._process()},\n" +
//		"finalize:function(a){a&&this._append(a);return this._doFinalize()},keySize:4,ivSize:4,_ENC_XFORM_MODE:1,_DEC_XFORM_MODE:2,_createHelper:function(m){return{encrypt:function(f,b,e){return(\"string\"==typeof b?c:a).encrypt(m,f,b,e)},decrypt:function(f,b,e){return(\"string\"==typeof b?c:a).decrypt(m,f,b,e)}}}});d.StreamCipher=v.extend({_doFinalize:function(){return this._process(!0)},blockSize:1});var b=l.mode={},x=function(a,f,b){var c=this._iv;c?this._iv=u:c=this._prevBlock;for(var e=0;e<b;e++)a[f+e]^=\n" +
//		"c[e]},r=(d.BlockCipherMode=n.extend({createEncryptor:function(a,f){return this.Encryptor.create(a,f)},createDecryptor:function(a,f){return this.Decryptor.create(a,f)},init:function(a,f){this._cipher=a;this._iv=f}})).extend();r.Encryptor=r.extend({processBlock:function(a,f){var b=this._cipher,c=b.blockSize;x.call(this,a,f,c);b.encryptBlock(a,f);this._prevBlock=a.slice(f,f+c)}});r.Decryptor=r.extend({processBlock:function(a,b){var c=this._cipher,e=c.blockSize,d=a.slice(b,b+e);c.decryptBlock(a,b);x.call(this,\n" +
//		"a,b,e);this._prevBlock=d}});b=b.CBC=r;r=(l.pad={}).Pkcs7={pad:function(a,b){for(var c=4*b,c=c-a.sigBytes%c,e=c<<24|c<<16|c<<8|c,d=[],l=0;l<c;l+=4)d.push(e);c=p.create(d,c);a.concat(c)},unpad:function(a){a.sigBytes-=a.words[a.sigBytes-1>>>2]&255}};d.BlockCipher=v.extend({cfg:v.cfg.extend({mode:b,padding:r}),reset:function(){v.reset.call(this);var a=this.cfg,c=a.iv,a=a.mode;if(this._xformMode==this._ENC_XFORM_MODE)var b=a.createEncryptor;else b=a.createDecryptor,this._minBufferSize=1;this._mode=b.call(a,\n" +
//		"this,c&&c.words)},_doProcessBlock:function(a,c){this._mode.processBlock(a,c)},_doFinalize:function(){var a=this.cfg.padding;if(this._xformMode==this._ENC_XFORM_MODE){a.pad(this._data,this.blockSize);var c=this._process(!0)}else c=this._process(!0),a.unpad(c);return c},blockSize:4});var e=d.CipherParams=n.extend({init:function(a){this.mixIn(a)},toString:function(a){return(a||this.formatter).stringify(this)}}),b=(l.format={}).OpenSSL={stringify:function(a){var c=a.ciphertext;a=a.salt;return(a?p.create([1398893684,\n" +
//		"1701076831]).concat(a).concat(c):c).toString(q)},parse:function(a){a=q.parse(a);var c=a.words;if(1398893684==c[0]&&1701076831==c[1]){var b=p.create(c.slice(2,4));c.splice(0,4);a.sigBytes-=16}return e.create({ciphertext:a,salt:b})}},a=d.SerializableCipher=n.extend({cfg:n.extend({format:b}),encrypt:function(a,c,b,d){d=this.cfg.extend(d);var l=a.createEncryptor(b,d);c=l.finalize(c);l=l.cfg;return e.create({ciphertext:c,key:b,iv:l.iv,algorithm:a,mode:l.mode,padding:l.padding,blockSize:a.blockSize,formatter:d.format})},\n" +
//		"decrypt:function(a,c,b,e){e=this.cfg.extend(e);c=this._parse(c,e.format);return a.createDecryptor(b,e).finalize(c.ciphertext)},_parse:function(a,c){return\"string\"==typeof a?c.parse(a,this):a}}),l=(l.kdf={}).OpenSSL={execute:function(a,c,b,d){d||(d=p.random(8));a=w.create({keySize:c+b}).compute(a,d);b=p.create(a.words.slice(c),4*b);a.sigBytes=4*c;return e.create({key:a,iv:b,salt:d})}},c=d.PasswordBasedCipher=a.extend({cfg:a.cfg.extend({kdf:l}),encrypt:function(c,b,e,d){d=this.cfg.extend(d);e=d.kdf.execute(e,\n" +
//		"c.keySize,c.ivSize);d.iv=e.iv;c=a.encrypt.call(this,c,b,e.key,d);c.mixIn(e);return c},decrypt:function(c,b,e,d){d=this.cfg.extend(d);b=this._parse(b,d.format);e=d.kdf.execute(e,c.keySize,c.ivSize,b.salt);d.iv=e.iv;return a.decrypt.call(this,c,b,e.key,d)}})}();\n" +
//		"(function(){function u(b,a){var c=(this._lBlock>>>b^this._rBlock)&a;this._rBlock^=c;this._lBlock^=c<<b}function l(b,a){var c=(this._rBlock>>>b^this._lBlock)&a;this._lBlock^=c;this._rBlock^=c<<b}var d=CryptoJS,n=d.lib,p=n.WordArray,n=n.BlockCipher,s=d.algo,q=[57,49,41,33,25,17,9,1,58,50,42,34,26,18,10,2,59,51,43,35,27,19,11,3,60,52,44,36,63,55,47,39,31,23,15,7,62,54,46,38,30,22,14,6,61,53,45,37,29,21,13,5,28,20,12,4],w=[14,17,11,24,1,5,3,28,15,6,21,10,23,19,12,4,26,8,16,7,27,20,13,2,41,52,31,37,47,\n" +
//		"55,30,40,51,45,33,48,44,49,39,56,34,53,46,42,50,36,29,32],v=[1,2,4,6,8,10,12,14,15,17,19,21,23,25,27,28],b=[{\"0\":8421888,268435456:32768,536870912:8421378,805306368:2,1073741824:512,1342177280:8421890,1610612736:8389122,1879048192:8388608,2147483648:514,2415919104:8389120,2684354560:33280,2952790016:8421376,3221225472:32770,3489660928:8388610,3758096384:0,4026531840:33282,134217728:0,402653184:8421890,671088640:33282,939524096:32768,1207959552:8421888,1476395008:512,1744830464:8421378,2013265920:2,\n" +
//		"2281701376:8389120,2550136832:33280,2818572288:8421376,3087007744:8389122,3355443200:8388610,3623878656:32770,3892314112:514,4160749568:8388608,1:32768,268435457:2,536870913:8421888,805306369:8388608,1073741825:8421378,1342177281:33280,1610612737:512,1879048193:8389122,2147483649:8421890,2415919105:8421376,2684354561:8388610,2952790017:33282,3221225473:514,3489660929:8389120,3758096385:32770,4026531841:0,134217729:8421890,402653185:8421376,671088641:8388608,939524097:512,1207959553:32768,1476395009:8388610,\n" +
//		"1744830465:2,2013265921:33282,2281701377:32770,2550136833:8389122,2818572289:514,3087007745:8421888,3355443201:8389120,3623878657:0,3892314113:33280,4160749569:8421378},{\"0\":1074282512,16777216:16384,33554432:524288,50331648:1074266128,67108864:1073741840,83886080:1074282496,100663296:1073758208,117440512:16,134217728:540672,150994944:1073758224,167772160:1073741824,184549376:540688,201326592:524304,218103808:0,234881024:16400,251658240:1074266112,8388608:1073758208,25165824:540688,41943040:16,58720256:1073758224,\n" +
//		"75497472:1074282512,92274688:1073741824,109051904:524288,125829120:1074266128,142606336:524304,159383552:0,176160768:16384,192937984:1074266112,209715200:1073741840,226492416:540672,243269632:1074282496,260046848:16400,268435456:0,285212672:1074266128,301989888:1073758224,318767104:1074282496,335544320:1074266112,352321536:16,369098752:540688,385875968:16384,402653184:16400,419430400:524288,436207616:524304,452984832:1073741840,469762048:540672,486539264:1073758208,503316480:1073741824,520093696:1074282512,\n" +
//		"276824064:540688,293601280:524288,310378496:1074266112,327155712:16384,343932928:1073758208,360710144:1074282512,377487360:16,394264576:1073741824,411041792:1074282496,427819008:1073741840,444596224:1073758224,461373440:524304,478150656:0,494927872:16400,511705088:1074266128,528482304:540672},{\"0\":260,1048576:0,2097152:67109120,3145728:65796,4194304:65540,5242880:67108868,6291456:67174660,7340032:67174400,8388608:67108864,9437184:67174656,10485760:65792,11534336:67174404,12582912:67109124,13631488:65536,\n" +
//		"14680064:4,15728640:256,524288:67174656,1572864:67174404,2621440:0,3670016:67109120,4718592:67108868,5767168:65536,6815744:65540,7864320:260,8912896:4,9961472:256,11010048:67174400,12058624:65796,13107200:65792,14155776:67109124,15204352:67174660,16252928:67108864,16777216:67174656,17825792:65540,18874368:65536,19922944:67109120,20971520:256,22020096:67174660,23068672:67108868,24117248:0,25165824:67109124,26214400:67108864,27262976:4,28311552:65792,29360128:67174400,30408704:260,31457280:65796,32505856:67174404,\n" +
//		"17301504:67108864,18350080:260,19398656:67174656,20447232:0,21495808:65540,22544384:67109120,23592960:256,24641536:67174404,25690112:65536,26738688:67174660,27787264:65796,28835840:67108868,29884416:67109124,30932992:67174400,31981568:4,33030144:65792},{\"0\":2151682048,65536:2147487808,131072:4198464,196608:2151677952,262144:0,327680:4198400,393216:2147483712,458752:4194368,524288:2147483648,589824:4194304,655360:64,720896:2147487744,786432:2151678016,851968:4160,917504:4096,983040:2151682112,32768:2147487808,\n" +
//		"98304:64,163840:2151678016,229376:2147487744,294912:4198400,360448:2151682112,425984:0,491520:2151677952,557056:4096,622592:2151682048,688128:4194304,753664:4160,819200:2147483648,884736:4194368,950272:4198464,1015808:2147483712,1048576:4194368,1114112:4198400,1179648:2147483712,1245184:0,1310720:4160,1376256:2151678016,1441792:2151682048,1507328:2147487808,1572864:2151682112,1638400:2147483648,1703936:2151677952,1769472:4198464,1835008:2147487744,1900544:4194304,1966080:64,2031616:4096,1081344:2151677952,\n" +
//		"1146880:2151682112,1212416:0,1277952:4198400,1343488:4194368,1409024:2147483648,1474560:2147487808,1540096:64,1605632:2147483712,1671168:4096,1736704:2147487744,1802240:2151678016,1867776:4160,1933312:2151682048,1998848:4194304,2064384:4198464},{\"0\":128,4096:17039360,8192:262144,12288:536870912,16384:537133184,20480:16777344,24576:553648256,28672:262272,32768:16777216,36864:537133056,40960:536871040,45056:553910400,49152:553910272,53248:0,57344:17039488,61440:553648128,2048:17039488,6144:553648256,\n" +
//		"10240:128,14336:17039360,18432:262144,22528:537133184,26624:553910272,30720:536870912,34816:537133056,38912:0,43008:553910400,47104:16777344,51200:536871040,55296:553648128,59392:16777216,63488:262272,65536:262144,69632:128,73728:536870912,77824:553648256,81920:16777344,86016:553910272,90112:537133184,94208:16777216,98304:553910400,102400:553648128,106496:17039360,110592:537133056,114688:262272,118784:536871040,122880:0,126976:17039488,67584:553648256,71680:16777216,75776:17039360,79872:537133184,\n" +
//		"83968:536870912,88064:17039488,92160:128,96256:553910272,100352:262272,104448:553910400,108544:0,112640:553648128,116736:16777344,120832:262144,124928:537133056,129024:536871040},{\"0\":268435464,256:8192,512:270532608,768:270540808,1024:268443648,1280:2097152,1536:2097160,1792:268435456,2048:0,2304:268443656,2560:2105344,2816:8,3072:270532616,3328:2105352,3584:8200,3840:270540800,128:270532608,384:270540808,640:8,896:2097152,1152:2105352,1408:268435464,1664:268443648,1920:8200,2176:2097160,2432:8192,\n" +
//		"2688:268443656,2944:270532616,3200:0,3456:270540800,3712:2105344,3968:268435456,4096:268443648,4352:270532616,4608:270540808,4864:8200,5120:2097152,5376:268435456,5632:268435464,5888:2105344,6144:2105352,6400:0,6656:8,6912:270532608,7168:8192,7424:268443656,7680:270540800,7936:2097160,4224:8,4480:2105344,4736:2097152,4992:268435464,5248:268443648,5504:8200,5760:270540808,6016:270532608,6272:270540800,6528:270532616,6784:8192,7040:2105352,7296:2097160,7552:0,7808:268435456,8064:268443656},{\"0\":1048576,\n" +
//		"16:33555457,32:1024,48:1049601,64:34604033,80:0,96:1,112:34603009,128:33555456,144:1048577,160:33554433,176:34604032,192:34603008,208:1025,224:1049600,240:33554432,8:34603009,24:0,40:33555457,56:34604032,72:1048576,88:33554433,104:33554432,120:1025,136:1049601,152:33555456,168:34603008,184:1048577,200:1024,216:34604033,232:1,248:1049600,256:33554432,272:1048576,288:33555457,304:34603009,320:1048577,336:33555456,352:34604032,368:1049601,384:1025,400:34604033,416:1049600,432:1,448:0,464:34603008,480:33554433,\n" +
//		"496:1024,264:1049600,280:33555457,296:34603009,312:1,328:33554432,344:1048576,360:1025,376:34604032,392:33554433,408:34603008,424:0,440:34604033,456:1049601,472:1024,488:33555456,504:1048577},{\"0\":134219808,1:131072,2:134217728,3:32,4:131104,5:134350880,6:134350848,7:2048,8:134348800,9:134219776,10:133120,11:134348832,12:2080,13:0,14:134217760,15:133152,2147483648:2048,2147483649:134350880,2147483650:134219808,2147483651:134217728,2147483652:134348800,2147483653:133120,2147483654:133152,2147483655:32,\n" +
//		"2147483656:134217760,2147483657:2080,2147483658:131104,2147483659:134350848,2147483660:0,2147483661:134348832,2147483662:134219776,2147483663:131072,16:133152,17:134350848,18:32,19:2048,20:134219776,21:134217760,22:134348832,23:131072,24:0,25:131104,26:134348800,27:134219808,28:134350880,29:133120,30:2080,31:134217728,2147483664:131072,2147483665:2048,2147483666:134348832,2147483667:133152,2147483668:32,2147483669:134348800,2147483670:134217728,2147483671:134219808,2147483672:134350880,2147483673:134217760,\n" +
//		"2147483674:134219776,2147483675:0,2147483676:133120,2147483677:2080,2147483678:131104,2147483679:134350848}],x=[4160749569,528482304,33030144,2064384,129024,8064,504,2147483679],r=s.DES=n.extend({_doReset:function(){for(var b=this._key.words,a=[],c=0;56>c;c++){var d=q[c]-1;a[c]=b[d>>>5]>>>31-d%32&1}b=this._subKeys=[];for(d=0;16>d;d++){for(var f=b[d]=[],l=v[d],c=0;24>c;c++)f[c/6|0]|=a[(w[c]-1+l)%28]<<31-c%6,f[4+(c/6|0)]|=a[28+(w[c+24]-1+l)%28]<<31-c%6;f[0]=f[0]<<1|f[0]>>>31;for(c=1;7>c;c++)f[c]>>>=\n" +
//		"4*(c-1)+3;f[7]=f[7]<<5|f[7]>>>27}a=this._invSubKeys=[];for(c=0;16>c;c++)a[c]=b[15-c]},encryptBlock:function(b,a){this._doCryptBlock(b,a,this._subKeys)},decryptBlock:function(b,a){this._doCryptBlock(b,a,this._invSubKeys)},_doCryptBlock:function(e,a,c){this._lBlock=e[a];this._rBlock=e[a+1];u.call(this,4,252645135);u.call(this,16,65535);l.call(this,2,858993459);l.call(this,8,16711935);u.call(this,1,1431655765);for(var d=0;16>d;d++){for(var f=c[d],n=this._lBlock,p=this._rBlock,q=0,r=0;8>r;r++)q|=b[r][((p^\n" +
//		"f[r])&x[r])>>>0];this._lBlock=p;this._rBlock=n^q}c=this._lBlock;this._lBlock=this._rBlock;this._rBlock=c;u.call(this,1,1431655765);l.call(this,8,16711935);l.call(this,2,858993459);u.call(this,16,65535);u.call(this,4,252645135);e[a]=this._lBlock;e[a+1]=this._rBlock},keySize:2,ivSize:2,blockSize:2});d.DES=n._createHelper(r);s=s.TripleDES=n.extend({_doReset:function(){var b=this._key.words;this._des1=r.createEncryptor(p.create(b.slice(0,2)));this._des2=r.createEncryptor(p.create(b.slice(2,4)));this._des3=\n" +
//		"r.createEncryptor(p.create(b.slice(4,6)))},encryptBlock:function(b,a){this._des1.encryptBlock(b,a);this._des2.decryptBlock(b,a);this._des3.encryptBlock(b,a)},decryptBlock:function(b,a){this._des3.decryptBlock(b,a);this._des2.encryptBlock(b,a);this._des1.decryptBlock(b,a)},keySize:6,ivSize:2,blockSize:2});d.TripleDES=n._createHelper(s)})();\n" +
//		"\n" +
//		"/*\n" +
//		"CryptoJS v3.1.2\n" +
//		"code.google.com/p/crypto-js\n" +
//		"(c) 2009-2013 by Jeff Mott. All rights reserved.\n" +
//		"code.google.com/p/crypto-js/wiki/License\n" +
//		"*/\n" +
//		"CryptoJS.mode.ECB=function(){var a=CryptoJS.lib.BlockCipherMode.extend();a.Encryptor=a.extend({processBlock:function(a,b){this._cipher.encryptBlock(a,b)}});a.Decryptor=a.extend({processBlock:function(a,b){this._cipher.decryptBlock(a,b)}});return a}();\n" +
//		"\n" +
//		"\n" +
//		"function _mv_addr(ciphertext) {\n" +
//		"    var keyHex = CryptoJS.enc.Utf8.parse('fdrpnohyp8');\n" +
//		"    var decrypted = CryptoJS.DES.decrypt({ciphertext: CryptoJS.enc.Base64.parse(ciphertext)},\n" +
//		"        keyHex, {mode: CryptoJS.mode.ECB,padding: CryptoJS.pad.Pkcs7});\n" +
//		"    return decrypted.toString(CryptoJS.enc.Utf8);\n" +
//		"}\n" +
//		"\n" +
//		"/**\n" +
//		" * @return {string}\n" +
//		" */\n" +
//		"function GetRealPlayUrl(request) {\n" +
//		"    /*\n" +
//		"    *\n" +
//		"    * request:传入到参数，json格式\n" +
//		"    * 必须的字段:provider 平台id,quality 清晰度, content_type 视频类型 times 调用次数 tv_type 频道类型 call_back_data 回调数据(脚本返回的call_back_data字段数据，原样加到request字段,json格式的字符串)\n" +
//		"    *\n" +
//		"    * provider : 1 cntv 2 腾讯 3 爱奇艺  4 优酷 5 芒果 6 搜狐 7 米咕 8 华数TV 9 好趣 10 电视家 11 韩剧TV 12 山寨米咕\n" +
//		"    * quality  : 1 流畅 2 标清 3 高清 4 720P 5 1080P\n" +
//		"    * tv_type  : 央视 卫视 地方\n" +
//		"    *\n" +
//		"    * content_type : 4 直播 2 点播\n" +
//		"    * content_type为0时以下为必填:\n" +
//		"    * channel : cntv频道名称 migu url\n" +
//		"    * back_title:回看节目名称\n" +
//		"    * start_time:回看\n" +
//		"    * end_time:回看\n" +
//		"    * content_type 为2时以下为必填:\n" +
//		"    * url : 网页地址\n" +
//		"    * html_data : 网页内容\n" +
//		"    *\n" +
//		"    *\n" +
//		"    * 返回报文:\n" +
//		"    *json格式\n" +
//		"    * done:为true表示调用结束，不需要再次调用,false表示还需再次调用\n" +
//		"    * fetch_url_new:done为false时返回，表示下次调用是需获取网页内容的网址信息，是一个json串\n" +
//		"    * fetch_url_new_new中url表示请求网页地址，method表示请求方法(get/post),headers是请求头,body是请求体请求体\n" +
//		"    * urls:done为true时返回，表示最终的真实播放地址，是一个数组，长度为1时是完整的视频播放地址，为其他时是分段的播放地址，所有播放地址合成一个完整的视频\n" +
//		"    * issupportback:是否支持回看,0不支持，1支持\n" +
//		"    *\n" +
//		"    * */\n" +
//		"    // sohu      http://tv.sohu.com/20160727/n461224585.shtml\n" +
//		"    // iqiyi     http://www.iqiyi.com/v_19rrhc74rc.html\n" +
//		"    // youku     http://v.youku.com/v_show/id_XMTIyMDcwOTY4.html?spm=a2h1n.8251845.0.0\n" +
//		"    // mgtv      https://www.mgtv.com/b/150279/3785231.html\n" +
//		"    // qq        https://v.qq.com/x/cover/5qp15yo127kpwxy.html\n" +
//		"    // cntv      http://tv.cctv.com/live/cctv8/\n" +
//		"    // migu1     http://m.miguvideo.com/wap/resource/migu/detail/Detail_live.jsp?cid=608807427\n" +
//		"    // hanjuTV   http://www.hanju.cc/hanju/148176/1.html\n" +
//		"    // haoqu     http://www.haoqu.net/2/hunanweishi.html\n" +
//		"\n" +
//		"    var json_data;\n" +
//		"    var response;\n" +
//		"    var quality;\n" +
//		"    var data;\n" +
//		"    var re;\n" +
//		"    var i = 0;\n" +
//		"    var video_id;\n" +
//		"    var stream_id;\n" +
//		"    var issupportback = 0;\n" +
//		"    var fetch_url_new;\n" +
//		"    var fetch_url;\n" +
//		"    var result;\n" +
//		"    var begin = 0;\n" +
//		"    var end = 0;\n" +
//		"    var real_url;\n" +
//		"    var found = false;\n" +
//		"\n" +
//		"    fetch_url_new = {method:\"get\",header:{user_agent:\"\",referer:\"\",content_type:\"\"},body:\"\",url:\"\"};\n" +
//		"    json_data = JSON.parse(request);\n" +
//		"    result = {done: false, fetch_url:\"\",fetch_url_new: fetch_url_new, urls: [], issupportback: issupportback, call_back_data: {}};\n" +
//		"\n" +
//		"    if (json_data.times > 50){\n" +
//		"        fetch_url_new.url = \"\";\n" +
//		"        fetch_url = \"\";\n" +
//		"        result.done = true;\n" +
//		"        result.error = \"错误次数过多\";\n" +
//		"        return JSON.stringify(result);\n" +
//		"    }\n" +
//		"\n" +
//		"    if(typeof(json_data.quality) === \"undefined\" ){\n" +
//		"        quality = 3;\n" +
//		"    }else{\n" +
//		"        quality = json_data.quality;\n" +
//		"    }\n" +
//		"\n" +
//		"    if((json_data.provider === 1 && json_data.content_type === 4 ) || json_data.provider === 7 && json_data.content_type === 4 ){\n" +
//		"        result.issupportback = 1;\n" +
//		"    }\n" +
//		"\n" +
//		"    if(typeof(json_data.content_type) === \"undefined\" || typeof(json_data.provider) === \"undefined\"){\n" +
//		"        fetch_url_new.url = \"\";\n" +
//		"        fetch_url = \"\";\n" +
//		"        result.done = true;\n" +
//		"        return JSON.stringify(result);\n" +
//		"    }\n" +
//		"\n" +
//		"    if(typeof(json_data.content_type) === 2){\n" +
//		"        if(typeof(json_data.url) === \"undefined\" || json_data.url === \"\"){\n" +
//		"            fetch_url_new.url = \"\";\n" +
//		"            fetch_url = \"\";\n" +
//		"            result.done = true;\n" +
//		"            return JSON.stringify(result);\n" +
//		"        }\n" +
//		"        if(typeof(json_data.html_data) === \"undefined\" || json_data.html_data === \"\"){\n" +
//		"            fetch_url_new.url = \"\";\n" +
//		"            fetch_url = \"\";\n" +
//		"            result.done = true;\n" +
//		"            return JSON.stringify(result);\n" +
//		"        }\n" +
//		"    }else if(typeof(json_data.content_type) === 4){\n" +
//		"        if(typeof(json_data.channel) === \"undefined\" || json_data.channel === \"\"){\n" +
//		"            fetch_url_new.url = \"\";\n" +
//		"            fetch_url = \"\";\n" +
//		"            result.done = true;\n" +
//		"            return JSON.stringify(result);\n" +
//		"        }\n" +
//		"    }\n" +
//		"\n" +
//		"    if(json_data.channel.indexOf(\".m3u8\") > -1 || json_data.channel.indexOf(\"rtmp\") > -1){\n" +
//		"       fetch_url_new.url = \"\";\n" +
//		"       fetch_url = fetch_url_new.url;\n" +
//		"       result.done = true;\n" +
//		"       result.urls.push(json_data.channel);\n" +
//		"       return JSON.stringify(result);\n" +
//		"    }\n" +
//		"\n" +
//		"\n" +
//		"    if(json_data.content_type === 4 && json_data.channel.indexOf(\"_\") > 0){\n" +
//		"        begin = json_data.channel.indexOf(\"_\") + 1;\n" +
//		"        json_data.channel = json_data.channel.slice(begin);\n" +
//		"    }\n" +
//		"\n" +
//		"    if(json_data.provider === 8){\n" +
//		"        fetch_url_new.method = \"post\";\n" +
//		"    }\n" +
//		"\n" +
//		"\n" +
//		"\n" +
//		"    if(json_data.provider === 1){ //cctv\n" +
//		"        // http://tv.cctv.com/live/cctv1/\n" +
//		"        json_data.channel = json_data.channel.slice(json_data.channel.indexOf(\"/live/\") + 6).replace(\"/\",\"\");\n" +
//		"        if(json_data.channel === \"btv2\" || json_data.channel === \"btv3\" || json_data.channel === \"btv4\" || json_data.channel === \"btv5\" ||\n" +
//		"            json_data.channel === \"btv6\" || json_data.channel === \"btv7\" || json_data.channel === \"btv8\" || json_data.channel === \"btv9\" ||\n" +
//		"            json_data.channel === \"xiamen113\"){\n" +
//		"            json_data.tv_type = \"卫视\";\n" +
//		"        }\n" +
//		"        if(json_data.content_type === 4 || json_data.content_type === 2){\n" +
//		"            if(json_data.times === 1){\n" +
//		"                try{\n" +
//		"\n" +
//		"                    if(json_data.tv_type === \"地方\"){\n" +
//		"                        var vdn_uid = \"\";\n" +
//		"                        var vdn_tsp =new Date().getTime().toString().slice(0,10);\n" +
//		"                        var vdn_vnHtml5 = \"2049\";\t\t\t\t\t\t\t\t\t\t//央视网页Html5V1.0\n" +
//		"                        var staticCheck_Html5_02 = \"47899B86370B879139C08EA3B5E88267\";\t//第二条验证码\n" +
//		"                        var vdn_vc = md5((vdn_tsp+vdn_vnHtml5+staticCheck_Html5_02+vdn_uid)).toUpperCase();//2017年7月31日11:17:11\n" +
//		"\n" +
//		"                        //var vdnUrl = \"http://vdn.live.cntv.cn/api2/liveHtml5.do?channel=\" + addr + \"&client=html5\"+\"&tsp=\"+vdn_tsp + \"&vn=\"+ vdn_vnHtml5 + \"&vc=\"+vdn_vc + \"&uid=\"+vdn_uid + \"&wlan=\"+vdn_wlan;\n" +
//		"                        fetch_url_new.url = \"http://vdn.live.cntv.cn/api2/live.do?channel=pa://cctv_p2p_hd\" + json_data.channel + \"&tsp=1521600610&uid=&vc=\" + vdn_vc + \"&vn=3&wlan=w\";\n" +
//		"                        result.done = false;\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                    }else{\n" +
//		"                         real_url = setFlvHtml5AliNewUrl(json_data.channel);\n" +
//		"\n" +
//		"                         if(typeof(json_data.start_time) !== \"undefined\" && json_data.start_time && typeof(json_data.end_time) !== \"undefined\" && json_data.end_time){\n" +
//		"                            if((json_data.start_time + \"\").length > 10){\n" +
//		"                                json_data.start_time = json_data.start_time / 1000;\n" +
//		"                            }\n" +
//		"                            if((json_data.end_time + \"\").length > 10){\n" +
//		"                                json_data.end_time = json_data.end_time / 1000;\n" +
//		"                            }\n" +
//		"                            real_url = real_url + \"&start=\" + json_data.start_time + \"&end=\" + json_data.end_time;\n" +
//		"                        }\n" +
//		"\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls.push(real_url);\n" +
//		"                    }\n" +
//		"\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 2) {\n" +
//		"                try{\n" +
//		"                    data = json_data.html_data;\n" +
//		"                    if(data.indexOf(\"var html5VideoData=\") > -1){\n" +
//		"                        data = data.replace(\"var html5VideoData='\", \"\").replace(\"';getHtml5VideoData(html5VideoData);\", \"\");\n" +
//		"                    }\n" +
//		"                    response = JSON.parse(data);\n" +
//		"                    if (typeof(response.hls_url.hls2) !== \"undefined\" && response.hls_url.hls2) {\n" +
//		"                        real_url = response.hls_url.hls2;\n" +
//		"                    } else if (typeof(response.hls_url.hls4) !== \"undefined\" && response.hls_url.hls4) {\n" +
//		"                        real_url = response.hls_url.hls4;\n" +
//		"                    } else {\n" +
//		"                        real_url = response.hls_url.hls1;\n" +
//		"                    }\n" +
//		"\n" +
//		"                    if(typeof(json_data.start_time) !== \"undefined\" && json_data.start_time && typeof(json_data.end_time) !== \"undefined\" && json_data.end_time){\n" +
//		"                        if((json_data.start_time + \"\").length > 10){\n" +
//		"                            json_data.start_time = json_data.start_time / 1000;\n" +
//		"                        }\n" +
//		"                        if((json_data.end_time + \"\").length > 10){\n" +
//		"                            json_data.end_time = json_data.end_time / 1000;\n" +
//		"                        }\n" +
//		"                        real_url = real_url + \"&start=\" + json_data.start_time + \"&end=\" + json_data.end_time;\n" +
//		"                    }\n" +
//		"\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.push(real_url);\n" +
//		"\n" +
//		"                }catch(e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 6){ //搜狐\n" +
//		"        if(json_data.content_type === 2){//点播\n" +
//		"            if(json_data.times === 1){\n" +
//		"                fetch_url_new.url = json_data.url;\n" +
//		"                result.done = false;\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"            }else if(json_data.times === 2){\n" +
//		"                try{\n" +
//		"                    re = \"var vid=\\\"(.+?)\\\";\";\n" +
//		"                    video_id = json_data.html_data.match(re)[1];\n" +
//		"                    fetch_url_new.url = \"https://m.tv.sohu.com/phone_playinfo?callback=jsonpx1517539948922_56_11&vid=\" + video_id + \"&site=1&appid=tv&api_key=f351515304020cad28c92f70f002261c&plat=17&sver=1.0&partner=1&uid=1711141102362091&muid=1517539772399202&_c=1&pt=5&qd=680&src=11060001&ssl=1&_=1517539948922\";\n" +
//		"                    result.done = false;\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 3){\n" +
//		"                try{\n" +
//		"                    //https://m.tv.sohu.com/phone_playinfo?callback=jsonpx1517539948922_56_11&vid=3055369&site=1&appid=tv&api_key=f351515304020cad28c92f70f002261c&plat=17&sver=1.0&partner=1&uid=1711141102362091&muid=1517539772399202&_c=1&pt=5&qd=680&src=11060001&ssl=1&_=1517539948922\n" +
//		"                    data = json_data.html_data;\n" +
//		"                    data = data.replace(\"jsonpx1517539948922_56_11(\",\"\");\n" +
//		"                    data = data.replace(\"})\",\"}\");\n" +
//		"                    response = JSON.parse(data);\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    if(quality === 2){ // nor 640 * 360\n" +
//		"                        result.urls.push(response.data.urls.m3u8.nor[0]);\n" +
//		"                    }else if(quality === 4){ //sup 1280 * 720\n" +
//		"                        result.urls.push(response.data.urls.m3u8.sup[0]);\n" +
//		"                    }else if(quality === 3){ // high 852 * 480\n" +
//		"                        result.urls.push(response.data.urls.m3u8.hig[0]);\n" +
//		"                    }else{\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls.splice(0,result.urls.length);\n" +
//		"                        return JSON.stringify(result);\n" +
//		"                    }\n" +
//		"\n" +
//		"                    if(result.urls[0] === \"\" || typeof(result.urls[0]) === \"undefined\"){\n" +
//		"                        result.urls.splice(0,result.urls.length);\n" +
//		"                        if(response.data.urls.m3u8.sup[0] !== \"\" && typeof(response.data.urls.m3u8.sup[0]) !== \"undefined\"){\n" +
//		"                            result.urls.push(response.data.urls.m3u8.sup[0]);\n" +
//		"                        }else if(response.data.urls.m3u8.hig[0] !== \"\" && typeof(response.data.urls.m3u8.hig[0]) !== \"undefined\"){\n" +
//		"                            result.urls.push(response.data.urls.m3u8.hig[0]);\n" +
//		"                        }else if(response.data.urls.m3u8.nor[0] !== \"\" && typeof(response.data.urls.m3u8.nor[0]) !== \"undefined\" ){\n" +
//		"                            result.urls.push(response.data.urls.m3u8.nor[0]);\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else{\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 3){ //爱奇艺\n" +
//		"        if(json_data.content_type === 2){\n" +
//		"            if(json_data.times === 1){\n" +
//		"                fetch_url_new.url = json_data.url;\n" +
//		"                result.done = false;\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"            }else if(json_data.times === 2){\n" +
//		"                try{\n" +
//		"                    re = \"tvId:[0-9]+\";\n" +
//		"                    var tvid;\n" +
//		"                    data = json_data.html_data;\n" +
//		"                    tvid = data.match(re)[0];\n" +
//		"                    if (tvid !== null && typeof(tvid) !== \"undefined\" && tvid !== \"\") {\n" +
//		"                        tvid = tvid.split(\":\")[1];\n" +
//		"                    }\n" +
//		"                    if (tvid === null && typeof(tvid) === \"undefined\" && tvid === \"\") {\n" +
//		"                        re = \"data-player-tvid=\\\"[0-9]+\\\"\";\n" +
//		"                        tvid = data.match(re)[0];\n" +
//		"                        if (tvid !== null && typeof(tvid) !== \"undefined\" && tvid !== \"\") {\n" +
//		"                            tvid = tvid.split(\"=\")[1];\n" +
//		"                            tvid = tvid.substr(1,tvid.length - 2);\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"\n" +
//		"                    if(json_data.html_data.indexOf(\"data-videolist-vid=\") > -1){\n" +
//		"                        re = \"data-videolist-vid=\\\"[0-9a-zA-Z]+\\\"\";\n" +
//		"                        video_id = data.match(re);\n" +
//		"                        if (video_id !== null && typeof(video_id) !== \"undefined\" && video_id !== \"\") {\n" +
//		"                            video_id = video_id[0];\n" +
//		"                            video_id = video_id.split(\"=\")[1];\n" +
//		"                        }\n" +
//		"                    }else if(json_data.html_data.indexOf(\"data-player-videoid=\") > -1){\n" +
//		"                        re = 'data-player-videoid=\"([^\"]+)\"';\n" +
//		"                        video_id = data.match(re)[0];\n" +
//		"\n" +
//		"                        if (video_id !== null && typeof(video_id) !== \"undefined\" && video_id !== \"\") {\n" +
//		"                            video_id = video_id.split(\"=\")[1];\n" +
//		"                        }\n" +
//		"                    }else{\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls.splice(0,result.urls.length);\n" +
//		"                        return JSON.stringify(result);\n" +
//		"                    }\n" +
//		"\n" +
//		"                    video_id = video_id.replace(\"\\\"\",\"\");\n" +
//		"                    video_id = video_id.replace(\"\\\"\",\"\");\n" +
//		"                    var src = \"76f90cbd92f94a2e925d83e8ccd22cb7\";\n" +
//		"                    var key = \"d5fb4bd9d50c4be6948c97edd7254b0e\";\n" +
//		"                    var t1 = parseInt(new Date().getTime() / 1000);\n" +
//		"                    data = t1 + key + video_id ;\n" +
//		"                    var sc = md5(data);\n" +
//		"\n" +
//		"                    fetch_url_new.url = \"http://cache.m.iqiyi.com/tmts/\" + tvid + \"/\" + video_id + \"/?t=\" + t1 + \"&sc=\" + sc + \"&src=\" + src;\n" +
//		"                    result.done = false;\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"\n" +
//		"            }else if(json_data.times === 3) try {\n" +
//		"                response = JSON.parse(json_data.html_data);\n" +
//		"                var streams = response.data.vidl;\n" +
//		"                var cdn_urls;\n" +
//		"                for (i = 0; i < streams.length; i++) {\n" +
//		"                    /*quality  : 1 流畅 2 标清 3 高清 4 720P 5 1080P*/\n" +
//		"                    //vd 2 848 * 480\n" +
//		"                    //vd 4 1280 * 720\n" +
//		"                    //vd 1 592 * 336\n" +
//		"                    //vd 96 416 * 240\n" +
//		"                    if (quality === 1 && (streams[i].vd === 1 || streams[i].vd === 96)) {\n" +
//		"                        cdn_urls = streams[i].m3u;\n" +
//		"                        break;\n" +
//		"                    } else if (quality === 3 && (streams[i].vd === 2)) {\n" +
//		"                        cdn_urls = streams[i].m3u;\n" +
//		"                        break;\n" +
//		"                    } else if (quality === 4 && (streams[i].vd === 4)) {\n" +
//		"                        cdn_urls = streams[i].m3u;\n" +
//		"                        break;\n" +
//		"                    }\n" +
//		"                }\n" +
//		"\n" +
//		"                if(cdn_urls === \"\" || typeof(cdn_urls) === \"undefined\"){\n" +
//		"                    cdn_urls = streams[0].m3u\n" +
//		"                }\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.push(cdn_urls);\n" +
//		"            } catch (e) {\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0, result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 5){ //芒果\n" +
//		"        if(json_data.content_type === 2){\n" +
//		"            if(json_data.times === 1){\n" +
//		"                fetch_url_new.url = json_data.url;\n" +
//		"                result.done = false;\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"            }else if(json_data.times === 2){\n" +
//		"                try{\n" +
//		"                    video_id = json_data.url.split(\"/\")[5];\n" +
//		"                    video_id = video_id.replace(\".html\",\"\");\n" +
//		"\n" +
//		"                    fetch_url_new.url = \"http://pcweb.api.mgtv.com/player/video?video_id=\" + video_id ;\n" +
//		"                    result.done = false;\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                }catch (e){\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 3){\n" +
//		"                try{\n" +
//		"                    response = JSON.parse(json_data.html_data);\n" +
//		"                    var domain = response.data.stream_domain[0];\n" +
//		"\n" +
//		"                    /*\n" +
//		"                    标清:848 * 480\n" +
//		"                    高清:1104 * 624\n" +
//		"                    超清:1280 * 720\n" +
//		"                    */\n" +
//		"                    for(i = 0 ; i < response.data.stream.length ; i++){\n" +
//		"                        if((response.data.stream[i].name === \"标清\" && quality === 2) || (response.data.stream[i].name === \"高清\" && quality === 3)||(response.data.stream[i].name === \"高清\" && quality === 4)){\n" +
//		"                            fetch_url_new.url = domain + response.data.stream[i].url;\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            found = true;\n" +
//		"                            break;\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"\n" +
//		"                    result.done = found === false;\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 4){\n" +
//		"                try{\n" +
//		"                    response = JSON.parse(json_data.html_data);\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.push(response.info);\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 2){ //腾讯\n" +
//		"        if(json_data.content_type === 2){\n" +
//		"            if(json_data.times === 1){\n" +
//		"                fetch_url_new.url = json_data.url;\n" +
//		"                result.done = false;\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"            }else if(json_data.times === 2){\n" +
//		"                try{\n" +
//		"                    re = \"<link rel=\\\"canonical\\\" href=\\\"(.+?)\\\" />\";\n" +
//		"                    video_id = json_data.html_data.match(re)[0];\n" +
//		"                    if (video_id !== null && typeof(video_id) !== \"undefined\" && video_id !== \"\") {\n" +
//		"                        re = \"href=\\\"(.+?)\\\"\";\n" +
//		"                        video_id = video_id.match(re)[0];\n" +
//		"                        video_id = video_id.replace(\"href=\",\"\").replace(\"\\\"\",\"\").replace(\"\\\"\",\"\");\n" +
//		"                        video_id = video_id.split(\"/\");\n" +
//		"                        video_id = video_id[video_id.length - 1];\n" +
//		"                        video_id = video_id.replace(\".html\",\"\");\n" +
//		"\n" +
//		"                        fetch_url_new.url = \"http://vv.video.qq.com/getinfo?otype=json&appver=3.2.19.333&platform=11&defnpayver=1&vid=\"+video_id;\n" +
//		"                        result.done = false;\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                    }else{\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls.splice(0,result.urls.length);\n" +
//		"                        return JSON.stringify(result);\n" +
//		"\n" +
//		"                    }\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 3){\n" +
//		"                try{\n" +
//		"                    data = json_data.html_data;\n" +
//		"                    data = data.replace(\"QZOutputJson=\",\"\");\n" +
//		"                    data = data.replace(\"}};\",\"}}\");\n" +
//		"                    response = JSON.parse(data);\n" +
//		"                    var seg_cnt = response.vl.vi[0].cl.fc;\n" +
//		"                    if(seg_cnt === 0){\n" +
//		"                        seg_cnt = 1;\n" +
//		"                    }\n" +
//		"                    var quality_count = response.fl.fi.length;\n" +
//		"                    var quality_id;\n" +
//		"                    if(quality === 0){\n" +
//		"                        quality_id = response.fl.fi[0].id;\n" +
//		"                    }else{\n" +
//		"                        quality_id = response.fl.fi[quality_count - 1].id;\n" +
//		"                    }\n" +
//		"\n" +
//		"                    var fn_pre = response.vl.vi[0].lnk;\n" +
//		"                    video_id = response.vl.vi[0].vid;\n" +
//		"                    var file_name = fn_pre + \".p\" + quality_id % 10000 + \".\" + 1 + \".mp4\";\n" +
//		"                    result.fetch_url_new.url = \"http://vv.video.qq.com/getkey?otype=json&platform=11&format=\" + quality_id + \"&vid=\" + video_id + \"&appver=3.2.19.333\" + \"&filename=\" + file_name;\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = false;\n" +
//		"                    result.call_back_data.count = seg_cnt - 1;\n" +
//		"                    result.call_back_data.urls = [];\n" +
//		"                    result.call_back_data.host = response.vl.vi[0].ul.ui[0].url;\n" +
//		"                    result.call_back_data.url = response.vl.vi[0].ul.ui[0].url;\n" +
//		"                    result.call_back_data.vkey = response.vl.vi[0].fvkey;\n" +
//		"                    result.call_back_data.fn_pre = fn_pre;\n" +
//		"                    result.call_back_data.file_name = file_name;\n" +
//		"                    result.call_back_data.quality_id = quality_id;\n" +
//		"                    result.call_back_data.video_id = video_id;\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"\n" +
//		"            }else if(json_data.times >= 4){\n" +
//		"                try{\n" +
//		"                    data = json_data.html_data;\n" +
//		"                    data = data.replace(\"QZOutputJson=\",\"\");\n" +
//		"                    data = data.replace(\"};\",\"}\");\n" +
//		"                    response = JSON.parse(data);\n" +
//		"                    if(json_data.call_back_data.count === 0){\n" +
//		"                        fetch_url_new.url = \"http://vv.video.qq.com/getkey?otype=json&platform=11&format=\" + json_data.call_back_data.quality_id + \"&vid=\" + json_data.call_back_data.video_id + \"&appver=3.2.19.333\" + \"&filename=\" + json_data.call_back_data.file_name;\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls = json_data.call_back_data.urls;\n" +
//		"                    }else{\n" +
//		"                        var vkey;\n" +
//		"                        // vkey = json_data.call_back_data.vkey;\n" +
//		"                        // real_url =  json_data.call_back_data.host + json_data.call_back_data.file_name + \"?vkey=\" + response.key;\n" +
//		"                        vkey = response.key;\n" +
//		"                        real_url = json_data.call_back_data.host + json_data.call_back_data.file_name + \"?vkey=\" + vkey;\n" +
//		"\n" +
//		"                        // if(typeof(response.key) !== \"undefined\" && response.key){\n" +
//		"                        //     vkey = response.key;\n" +
//		"                        //     real_url = json_data.call_back_data.host + json_data.call_back_data.file_name + \"?vkey=\" + vkey;\n" +
//		"                        // }else{\n" +
//		"                        //     vkey = json_data.call_back_data.vkey;\n" +
//		"                        //     real_url =  json_data.call_back_data.host + json_data.call_back_data.fn_pre  + \".mp4?vkey=\" + vkey;\n" +
//		"                        // }\n" +
//		"                        json_data.call_back_data.urls.push(real_url);\n" +
//		"                        result.call_back_data.urls = json_data.call_back_data.urls;\n" +
//		"                        //fetch_url_new.url = key_api;\n" +
//		"                        result.done = false;\n" +
//		"\n" +
//		"                        var file_name = json_data.call_back_data.fn_pre + \".p\" + (json_data.call_back_data.quality_id % 10000) + \".\" + (json_data.times - 2) + \".mp4\";\n" +
//		"                        fetch_url_new.url = \"http://vv.video.qq.com/getkey?otype=json&platform=11&format=\" + json_data.call_back_data.quality_id + \"&vid=\" + json_data.call_back_data.video_id + \"&appver=3.2.19.333\" + \"&filename=\" + file_name;\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.call_back_data.count = json_data.call_back_data.count - 1;\n" +
//		"                        result.call_back_data.host = json_data.call_back_data.host;\n" +
//		"                        result.call_back_data.url = json_data.call_back_data.url;\n" +
//		"                        result.call_back_data.vkey = json_data.call_back_data.vkey;\n" +
//		"                        result.call_back_data.fn_pre = json_data.call_back_data.fn_pre;\n" +
//		"                        result.call_back_data.file_name = file_name;\n" +
//		"                        result.call_back_data.quality_id = json_data.call_back_data.quality_id;\n" +
//		"                        result.call_back_data.video_id = json_data.call_back_data.video_id;\n" +
//		"                    }\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 12){ //山寨米咕直播\n" +
//		"        if(json_data.content_type === 4){\n" +
//		"            if(json_data.times === 1) {\n" +
//		"                try{\n" +
//		"                    if(json_data.tv_type === \"央视\"){//央视\n" +
//		"                        fetch_url_new.url = \"http://mini.chinasosuo.cc/wap/cctv/\"+ json_data.channel + \".html\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = false;\n" +
//		"                        result.call_back_data = {type:\"cctv\"};\n" +
//		"                    }else if(json_data.tv_type === \"卫视\"){  //卫视\n" +
//		"                        fetch_url_new.url = \"http://mini.chinasosuo.cc/wap/tv/\"+ json_data.channel + \".html\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = false;\n" +
//		"                        result.call_back_data = {type:\"weishi\"};\n" +
//		"                    }else{//地方\n" +
//		"                        //http://mini.chinasosuo.cc/zhibo/difang/hubei/wuhan2.html\n" +
//		"                        fetch_url_new.url = \"http://mini.chinasosuo.cc/zhibo/difang/\" + json_data.provice + \"/\" + json_data.channel + \".html\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = false;\n" +
//		"                        result.call_back_data = {province:json_data.province};\n" +
//		"                    }\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 2){\n" +
//		"                try{\n" +
//		"                    if(json_data.tv_type === \"央视\"){\n" +
//		"                        if(json_data.html_data.indexOf(\"<video\") !== -1) {\n" +
//		"                            re = \"http.*m3u8\";\n" +
//		"                            real_url = json_data.html_data.match(re)[0];\n" +
//		"\n" +
//		"                            fetch_url_new.url = \"\";\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            result.done = true;\n" +
//		"                            result.urls.push(real_url);\n" +
//		"                        }else if(json_data.html_data.indexOf(\"iframe\") !== -1) {\n" +
//		"                            re = \"http.*type=cq\";\n" +
//		"                            fetch_url_new.url = json_data.html_data.match(re)[0];\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            result.done = false;\n" +
//		"                        }\n" +
//		"                    }else if(json_data.tv_type === \"卫视\"){\n" +
//		"                        if(json_data.html_data.indexOf(\"<video\") !== -1) {\n" +
//		"                            re = \"http.*m3u8\";\n" +
//		"                            real_url = json_data.html_data.match(re)[0];\n" +
//		"\n" +
//		"                            fetch_url_new.url = \"\";\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            result.done = true;\n" +
//		"                            result.urls.push(real_url);\n" +
//		"                        }else if(json_data.html_data.indexOf(\"iframe\") !== -1) {\n" +
//		"                            re = \"http.*type=cq\";\n" +
//		"                            fetch_url_new.url = json_data.html_data.match(re)[0];\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            result.done = false;\n" +
//		"                        }\n" +
//		"                    }else if(json_data.tv_type === \"地方\"){\n" +
//		"                        if(json_data.html_data.indexOf(\"url\") !== -1){\n" +
//		"                            re = \"url=http.*flv\";\n" +
//		"                            real_url = json_data.html_data.match(re)[0];\n" +
//		"                            real_url = real_url.slice(5);\n" +
//		"\n" +
//		"                            fetch_url_new.url = \"\";\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            result.done = true;\n" +
//		"                            result.urls.push(real_url);\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }else if(json_data.times === 3){\n" +
//		"                try{\n" +
//		"                    if(json_data.html_data.indexOf(\"encodeURIComponent\") !== -1){\n" +
//		"                        re = \"encodeURIComponent(.*)\";\n" +
//		"                        real_url = json_data.html_data.match(re)[0];\n" +
//		"                        begin = real_url.indexOf(\"(\");\n" +
//		"                        end = real_url.indexOf(\")\");\n" +
//		"                        real_url = real_url.slice(begin + 2,end - 1);\n" +
//		"\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls.push(real_url);\n" +
//		"                    }else{\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        fetch_url = fetch_url_new.url;\n" +
//		"                        result.done = true;\n" +
//		"                        result.urls.splice(0,result.urls.length);\n" +
//		"                        return JSON.stringify(result);\n" +
//		"                    }\n" +
//		"                }catch (e){\n" +
//		"                    result.error = e.toString();\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 7){//咪咕直播\n" +
//		"        // http://m.miguvideo.com/wap/resource/migu/detail/DetailLive_data.jsp?cid=609154254&range=0\n" +
//		"        //http://m.miguvideo.com/wap/resource/migu/detail/DetailLiveBackSee_data.jsp?playbillId=60915425420180314006&indx=1\n" +
//		"        if(json_data.times === 1){\n" +
//		"            try{\n" +
//		"                begin = json_data.channel.indexOf(\"cid=\");\n" +
//		"                var cid = json_data.channel.slice(begin + 4);\n" +
//		"                fetch_url_new.url = \"http://h5spdegrade.miguvideo.com/wap/resource/migu/detail/detail_Live_data.jsp?cid=\" + cid + \"&range=0\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = false;\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                data = json_data.html_data;\n" +
//		"                response = JSON.parse(data);\n" +
//		"\n" +
//		"                var isback = 0;\n" +
//		"                if(typeof(json_data.start_time) !== \"undefined\" && json_data.start_time && typeof(json_data.end_time) !== \"undefined\" && json_data.end_time) {\n" +
//		"                    isback = 1;\n" +
//		"                }\n" +
//		"\n" +
//		"                var t = new Date(json_data.start_time);\n" +
//		"                var start_time = getTime(t);\n" +
//		"\n" +
//		"                for(i = 0 ; i < response.length ; i ++){\n" +
//		"                    if(isback === 0){\n" +
//		"                        if(response[i].billtime === response[i].nowtime){\n" +
//		"                            //1 流畅 2 标清 3 高清 4 720P 5 1080P\n" +
//		"                            //idx: 1 标清 2 高清 3 超清\n" +
//		"                            var indx = 1;\n" +
//		"                            if(quality === 4){\n" +
//		"                                indx = 3;\n" +
//		"                            }else if(quality === 3){\n" +
//		"                                indx = 2;\n" +
//		"                            }else if(quality === 2){\n" +
//		"                                indx = 1;\n" +
//		"                            }else{\n" +
//		"                                fetch_url_new.url = \"\";\n" +
//		"                                fetch_url = fetch_url_new.url;\n" +
//		"                                result.done = true;\n" +
//		"                                result.urls.splice(0,result.urls.length);\n" +
//		"                                return JSON.stringify(result);\n" +
//		"                            }\n" +
//		"\n" +
//		"                            fetch_url_new.url = \"http://h5spdegrade.miguvideo.com/wap/resource/migu/detail/detail_LiveBackSee_data.jsp?playbillId=\" + response[i].playbillId + \"&indx=\" + indx;\n" +
//		"                            result.done = false;\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"\n" +
//		"                            break;\n" +
//		"                        }\n" +
//		"                    }else if(isback === 1){\n" +
//		"                        //节目名称一样并且开始时间小于1分钟认为是同一个epg\n" +
//		"                        if(json_data.back_title === response[i].playName){\n" +
//		"                            if(start_time === response[i].billtime){\n" +
//		"                                fetch_url_new.url = \"\";\n" +
//		"                                fetch_url = fetch_url_new.url;\n" +
//		"                                result.done = true;\n" +
//		"                                result.urls.push(response[i].url);\n" +
//		"                                break;\n" +
//		"                            }else if (Math.abs(parseInt(start_time.slice(8)) - parseInt(response[i].billtime.slice(8)) < 60)){\n" +
//		"                                fetch_url_new.urls = \"\";\n" +
//		"                                fetch_url = fetch_url_new.url;\n" +
//		"                                result.done = true;\n" +
//		"                                result.urls.push(response[i].url);\n" +
//		"                                break;\n" +
//		"                            }\n" +
//		"                        }\n" +
//		"\n" +
//		"                        response[i].playStartTime = response[i].playStartTime.replace(\":\",\"\");\n" +
//		"                        response[i].playEndTime = response[i].playEndTime.replace(\":\",\"\");\n" +
//		"\n" +
//		"                        //如果系统的epg开始时间在咪咕的epg，则认为是同一个epg\n" +
//		"                        if(start_time.slice(8,12) >= response[i].playStartTime && start_time.slice(8,12) < response[i].playEndTime){\n" +
//		"                            fetch_url_new.url = \"\";\n" +
//		"                            fetch_url = fetch_url_new.url;\n" +
//		"                            result.done = true;\n" +
//		"                            eval(response[i].deEncrptJsFunc);\n" +
//		"                            result.urls.push(_mv_addr(response[i].url));\n" +
//		"                            break;\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"                }\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"\n" +
//		"        }else if(json_data.times === 3){\n" +
//		"            try{\n" +
//		"                data = json_data.html_data;\n" +
//		"                response = JSON.parse(data);\n" +
//		"\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                eval(response.deEncrptJsFunc);\n" +
//		"                result.urls.push(_mv_addr(response.liveBackPlayUrl));\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 11){//hanjuTV\n" +
//		"        if(json_data.times === 1){\n" +
//		"            fetch_url_new.url = json_data.url;\n" +
//		"            result.done = false;\n" +
//		"            fetch_url = fetch_url_new.url;\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                //vid='https://www3.youku00.com/20180304/KJaKFqDw/index.m3u8'\n" +
//		"                re = \"vid='.*.m3u8\";\n" +
//		"                real_url = json_data.html_data.match(re)[0];\n" +
//		"                real_url = real_url.slice(5);\n" +
//		"\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.push(real_url);\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 10){//电视家\n" +
//		"        if(json_data.times === 1){\n" +
//		"            fetch_url_new.url = \"http://cdn.idianshijia.com/api/channel/groupSimplifiedChinese_162\";\n" +
//		"            fetch_url = fetch_url_new.url;\n" +
//		"            result.done = false;\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                data = json_data.html_data;\n" +
//		"                response = JSON.parse(data);\n" +
//		"                for(i = 0 ; i < response.data.length ; i++){\n" +
//		"                    if(found === true){\n" +
//		"                        break;\n" +
//		"                    }\n" +
//		"                    for(j = 0 ; j < response.data[i].channels.length ; j++){\n" +
//		"                        if(found === true){\n" +
//		"                            break;\n" +
//		"                        }\n" +
//		"                        if(json_data.channel === response.data[i].channels[j].name){\n" +
//		"                            for(k = 0 ; k < response.data[i].channels[j].streams.length ; k++){\n" +
//		"                                if(response.data[i].channels[j].streams[k].url.indexOf(\"m3u8\") > 0){\n" +
//		"                                    real_url = response.data[i].channels[j].streams[k].url;\n" +
//		"                                    found = true;\n" +
//		"                                    break\n" +
//		"                                }\n" +
//		"                            }\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"                }\n" +
//		"\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.push(real_url);\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 8){//华数Phone\n" +
//		"        if(json_data.times === 1){\n" +
//		"            try{\n" +
//		"                var quality_code;\n" +
//		"                if(quality === 3){\n" +
//		"                    quality_code = 1000995; /*标清 720 * 576*/\n" +
//		"                }else if(quality === 4){\n" +
//		"                    /*高清 1280* 720*/\n" +
//		"                    quality_code = 1000996;\n" +
//		"                }else{\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"                fetch_url_new.header = {\"Content-Type\":\"application/xml\"};\n" +
//		"                fetch_url_new.body = \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n\" +\n" +
//		"                    \"<message module=\\\"CATALOG_SERVICE\\\" version=\\\"1.0\\\">\\n\" +\n" +
//		"                    \"\\t<header action=\\\"REQUEST\\\" command=\\\"CONTENT_QUERY\\\" sequence=\\\"20121030212732_103861\\\" component-id=\\\"SYSTEM2\\\" component-type=\\\"THIRD_PARTY_SYSTEM\\\" />\\n\" +\n" +
//		"                    \"\\t<body>\\n\" +\n" +
//		"                    \"\\t\\t<contents>\\n\" +\n" +
//		"                    \"\\t\\t\\t<content>\\n\" +\n" +
//		"                    \"\\t\\t\\t\\t<code>\" + json_data.channel + \"</code>\\n\" +\n" +
//		"                    \"\\t\\t\\t\\t<site-code>1000889</site-code>\\n\" +\n" +
//		"                    \"\\t\\t\\t\\t<items-index>-1</items-index>\\n\" +\n" +
//		"                    \"\\t\\t\\t\\t<folder-code>\" + quality_code + \"</folder-code>\\n\" +\n" +
//		"                    \"\\t\\t\\t\\t<format>-1</format>\\n\" +\n" +
//		"                    \"\\t\\t\\t</content>\\n\" +\n" +
//		"                    \"\\t\\t</contents>\\n\" +\n" +
//		"                    \"\\t</body>\\n\" +\n" +
//		"                    \"</message>\";\n" +
//		"                fetch_url_new.url = \"http://101.71.69.172:8080/wasu_catalog/catalog\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = false;\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                // var parser=new DOMParser();//创建文档对象\n" +
//		"                // var xmldoc = parser.parseFromString(json_data.html_data,\"text/xml\");\n" +
//		"                // //提取数据\n" +
//		"                // var playUrl = xmldoc.getElementsByTagName(\"playUrl\")[0].textContent;\n" +
//		"                // playUrl = playUrl.replace(/\\s+/g,\"\");\n" +
//		"                re = \"<playUrl>.*</playUrl>\";\n" +
//		"                var players = json_data.html_data.match(re)[0].toString();\n" +
//		"                players = players.replace(\"<![CDATA[\",\"\").replace(\"]]>\",\"\").replace(/\\s+/g,\"\").replace(\"<playUrl>\",\"\").replace(\"</playUrl>\",\"\");\n" +
//		"\n" +
//		"                fetch_url_new.url = \"\"\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.push(players);\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString()\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 9){ //好趣\n" +
//		"        if(json_data.times === 1){\n" +
//		"            fetch_url_new.url = json_data.channel;\n" +
//		"            fetch_url = fetch_url_new.url;\n" +
//		"            result.done = false;\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                re = \"<li data-player=\\\"[0-9]*\\\".*</span></li>\";\n" +
//		"                var playlist = json_data.html_data.match(re)[0];\n" +
//		"                var bq = playlist.split(\"</li>\");\n" +
//		"                result.call_back_data.fetch_stream_ids = [];\n" +
//		"                for(i = 0 ; i < bq.length - 1 ; i++ ){\n" +
//		"                    // <li data-player=\"6926\" class=\"tab-item btn btn-syc btn-syc-select\"><span class=\"s\">WASU标清</span>\n" +
//		"                    stream_id = bq[i].split(\"=\")[1].split(\" \")[0].replace('\"','').replace('\"','');\n" +
//		"                    var quality_name = bq[i].split(\">\")[2].split(\"<\")[0].replace('\"','').replace('\"','');\n" +
//		"                    //result.call_back_data.fetch_stream_ids = [];\n" +
//		"                    if((quality === 2 && quality_name.indexOf(\"标清\") > -1) || (quality === 3 && quality_name.indexOf(\"高清\") > -1) || quality_name.indexOf(\"\") > -1){\n" +
//		"                        result.call_back_data.fetch_stream_ids.push(stream_id);\n" +
//		"                    }\n" +
//		"                }\n" +
//		"\n" +
//		"                if(result[\"call_back_data\"][\"fetch_stream_ids\"].length === 0){\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"\n" +
//		"                fetch_url_new.url = \"http://www.haoqu.net/e/extend/tv.php?id=\" + result.call_back_data.fetch_stream_ids[0];\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = false;\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString()\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }else if(json_data.times > 2){\n" +
//		"            try{\n" +
//		"                //signal = 'WASU����$rtmp://qmtvrtmp-al.wasu.cn/live10/hd_hnws$flv';\n" +
//		"                if(json_data.times === (json_data.call_back_data.fetch_stream_ids.length - 2)){\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"                re = \"signal = '.*';\";\n" +
//		"                var players = json_data.html_data.match(re)[0];\n" +
//		"                if(players.indexOf(\"m3u8\") < 0 && players.indexOf(\"rtmp\") < 0){\n" +
//		"                    if(json_data.times - 2 > json_data.call_back_data.fetch_stream_ids.length){\n" +
//		"                        fetch_url_new.url = \"\";\n" +
//		"                        result.done = true;\n" +
//		"                    }else{\n" +
//		"                        fetch_url_new.url = \"http://www.haoqu.net/e/extend/tv.php?id=\" + json_data.call_back_data.fetch_stream_ids[json_data.times - 2];\n" +
//		"                        result.done = false;\n" +
//		"                    }\n" +
//		"\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.call_back_data = json_data.call_back_data;\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"                var url = players.split(\"$\")[1].toString();\n" +
//		"                if(url === null || url === \"\" || typeof(url) === \"undefined\" ){\n" +
//		"                    fetch_url_new.url = \"http://www.haoqu.net/e/extend/tv.php?id=\" + json_data.call_back_data.fetch_stream_ids[json_data.times - 2];\n" +
//		"                    result.done = false;\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.call_back_data = json_data.call_back_data;\n" +
//		"                }else{\n" +
//		"                    fetch_url_new.url = \"\";\n" +
//		"                    fetch_url = fetch_url_new.url;\n" +
//		"                    result.done = true;\n" +
//		"                    //http://player.haoqu.net/swf.html?id=http://www.panjinhonghaitan.com/listen/flash/flvplayer.swf?flashvars=&src=rtmp://60.22.145.101:1935/live/live3&autoHideControlBar=true&streamType=live&autoPlay=true&verbose=true\n" +
//		"                    if(url.indexOf(\"http\") > -1 && url.indexOf(\"rtmp\") > -1){\n" +
//		"                        begin = url.indexOf(\"rtmp\")\n" +
//		"                        url = url.slice(begin)\n" +
//		"                        if(url.indexOf(\"&\") > -1){\n" +
//		"                            end = url.indexOf(\"&\")\n" +
//		"                            url = url.slice(0,end)\n" +
//		"                        }\n" +
//		"                    }\n" +
//		"                    result.urls.push(url);\n" +
//		"                }\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 13){//华数TV\n" +
//		"        if(json_data.times === 1){\n" +
//		"            fetch_url_new.url = \"http://clientapi.wasu.cn/Phone/liveinfo/id/\" + json_data.channel;\n" +
//		"            result.done = false;\n" +
//		"            fetch_url = fetch_url_new.url;\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                data = json_data.html_data;\n" +
//		"                response = JSON.parse(data);\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.push(response[0].p2pinfo.videoid);\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"    }else if(json_data.provider === 4){//优酷\n" +
//		"        if(json_data.times === 1){\n" +
//		"            var ccode;\n" +
//		"            if(json_data.provider === 4){\n" +
//		"                ccode = \"0502\";\n" +
//		"            }else{\n" +
//		"                ccode = \"0512\"; //土豆\n" +
//		"            }\n" +
//		"            var ckey = \"DIl58SLFxFNndSV1GFNnMQVYkx1PP5tKe1siZu/86PR1u/Wh1Ptd%2BWOZsHHWxysSfAOhNJpdVWsdVJNsfJ8Sxd8WKVvNfAS8aS8fAOzYARzPyPc3JvtnPHjTdKfESTdnuTW6ZPvk2pNDh4uFzotgdMEFkzQ5wZVXl2Pf1/Y6hLK0OnCNxBj3%2Bnb0v72gZ6b0td%2BWOZsHHWxysSo/0y9D2K42SaB8Y/%2BaD2K42SaB8Y/%2BahU%2BWOZsHcrxysooUeND\";\n" +
//		"            var b64p = '([a-zA-Z0-9=]+)'\n" +
//		"            var p_list = ['youku\\.com/v_show/id_' + b64p,\n" +
//		"                      'player\\.youku\\.com/player\\.php/sid/' + b64p + '/v\\.swf',\n" +
//		"                      'loader\\.swf\\?VideoIDS=' + b64p,\n" +
//		"                      'tudou.com/v/' + b64p,\n" +
//		"                      'player\\.youku\\.com/embed/' + b64p]\n" +
//		"\n" +
//		"            var vid;\n" +
//		"            for(var i = 0 ; i < p_list.length ; i++ ){\n" +
//		"                var hit = json_data.url.match(p_list[i]);\n" +
//		"                if (hit !== null && typeof(hit) !== \"undefined\") {\n" +
//		"                    vid = hit[1];\n" +
//		"                    break;\n" +
//		"                }\n" +
//		"            }\n" +
//		"\n" +
//		"            if(vid === \"\" || typeof(vid) === \"undefined\"){\n" +
//		"                result.error = \"get vid error!!!\";\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"\n" +
//		"            var utid = \"caSJEpwANg0CAatxaKRHU3CU\";//getCookie(\"cna\");\n" +
//		"            var url = \"https://ups.youku.com/ups/get.json?vid=\" + vid + \"&ccode=\" + ccode;\n" +
//		"            url += \"&client_ip=192.168.1.1\";\n" +
//		"            url += \"&utid=\" + utid;\n" +
//		"            url += \"&client_ts=\" + Date.parse(new Date()).toString().slice(0,10);\n" +
//		"            url += \"&ckey=\" + ckey;\n" +
//		"\n" +
//		"            //result.fetch_url_new.header = {\"UserAgent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36\",\"Referer\":\"http://v.youku.com\"};\n" +
//		"            result.fetch_url_new.header.user_agent = \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36\";\n" +
//		"            result.fetch_url_new.header.referer = \"http://v.youku.com\";\n" +
//		"            //result.fetch_url_new.header = {user_agent:\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36\",referer:\"http://v.youku.com\"};\n" +
//		"            result.fetch_url_new.method = \"get\";\n" +
//		"            result.fetch_url_new.url = url;\n" +
//		"            result.done = false;\n" +
//		"            result.fetch_url = result.fetch_url_new.url;\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                var api_data = JSON.parse(json_data.html_data);\n" +
//		"\n" +
//		"                var youku_quality;\n" +
//		"                var quality_id;\n" +
//		"                for(var i = 0 ; i < api_data.data.stream.length ; i++ ){\n" +
//		"                    youku_quality = api_data.data.stream[i].stream_type;\n" +
//		"                    if(youku_quality === \"hd3\" || youku_quality === \"hd3v2\" || youku_quality === \"mp4hd3\" || youku_quality === \"mp4hd3v2\"){\n" +
//		"                         quality_id = 5;\n" +
//		"                    }else if(youku_quality === \"hd2\" || youku_quality === \"hd2v2\" || youku_quality === \"mp4hd2\" || youku_quality === \"mp4hd2v2\"){\n" +
//		"                         quality_id = 4;\n" +
//		"                    }else if(youku_quality === \"mp4hd\"){\n" +
//		"                         quality_id = 3;\n" +
//		"                    }else if(youku_quality === \"mp4sd\" || youku_quality === \"flv\" || youku_quality === \"mp4\"){\n" +
//		"                         quality_id = 2;\n" +
//		"                    }else if(youku_quality === \"flvhd\" || youku_quality === \"3gphd\"){\n" +
//		"                         quality_id = 1;\n" +
//		"                    }\n" +
//		"                    if(quality_id === quality){\n" +
//		"                         result.urls.push(api_data.data.stream[i].m3u8_url);\n" +
//		"                         break;\n" +
//		"                    }\n" +
//		"                }\n" +
//		"\n" +
//		"                if(result.urls[0] === \"\" || typeof(result.urls[0]) === \"undefined\") {\n" +
//		"                    result.urls.push(api_data.data.stream[0].m3u8_url);\n" +
//		"                }\n" +
//		"\n" +
//		"                result.done = true;\n" +
//		"                result.fetch_url_new.url = \"\";\n" +
//		"                result.fetch_url = \"\";\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url_new.url = \"\";\n" +
//		"                fetch_url = fetch_url_new.url;\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"\n" +
//		"        }\n" +
//		"    }\n" +
//		"    return JSON.stringify(result);\n" +
//		"}\n" +
//		"\n" +
//		"function getTime(date) {\n" +
//		"    var y = date.getFullYear();\n" +
//		"    var m = date.getMonth() + 1;\n" +
//		"    m = m < 10 ? ('0' + m) : m;\n" +
//		"    var d = date.getDate();\n" +
//		"    d = d < 10 ? ('0' + d) : d;\n" +
//		"    var h = date.getHours();\n" +
//		"    h=h < 10 ? ('0' + h) : h;\n" +
//		"    var minute = date.getMinutes();\n" +
//		"    minute = minute < 10 ? ('0' + minute) : minute;\n" +
//		"    var second=date.getSeconds();\n" +
//		"    second=second < 10 ? ('0' + second) : second;\n" +
//		"    return y + m + d + h + minute + second;\n" +
//		"}\n" +
//		"\n" +
//		"function datetimeToUnix(datetime) {\n" +
//		"    var tmp_datetime = datetime.replace(/:/g, '-');\n" +
//		"    tmp_datetime = tmp_datetime.replace(/ /g, '-');\n" +
//		"    var arr = tmp_datetime.split(\"-\");\n" +
//		"    var now = new Date(Date.UTC(arr[0], arr[1] - 1, arr[2], arr[3] - 8, arr[4], arr[5]));\n" +
//		"    return parseInt(now.getTime() / 1000);\n" +
//		"}\n"
//
//	return jsCode
//}
//
