package script

import (
	"background/common/logger"
	"github.com/robertkrimen/otto"
	"encoding/json"
	"io/ioutil"
	"net/http"
)
func GetIqiyiRealPlayUrl(url string)(string){
	jsCode := GetiqiyiJsCode()

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
			logger.Debug(string(recv))
			reqParam.HtmlData = string(recv)

			//logger.Debug("html_data:",reqParam.HtmlData)

		}

		b, _ := json.Marshal(reqParam)
		logger.Debug(string(b))

		vm := otto.New()
		vm.Run(jsCode)
		data ,err := vm.Call("GetIqiyiRealPlayUrl",nil,string(b))
		if err != nil{
			logger.Error(err)
			return ""
		}
		result = data.String()
		logger.Debug(result)

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

func GetiqiyiJsCode()(string){
	jsCode := "!function(a){\"use strict\";function b(a,b){var c=(65535&a)+(65535&b),d=(a>>16)+(b>>16)+(c>>16);return d<<16|65535&c}function c(a,b){return a<<b|a>>>32-b}function d(a,d,e,f,g,h){return b(c(b(b(d,a),b(f,h)),g),e)}function e(a,b,c,e,f,g,h){return d(b&c|~b&e,a,b,f,g,h)}function f(a,b,c,e,f,g,h){return d(b&e|c&~e,a,b,f,g,h)}function g(a,b,c,e,f,g,h){return d(b^c^e,a,b,f,g,h)}function h(a,b,c,e,f,g,h){return d(c^(b|~e),a,b,f,g,h)}function i(a,c){a[c>>5]|=128<<c%32,a[(c+64>>>9<<4)+14]=c;var d,i,j,k,l,m=1732584193,n=-271733879,o=-1732584194,p=271733878;for(d=0;d<a.length;d+=16)i=m,j=n,k=o,l=p,m=e(m,n,o,p,a[d],7,-680876936),p=e(p,m,n,o,a[d+1],12,-389564586),o=e(o,p,m,n,a[d+2],17,606105819),n=e(n,o,p,m,a[d+3],22,-1044525330),m=e(m,n,o,p,a[d+4],7,-176418897),p=e(p,m,n,o,a[d+5],12,1200080426),o=e(o,p,m,n,a[d+6],17,-1473231341),n=e(n,o,p,m,a[d+7],22,-45705983),m=e(m,n,o,p,a[d+8],7,1770035416),p=e(p,m,n,o,a[d+9],12,-1958414417),o=e(o,p,m,n,a[d+10],17,-42063),n=e(n,o,p,m,a[d+11],22,-1990404162),m=e(m,n,o,p,a[d+12],7,1804603682),p=e(p,m,n,o,a[d+13],12,-40341101),o=e(o,p,m,n,a[d+14],17,-1502002290),n=e(n,o,p,m,a[d+15],22,1236535329),m=f(m,n,o,p,a[d+1],5,-165796510),p=f(p,m,n,o,a[d+6],9,-1069501632),o=f(o,p,m,n,a[d+11],14,643717713),n=f(n,o,p,m,a[d],20,-373897302),m=f(m,n,o,p,a[d+5],5,-701558691),p=f(p,m,n,o,a[d+10],9,38016083),o=f(o,p,m,n,a[d+15],14,-660478335),n=f(n,o,p,m,a[d+4],20,-405537848),m=f(m,n,o,p,a[d+9],5,568446438),p=f(p,m,n,o,a[d+14],9,-1019803690),o=f(o,p,m,n,a[d+3],14,-187363961),n=f(n,o,p,m,a[d+8],20,1163531501),m=f(m,n,o,p,a[d+13],5,-1444681467),p=f(p,m,n,o,a[d+2],9,-51403784),o=f(o,p,m,n,a[d+7],14,1735328473),n=f(n,o,p,m,a[d+12],20,-1926607734),m=g(m,n,o,p,a[d+5],4,-378558),p=g(p,m,n,o,a[d+8],11,-2022574463),o=g(o,p,m,n,a[d+11],16,1839030562),n=g(n,o,p,m,a[d+14],23,-35309556),m=g(m,n,o,p,a[d+1],4,-1530992060),p=g(p,m,n,o,a[d+4],11,1272893353),o=g(o,p,m,n,a[d+7],16,-155497632),n=g(n,o,p,m,a[d+10],23,-1094730640),m=g(m,n,o,p,a[d+13],4,681279174),p=g(p,m,n,o,a[d],11,-358537222),o=g(o,p,m,n,a[d+3],16,-722521979),n=g(n,o,p,m,a[d+6],23,76029189),m=g(m,n,o,p,a[d+9],4,-640364487),p=g(p,m,n,o,a[d+12],11,-421815835),o=g(o,p,m,n,a[d+15],16,530742520),n=g(n,o,p,m,a[d+2],23,-995338651),m=h(m,n,o,p,a[d],6,-198630844),p=h(p,m,n,o,a[d+7],10,1126891415),o=h(o,p,m,n,a[d+14],15,-1416354905),n=h(n,o,p,m,a[d+5],21,-57434055),m=h(m,n,o,p,a[d+12],6,1700485571),p=h(p,m,n,o,a[d+3],10,-1894986606),o=h(o,p,m,n,a[d+10],15,-1051523),n=h(n,o,p,m,a[d+1],21,-2054922799),m=h(m,n,o,p,a[d+8],6,1873313359),p=h(p,m,n,o,a[d+15],10,-30611744),o=h(o,p,m,n,a[d+6],15,-1560198380),n=h(n,o,p,m,a[d+13],21,1309151649),m=h(m,n,o,p,a[d+4],6,-145523070),p=h(p,m,n,o,a[d+11],10,-1120210379),o=h(o,p,m,n,a[d+2],15,718787259),n=h(n,o,p,m,a[d+9],21,-343485551),m=b(m,i),n=b(n,j),o=b(o,k),p=b(p,l);return[m,n,o,p]}function j(a){var b,c=\"\";for(b=0;b<32*a.length;b+=8)c+=String.fromCharCode(a[b>>5]>>>b%32&255);return c}function k(a){var b,c=[];for(c[(a.length>>2)-1]=void 0,b=0;b<c.length;b+=1)c[b]=0;for(b=0;b<8*a.length;b+=8)c[b>>5]|=(255&a.charCodeAt(b/8))<<b%32;return c}function l(a){return j(i(k(a),8*a.length))}function m(a,b){var c,d,e=k(a),f=[],g=[];for(f[15]=g[15]=void 0,e.length>16&&(e=i(e,8*a.length)),c=0;16>c;c+=1)f[c]=909522486^e[c],g[c]=1549556828^e[c];return d=i(f.concat(k(b)),512+8*b.length),j(i(g.concat(d),640))}function n(a){var b,c,d=\"0123456789abcdef\",e=\"\";for(c=0;c<a.length;c+=1)b=a.charCodeAt(c),e+=d.charAt(b>>>4&15)+d.charAt(15&b);return e}function o(a){return unescape(encodeURIComponent(a))}function p(a){return l(o(a))}function q(a){return n(p(a))}function r(a,b){return m(o(a),o(b))}function s(a,b){return n(r(a,b))}function t(a,b,c){return b?c?r(b,a):s(b,a):c?p(a):q(a)}\"function\"==typeof define&&define.amd?define(function(){return t}):a.md5=t}(this);\n" +
		"\n" +
		"/**\n" +
		" * @return {string}\n" +
		" */\n" +
		"function GetIqiyiRealPlayUrl(request) {\n" +
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
		"            var re = \"tvId:[0-9]+\";\n" +
		"            var data = json_data.html_data;\n" +
		"            var tvid = data.match(re)[0];\n" +
		"            if (tvid !== null && typeof(tvid) !== \"undefined\" && tvid !== \"\") {\n" +
		"                tvid = tvid.split(\":\")[1];\n" +
		"            }\n" +
		"            if (tvid === null && typeof(tvid) === \"undefined\" && tvid === \"\") {\n" +
		"                re = \"data-player-tvid=\\\"[0-9]+\\\"\";\n" +
		"                tvid = data.match(re)[0];\n" +
		"                if (tvid !== null && typeof(tvid) !== \"undefined\" && tvid !== \"\") {\n" +
		"                    tvid = tvid.split(\"=\")[1];\n" +
		"                    tvid = tvid.substr(1,tvid.length - 2);\n" +
		"                }\n" +
		"            }\n" +
		"\n" +
		"            var video_id;\n" +
		"            if(json_data.html_data.indexOf(\"data-player-videoid=\") > -1){\n" +
		"                re = \"data-player-videoid=\\\"[0-9a-zA-Z]+\\\"\";\n" +
		"                video_id = data.match(re);\n" +
		"\n" +
		"                if (video_id !== null && typeof(video_id) !== \"undefined\" && video_id !== \"\") {\n" +
		"                    video_id = video_id[0];\n" +
		"                    video_id = video_id.split(\"=\")[1];\n" +
		"                }\n" +
		"            }else if(json_data.html_data.indexOf(\"data-videolist-vid=\") > -1){\n" +
		"                re = 'data-player-videoid=\"([^\"]+)\"';\n" +
		"                video_id = data.match(re)[0];\n" +
		"\n" +
		"                if (video_id !== null && typeof(video_id) !== \"undefined\" && video_id !== \"\") {\n" +
		"                    video_id = video_id.split(\"=\")[1];\n" +
		"                }\n" +
		"            }else{\n" +
		"                fetch_url.url = \"\";\n" +
		"                result.done = true;\n" +
		"                result.urls.splice(0,result.urls.length);\n" +
		"                return JSON.stringify(result);\n" +
		"            }\n" +
		"\n" +
		"            video_id = video_id.replace(\"\\\"\",\"\");\n" +
		"            video_id = video_id.replace(\"\\\"\",\"\");\n" +
		"            var src = \"76f90cbd92f94a2e925d83e8ccd22cb7\";\n" +
		"            var key = \"d5fb4bd9d50c4be6948c97edd7254b0e\";\n" +
		"            var t1 = parseInt(new Date().getTime() / 1000);\n" +
		"            data = t1 + key + video_id ;\n" +
		"            var sc = md5(data);\n" +
		"\n" +
		"            fetch_url.url = \"http://cache.m.iqiyi.com/tmts/\" + tvid + \"/\" + video_id + \"/?t=\" + t1 + \"&sc=\" + sc + \"&src=\" + src;\n" +
		"            result.done = false;\n" +
		"        }catch (e){\n" +
		"            result.error = e.toString();\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            result.urls.splice(0,result.urls.length);\n" +
		"            return JSON.stringify(result);\n" +
		"        }\n" +
		"\n" +
		"    }else if(json_data.times === 3) try {\n" +
		"        var response = JSON.parse(json_data.html_data);\n" +
		"        var streams = response.data.vidl;\n" +
		"        var cdn_urls;\n" +
		"\n" +
		"        for (i = 0; i < streams.length; i++) {\n" +
		"            /*quality  : 1 流畅 2 标清 3 高清 4 720P 5 1080P*/\n" +
		"            //vd 2 848 * 480\n" +
		"            //vd 4 1280 * 720\n" +
		"            //vd 1 592 * 336\n" +
		"            //vd 96 416 * 240\n" +
		"            if (quality === 1 && (streams[i].vd === 1 || streams[i].vd === 96)) {\n" +
		"                cdn_urls = streams[i].m3u;\n" +
		"                break;\n" +
		"            } else if (quality === 3 && (streams[i].vd === 2)) {\n" +
		"                cdn_urls = streams[i].m3u;\n" +
		"                break;\n" +
		"            } else if (quality === 4 && (streams[i].vd === 4)) {\n" +
		"                cdn_urls = streams[i].m3u;\n" +
		"                break;\n" +
		"            }\n" +
		"        }\n" +
		"\n" +
		"        if(cdn_urls === \"\" || typeof(cdn_urls) === \"undefined\"){\n" +
		"            cdn_urls = streams[0].m3u\n" +
		"        }\n" +
		"        fetch_url.url = \"\";\n" +
		"        result.done = true;\n" +
		"        result.urls.push(cdn_urls);\n" +
		"    } catch (e) {\n" +
		"        result.error = e.toString();\n" +
		"        fetch_url.url = \"\";\n" +
		"        result.done = true;\n" +
		"        result.urls.splice(0, result.urls.length);\n" +
		"        return JSON.stringify(result);\n" +
		"    }\n" +
		"    return JSON.stringify(result);\n" +
		"}\n" +
		"\n"
	return jsCode
}