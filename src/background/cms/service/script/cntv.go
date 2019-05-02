package script

import (
	"background/common/logger"
	"github.com/robertkrimen/otto"
	"encoding/json"
	"io/ioutil"
	"net/http"
)
func GetCntvRealPlayUrl(url string)(string){
	jsCode := GetCntvJsCode()

	type CallBackDatas struct{
		FetchUrl   string      `gorm:"fetch_url" json:"fetch_url"`
	}

	type ReqParam struct{
		Times             uint32    `gorm:"times" json:"times"`
		Quality           uint8     `gorm:"quality" json:"quality"`
		ContentType       uint32    `gorm:"content_type" json:"content_type"`
		HtmlData          string    `gorm:"html_data" json:"html_data"`
		Channel           string    `gorm:"channel" json:"channel"`
		CallBackData      CallBackDatas  `gorm:"call_back_data" json:"call_back_data"`
		TvType            string    `gorm:"tv_type" json:"tv_type"`
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
		reqParam.Channel = url
		reqParam.TvType = "央视"
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
		data ,err := vm.Call("GetCntvRealPlayUrl",nil,string(b))
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

func GetCntvJsCode()(string){
	jsCode := "!function(a){\"use strict\";function b(a,b){var c=(65535&a)+(65535&b),d=(a>>16)+(b>>16)+(c>>16);return d<<16|65535&c}function c(a,b){return a<<b|a>>>32-b}function d(a,d,e,f,g,h){return b(c(b(b(d,a),b(f,h)),g),e)}function e(a,b,c,e,f,g,h){return d(b&c|~b&e,a,b,f,g,h)}function f(a,b,c,e,f,g,h){return d(b&e|c&~e,a,b,f,g,h)}function g(a,b,c,e,f,g,h){return d(b^c^e,a,b,f,g,h)}function h(a,b,c,e,f,g,h){return d(c^(b|~e),a,b,f,g,h)}function i(a,c){a[c>>5]|=128<<c%32,a[(c+64>>>9<<4)+14]=c;var d,i,j,k,l,m=1732584193,n=-271733879,o=-1732584194,p=271733878;for(d=0;d<a.length;d+=16)i=m,j=n,k=o,l=p,m=e(m,n,o,p,a[d],7,-680876936),p=e(p,m,n,o,a[d+1],12,-389564586),o=e(o,p,m,n,a[d+2],17,606105819),n=e(n,o,p,m,a[d+3],22,-1044525330),m=e(m,n,o,p,a[d+4],7,-176418897),p=e(p,m,n,o,a[d+5],12,1200080426),o=e(o,p,m,n,a[d+6],17,-1473231341),n=e(n,o,p,m,a[d+7],22,-45705983),m=e(m,n,o,p,a[d+8],7,1770035416),p=e(p,m,n,o,a[d+9],12,-1958414417),o=e(o,p,m,n,a[d+10],17,-42063),n=e(n,o,p,m,a[d+11],22,-1990404162),m=e(m,n,o,p,a[d+12],7,1804603682),p=e(p,m,n,o,a[d+13],12,-40341101),o=e(o,p,m,n,a[d+14],17,-1502002290),n=e(n,o,p,m,a[d+15],22,1236535329),m=f(m,n,o,p,a[d+1],5,-165796510),p=f(p,m,n,o,a[d+6],9,-1069501632),o=f(o,p,m,n,a[d+11],14,643717713),n=f(n,o,p,m,a[d],20,-373897302),m=f(m,n,o,p,a[d+5],5,-701558691),p=f(p,m,n,o,a[d+10],9,38016083),o=f(o,p,m,n,a[d+15],14,-660478335),n=f(n,o,p,m,a[d+4],20,-405537848),m=f(m,n,o,p,a[d+9],5,568446438),p=f(p,m,n,o,a[d+14],9,-1019803690),o=f(o,p,m,n,a[d+3],14,-187363961),n=f(n,o,p,m,a[d+8],20,1163531501),m=f(m,n,o,p,a[d+13],5,-1444681467),p=f(p,m,n,o,a[d+2],9,-51403784),o=f(o,p,m,n,a[d+7],14,1735328473),n=f(n,o,p,m,a[d+12],20,-1926607734),m=g(m,n,o,p,a[d+5],4,-378558),p=g(p,m,n,o,a[d+8],11,-2022574463),o=g(o,p,m,n,a[d+11],16,1839030562),n=g(n,o,p,m,a[d+14],23,-35309556),m=g(m,n,o,p,a[d+1],4,-1530992060),p=g(p,m,n,o,a[d+4],11,1272893353),o=g(o,p,m,n,a[d+7],16,-155497632),n=g(n,o,p,m,a[d+10],23,-1094730640),m=g(m,n,o,p,a[d+13],4,681279174),p=g(p,m,n,o,a[d],11,-358537222),o=g(o,p,m,n,a[d+3],16,-722521979),n=g(n,o,p,m,a[d+6],23,76029189),m=g(m,n,o,p,a[d+9],4,-640364487),p=g(p,m,n,o,a[d+12],11,-421815835),o=g(o,p,m,n,a[d+15],16,530742520),n=g(n,o,p,m,a[d+2],23,-995338651),m=h(m,n,o,p,a[d],6,-198630844),p=h(p,m,n,o,a[d+7],10,1126891415),o=h(o,p,m,n,a[d+14],15,-1416354905),n=h(n,o,p,m,a[d+5],21,-57434055),m=h(m,n,o,p,a[d+12],6,1700485571),p=h(p,m,n,o,a[d+3],10,-1894986606),o=h(o,p,m,n,a[d+10],15,-1051523),n=h(n,o,p,m,a[d+1],21,-2054922799),m=h(m,n,o,p,a[d+8],6,1873313359),p=h(p,m,n,o,a[d+15],10,-30611744),o=h(o,p,m,n,a[d+6],15,-1560198380),n=h(n,o,p,m,a[d+13],21,1309151649),m=h(m,n,o,p,a[d+4],6,-145523070),p=h(p,m,n,o,a[d+11],10,-1120210379),o=h(o,p,m,n,a[d+2],15,718787259),n=h(n,o,p,m,a[d+9],21,-343485551),m=b(m,i),n=b(n,j),o=b(o,k),p=b(p,l);return[m,n,o,p]}function j(a){var b,c=\"\";for(b=0;b<32*a.length;b+=8)c+=String.fromCharCode(a[b>>5]>>>b%32&255);return c}function k(a){var b,c=[];for(c[(a.length>>2)-1]=void 0,b=0;b<c.length;b+=1)c[b]=0;for(b=0;b<8*a.length;b+=8)c[b>>5]|=(255&a.charCodeAt(b/8))<<b%32;return c}function l(a){return j(i(k(a),8*a.length))}function m(a,b){var c,d,e=k(a),f=[],g=[];for(f[15]=g[15]=void 0,e.length>16&&(e=i(e,8*a.length)),c=0;16>c;c+=1)f[c]=909522486^e[c],g[c]=1549556828^e[c];return d=i(f.concat(k(b)),512+8*b.length),j(i(g.concat(d),640))}function n(a){var b,c,d=\"0123456789abcdef\",e=\"\";for(c=0;c<a.length;c+=1)b=a.charCodeAt(c),e+=d.charAt(b>>>4&15)+d.charAt(15&b);return e}function o(a){return unescape(encodeURIComponent(a))}function p(a){return l(o(a))}function q(a){return n(p(a))}function r(a,b){return m(o(a),o(b))}function s(a,b){return n(r(a,b))}function t(a,b,c){return b?c?r(b,a):s(b,a):c?p(a):q(a)}\"function\"==typeof define&&define.amd?define(function(){return t}):a.md5=t}(this);\n" +
		"\n" +
		"var _0xda49=[\"\x75\x64\x6B\x32\x37\x76\x6E\x38\x6C\x64\x66\x33\x6C\x63\x76\x31\x73\x70\",\"\x70\x61\x72\x73\x65\",\"\x30\",\"\x72\x61\x6E\x64\x6F\x6D\",\"\x66\x6C\x6F\x6F\x72\",\n" +
		"    \"\x2F\x63\x6E\x74\x76\x6C\x69\x76\x65\x2F\",\"\x6D\x64\x2E\x6D\x33\x75\x38\",\"\x2D\",\n" +
		"    \"\x64\x6C\x39\x32\x66\x39\x63\x6A\x68\x33\x68\x38\x32\x76\x63\x36\x32\x6B\x78\x61\x6C\x69\x77\x6C\x31\x66\",\"\x73\x75\x62\x73\x74\x72\x69\x6E\x67\",\n" +
		"    \"\x68\x74\x74\x70\x3A\x2F\x2F\x68\x6C\x73\x32\x2E\x63\x6E\x74\x76\x2E\x6D\x79\x61\x6C\x69\x63\x64\x6E\x2E\x63\x6F\x6D\",\"\x3F\x61\x75\x74\x68\x5F\x6B\x65\x79\x3D\",\n" +
		"    \"\x4C\x49\x56\x45\x2D\x48\x4C\x53\x2D\x43\x44\x4E\x2D\x54\x58\x59\",\n" +
		"    \"\x68\x74\x74\x70\x3A\x2F\x2F\x63\x63\x74\x76\x35\x2E\x74\x78\x74\x79\x2E\x35\x32\x31\x33\x2E\x6C\x69\x76\x65\x70\x6C\x61\x79\x2E\x6D\x79\x71\x63\x6C\x6F\x75\x64\x2E\x63\x6F\x6D\x2F\x6C\x69\x76\x65\x2F\",\n" +
		"    \"\x5F\x74\x78\x74\x79\x2E\x6D\x33\x75\x38\",\"\x31\",\"\x73\x72\x63\",\"\x5F\x68\x74\x6D\x6C\x35\x50\x6C\x61\x79\x65\x72\",\"\x68\x75\x61\x77\x65\x69\",\n" +
		"    \"\x69\x6E\x64\x65\x78\x4F\x66\",\"\x74\x6F\x4C\x6F\x77\x65\x72\x43\x61\x73\x65\",\"\x75\x73\x65\x72\x41\x67\x65\x6E\x74\",\"\x6C\x6F\x61\x64\",\n" +
		"    \"\x64\x69\x61\x6E\x70\x69\x61\x6E\x2E\x6D\x70\x34\",\"\x75\x6E\x64\x65\x66\x69\x6E\x65\x64\",\"\x6C\x65\x6E\x67\x74\x68\",\"\x70\x6F\x73\x74\x65\x72\",\n" +
		"    \"\x4C\x49\x56\x45\x2D\x48\x44\x53\x2D\x43\x44\x4E\x2D\x41\x4C\x49\"];\n" +
		"\n" +
		"function setFlvHtml5AliNewUrl(videoName) {\n" +
		"    var _0x56eax3 = _0xda49[0];\n" +
		"    var _0x56eax4 = Date[_0xda49[1]](new Date()) / 1000;\n" +
		"    var _0x56eax5 = _0xda49[2];\n" +
		"    var _0x56eax6 = Math[_0xda49[4]](Math[_0xda49[3]]() * 1000);\n" +
		"    var _0x56eax7 = _0xda49[5] + videoName + _0xda49[6];\n" +
		"    var _0x1576=[\"\x63\x75\x6C\x33\", \"\x63\x71\x6C\x36\"];\n" +
		"    var html5Aauth= _0x1576[0];\n" +
		"    _0x1576[1]=\"\x63\x67\x6C\x35\";\n" +
		"    var html5CdnStr = _0x1576[1] + html5Aauth;\n" +
		"    var _0x56eax8 = md5(_0x56eax7 + _0xda49[7] + _0x56eax4 + _0xda49[7] + _0x56eax6 + _0xda49[7] + _0x56eax5 + _0xda49[7] + _0xda49[8][19] + html5Aauth + _0x56eax3[_0xda49[9]](5, 11) + html5CdnStr);\n" +
		"    var _0x56eax9 = _0x56eax4 + _0xda49[7] + _0x56eax6 + _0xda49[7] + _0x56eax5 + _0xda49[7] + _0x56eax8;\n" +
		"    var _0x56eaxa = _0xda49[10] + _0x56eax7 + _0xda49[11] + _0x56eax9;\n" +
		"    var _0x56eaxb = _0x56eaxa;\n" +
		"    return _0x56eaxb;\n" +
		"}\n" +
		"\n" +
		"/**\n" +
		" * @return {string}\n" +
		" */\n" +
		"function GetCntvRealPlayUrl(request) {\n" +
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
		"    json_data.channel = json_data.channel.slice(json_data.channel.indexOf(\"/live/\") + 6).replace(\"/\",\"\");\n" +
		"    if(json_data.channel === \"btv2\" || json_data.channel === \"btv3\" || json_data.channel === \"btv4\" || json_data.channel === \"btv5\" ||\n" +
		"        json_data.channel === \"btv6\" || json_data.channel === \"btv7\" || json_data.channel === \"btv8\" || json_data.channel === \"btv9\" ||\n" +
		"        json_data.channel === \"xiamen113\" || json_data.channel === \"cztv1\"){\n" +
		"        json_data.tv_type = \"卫视\";\n" +
		"    }\n" +
		"    if(json_data.times === 1){\n" +
		"        try{\n" +
		"\n" +
		"            if(json_data.tv_type === \"地方\"){\n" +
		"                var vdn_uid = \"\";\n" +
		"                var vdn_tsp =new Date().getTime().toString().slice(0,10);\n" +
		"                var vdn_vnHtml5 = \"2049\";\t\t\t\t\t\t\t\t\t\t//央视网页Html5V1.0\n" +
		"                var staticCheck_Html5_02 = \"47899B86370B879139C08EA3B5E88267\";\t//第二条验证码\n" +
		"                var vdn_vc = md5((vdn_tsp+vdn_vnHtml5+staticCheck_Html5_02+vdn_uid)).toUpperCase();//2017年7月31日11:17:11\n" +
		"\n" +
		"                //var vdnUrl = \"http://vdn.live.cntv.cn/api2/liveHtml5.do?channel=\" + addr + \"&client=html5\"+\"&tsp=\"+vdn_tsp + \"&vn=\"+ vdn_vnHtml5 + \"&vc=\"+vdn_vc + \"&uid=\"+vdn_uid + \"&wlan=\"+vdn_wlan;\n" +
		"                fetch_url.url = \"http://vdn.live.cntv.cn/api2/live.do?channel=pa://cctv_p2p_hd\" + json_data.channel + \"&tsp=1521600610&uid=&vc=\" + vdn_vc + \"&vn=3&wlan=w\";\n" +
		"                result.done = false;\n" +
		"            }else{\n" +
		"                 var real_url = setFlvHtml5AliNewUrl(json_data.channel);\n" +
		"\n" +
		"                 if((json_data.start_time + \"\").length > 10){\n" +
		"                     json_data.start_time = json_data.start_time / 1000;\n" +
		"                     json_data.end_time = 9999999999;\n" +
		"                     real_url = real_url + \"&start=\" + json_data.start_time + \"&end=\" + json_data.end_time;\n" +
		"                 }\n" +
		"\n" +
		"                 fetch_url.url = \"\";\n" +
		"                 result.done = true;\n" +
		"                 result.urls.push(real_url);\n" +
		"            }\n" +
		"\n" +
		"        }catch (e){\n" +
		"            result.error = e.toString();\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            return JSON.stringify(result);\n" +
		"        }\n" +
		"    }else if(json_data.times === 2) {\n" +
		"        try{\n" +
		"            var data = json_data.html_data;\n" +
		"            if(data.indexOf(\"var html5VideoData=\") > -1){\n" +
		"                data = data.replace(\"var html5VideoData='\", \"\").replace(\"';getHtml5VideoData(html5VideoData);\", \"\");\n" +
		"            }\n" +
		"            var real_url;\n" +
		"            var response = JSON.parse(data);\n" +
		"            if (typeof(response.hls_url.hls2) !== \"undefined\" && response.hls_url.hls2) {\n" +
		"                real_url = response.hls_url.hls2;\n" +
		"            } else if (typeof(response.hls_url.hls4) !== \"undefined\" && response.hls_url.hls4) {\n" +
		"                real_url = response.hls_url.hls4;\n" +
		"            } else {\n" +
		"                real_url = response.hls_url.hls1;\n" +
		"            }\n" +
		"\n" +
		"            if(typeof(json_data.start_time) !== \"undefined\" || json_data.start_time){\n" +
		"                if((json_data.start_time + \"\").length > 10){\n" +
		"                    json_data.start_time = json_data.start_time / 1000;\n" +
		"                }\n" +
		"                json_data.end_time = 9999999999;\n" +
		"\n" +
		"                real_url = real_url + \"&start=\" + json_data.start_time + \"&end=\" + json_data.end_time;\n" +
		"            }\n" +
		"\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            result.urls.push(real_url);\n" +
		"\n" +
		"        }catch(e){\n" +
		"            result.error = e.toString();\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            return JSON.stringify(result);\n" +
		"        }\n" +
		"    }\n" +
		"\n" +
		"    return JSON.stringify(result);\n" +
		"}\n" +
		"\n"

	return jsCode
}