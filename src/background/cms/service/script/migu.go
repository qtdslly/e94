package script

import (
	"background/common/logger"
	"github.com/robertkrimen/otto"
	"encoding/json"
	"io/ioutil"
	"net/http"
)
func GetMiguRealPlayUrl(content_type uint32,url string)(string){
	jsCode := GetMiguJsCode()

	type CallBackDatas struct{
		FetchUrlNew   string      `gorm:"fetch_url_new" json:"fetch_url_new"`
		VId           string      `gorm:"vId" json:"vId"`
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
		reqParam.ContentType = content_type
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
		data ,err := vm.Call("GetMiguRealPlayUrl",nil,string(b))
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

func GetMiguJsCode()(string){

	jsCode := "!function (a) {\n" +
		"    \"use strict\";\n" +
		"\n" +
		"    function b(a, b) {\n" +
		"        var c = (65535 & a) + (65535 & b), d = (a >> 16) + (b >> 16) + (c >> 16);\n" +
		"        return d << 16 | 65535 & c\n" +
		"    }\n" +
		"\n" +
		"    function c(a, b) {\n" +
		"        return a << b | a >>> 32 - b\n" +
		"    }\n" +
		"\n" +
		"    function d(a, d, e, f, g, h) {\n" +
		"        return b(c(b(b(d, a), b(f, h)), g), e)\n" +
		"    }\n" +
		"\n" +
		"    function e(a, b, c, e, f, g, h) {\n" +
		"        return d(b & c | ~b & e, a, b, f, g, h)\n" +
		"    }\n" +
		"\n" +
		"    function f(a, b, c, e, f, g, h) {\n" +
		"        return d(b & e | c & ~e, a, b, f, g, h)\n" +
		"    }\n" +
		"\n" +
		"    function g(a, b, c, e, f, g, h) {\n" +
		"        return d(b ^ c ^ e, a, b, f, g, h)\n" +
		"    }\n" +
		"\n" +
		"    function h(a, b, c, e, f, g, h) {\n" +
		"        return d(c ^ (b | ~e), a, b, f, g, h)\n" +
		"    }\n" +
		"\n" +
		"    function i(a, c) {\n" +
		"        a[c >> 5] |= 128 << c % 32, a[(c + 64 >>> 9 << 4) + 14] = c;\n" +
		"        var d, i, j, k, l, m = 1732584193, n = -271733879, o = -1732584194, p = 271733878;\n" +
		"        for (d = 0; d < a.length; d += 16) i = m, j = n, k = o, l = p, m = e(m, n, o, p, a[d], 7, -680876936), p = e(p, m, n, o, a[d + 1], 12, -389564586), o = e(o, p, m, n, a[d + 2], 17, 606105819), n = e(n, o, p, m, a[d + 3], 22, -1044525330), m = e(m, n, o, p, a[d + 4], 7, -176418897), p = e(p, m, n, o, a[d + 5], 12, 1200080426), o = e(o, p, m, n, a[d + 6], 17, -1473231341), n = e(n, o, p, m, a[d + 7], 22, -45705983), m = e(m, n, o, p, a[d + 8], 7, 1770035416), p = e(p, m, n, o, a[d + 9], 12, -1958414417), o = e(o, p, m, n, a[d + 10], 17, -42063), n = e(n, o, p, m, a[d + 11], 22, -1990404162), m = e(m, n, o, p, a[d + 12], 7, 1804603682), p = e(p, m, n, o, a[d + 13], 12, -40341101), o = e(o, p, m, n, a[d + 14], 17, -1502002290), n = e(n, o, p, m, a[d + 15], 22, 1236535329), m = f(m, n, o, p, a[d + 1], 5, -165796510), p = f(p, m, n, o, a[d + 6], 9, -1069501632), o = f(o, p, m, n, a[d + 11], 14, 643717713), n = f(n, o, p, m, a[d], 20, -373897302), m = f(m, n, o, p, a[d + 5], 5, -701558691), p = f(p, m, n, o, a[d + 10], 9, 38016083), o = f(o, p, m, n, a[d + 15], 14, -660478335), n = f(n, o, p, m, a[d + 4], 20, -405537848), m = f(m, n, o, p, a[d + 9], 5, 568446438), p = f(p, m, n, o, a[d + 14], 9, -1019803690), o = f(o, p, m, n, a[d + 3], 14, -187363961), n = f(n, o, p, m, a[d + 8], 20, 1163531501), m = f(m, n, o, p, a[d + 13], 5, -1444681467), p = f(p, m, n, o, a[d + 2], 9, -51403784), o = f(o, p, m, n, a[d + 7], 14, 1735328473), n = f(n, o, p, m, a[d + 12], 20, -1926607734), m = g(m, n, o, p, a[d + 5], 4, -378558), p = g(p, m, n, o, a[d + 8], 11, -2022574463), o = g(o, p, m, n, a[d + 11], 16, 1839030562), n = g(n, o, p, m, a[d + 14], 23, -35309556), m = g(m, n, o, p, a[d + 1], 4, -1530992060), p = g(p, m, n, o, a[d + 4], 11, 1272893353), o = g(o, p, m, n, a[d + 7], 16, -155497632), n = g(n, o, p, m, a[d + 10], 23, -1094730640), m = g(m, n, o, p, a[d + 13], 4, 681279174), p = g(p, m, n, o, a[d], 11, -358537222), o = g(o, p, m, n, a[d + 3], 16, -722521979), n = g(n, o, p, m, a[d + 6], 23, 76029189), m = g(m, n, o, p, a[d + 9], 4, -640364487), p = g(p, m, n, o, a[d + 12], 11, -421815835), o = g(o, p, m, n, a[d + 15], 16, 530742520), n = g(n, o, p, m, a[d + 2], 23, -995338651), m = h(m, n, o, p, a[d], 6, -198630844), p = h(p, m, n, o, a[d + 7], 10, 1126891415), o = h(o, p, m, n, a[d + 14], 15, -1416354905), n = h(n, o, p, m, a[d + 5], 21, -57434055), m = h(m, n, o, p, a[d + 12], 6, 1700485571), p = h(p, m, n, o, a[d + 3], 10, -1894986606), o = h(o, p, m, n, a[d + 10], 15, -1051523), n = h(n, o, p, m, a[d + 1], 21, -2054922799), m = h(m, n, o, p, a[d + 8], 6, 1873313359), p = h(p, m, n, o, a[d + 15], 10, -30611744), o = h(o, p, m, n, a[d + 6], 15, -1560198380), n = h(n, o, p, m, a[d + 13], 21, 1309151649), m = h(m, n, o, p, a[d + 4], 6, -145523070), p = h(p, m, n, o, a[d + 11], 10, -1120210379), o = h(o, p, m, n, a[d + 2], 15, 718787259), n = h(n, o, p, m, a[d + 9], 21, -343485551), m = b(m, i), n = b(n, j), o = b(o, k), p = b(p, l);\n" +
		"        return [m, n, o, p]\n" +
		"    }\n" +
		"\n" +
		"    function j(a) {\n" +
		"        var b, c = \"\";\n" +
		"        for (b = 0; b < 32 * a.length; b += 8) c += String.fromCharCode(a[b >> 5] >>> b % 32 & 255);\n" +
		"        return c\n" +
		"    }\n" +
		"\n" +
		"    function k(a) {\n" +
		"        var b, c = [];\n" +
		"        for (c[(a.length >> 2) - 1] = void 0, b = 0; b < c.length; b += 1) c[b] = 0;\n" +
		"        for (b = 0; b < 8 * a.length; b += 8) c[b >> 5] |= (255 & a.charCodeAt(b / 8)) << b % 32;\n" +
		"        return c\n" +
		"    }\n" +
		"\n" +
		"    function l(a) {\n" +
		"        return j(i(k(a), 8 * a.length))\n" +
		"    }\n" +
		"\n" +
		"    function m(a, b) {\n" +
		"        var c, d, e = k(a), f = [], g = [];\n" +
		"        for (f[15] = g[15] = void 0, e.length > 16 && (e = i(e, 8 * a.length)), c = 0; 16 > c; c += 1) f[c] = 909522486 ^ e[c], g[c] = 1549556828 ^ e[c];\n" +
		"        return d = i(f.concat(k(b)), 512 + 8 * b.length), j(i(g.concat(d), 640))\n" +
		"    }\n" +
		"\n" +
		"    function n(a) {\n" +
		"        var b, c, d = \"0123456789abcdef\", e = \"\";\n" +
		"        for (c = 0; c < a.length; c += 1) b = a.charCodeAt(c), e += d.charAt(b >>> 4 & 15) + d.charAt(15 & b);\n" +
		"        return e\n" +
		"    }\n" +
		"\n" +
		"    function o(a) {\n" +
		"        return unescape(encodeURIComponent(a))\n" +
		"    }\n" +
		"\n" +
		"    function p(a) {\n" +
		"        return l(o(a))\n" +
		"    }\n" +
		"\n" +
		"    function q(a) {\n" +
		"        return n(p(a))\n" +
		"    }\n" +
		"\n" +
		"    function r(a, b) {\n" +
		"        return m(o(a), o(b))\n" +
		"    }\n" +
		"\n" +
		"    function s(a, b) {\n" +
		"        return n(r(a, b))\n" +
		"    }\n" +
		"\n" +
		"    function t(a, b, c) {\n" +
		"        return b ? c ? r(b, a) : s(b, a) : c ? p(a) : q(a)\n" +
		"    }\n" +
		"\n" +
		"    \"function\" == typeof define && define.amd ? define(function () {\n" +
		"        return t\n" +
		"    }) : a.md5 = t\n" +
		"}(this);\n" +
		"\n" +
		"/*REPLACE_BEGIN*/\n" +
		"/**********替换部分的代码(CNTV)***********/\n" +
		"var _0xda49 = [\"\x75\x64\x6B\x32\x37\x76\x6E\x38\x6C\x64\x66\x33\x6C\x63\x76\x31\x73\x70\", \"\x70\x61\x72\x73\x65\", \"\x30\", \"\x72\x61\x6E\x64\x6F\x6D\", \"\x66\x6C\x6F\x6F\x72\",\n" +
		"    \"\x2F\x63\x6E\x74\x76\x6C\x69\x76\x65\x2F\", \"\x6D\x64\x2E\x6D\x33\x75\x38\", \"\x2D\",\n" +
		"    \"\x64\x6C\x39\x32\x66\x39\x63\x6A\x68\x33\x68\x38\x32\x76\x63\x36\x32\x6B\x78\x61\x6C\x69\x77\x6C\x31\x66\", \"\x73\x75\x62\x73\x74\x72\x69\x6E\x67\",\n" +
		"    \"\x68\x74\x74\x70\x3A\x2F\x2F\x68\x6C\x73\x32\x2E\x63\x6E\x74\x76\x2E\x6D\x79\x61\x6C\x69\x63\x64\x6E\x2E\x63\x6F\x6D\", \"\x3F\x61\x75\x74\x68\x5F\x6B\x65\x79\x3D\",\n" +
		"    \"\x4C\x49\x56\x45\x2D\x48\x4C\x53\x2D\x43\x44\x4E\x2D\x54\x58\x59\",\n" +
		"    \"\x68\x74\x74\x70\x3A\x2F\x2F\x63\x63\x74\x76\x35\x2E\x74\x78\x74\x79\x2E\x35\x32\x31\x33\x2E\x6C\x69\x76\x65\x70\x6C\x61\x79\x2E\x6D\x79\x71\x63\x6C\x6F\x75\x64\x2E\x63\x6F\x6D\x2F\x6C\x69\x76\x65\x2F\",\n" +
		"    \"\x5F\x74\x78\x74\x79\x2E\x6D\x33\x75\x38\", \"\x31\", \"\x73\x72\x63\", \"\x5F\x68\x74\x6D\x6C\x35\x50\x6C\x61\x79\x65\x72\", \"\x68\x75\x61\x77\x65\x69\",\n" +
		"    \"\x69\x6E\x64\x65\x78\x4F\x66\", \"\x74\x6F\x4C\x6F\x77\x65\x72\x43\x61\x73\x65\", \"\x75\x73\x65\x72\x41\x67\x65\x6E\x74\", \"\x6C\x6F\x61\x64\",\n" +
		"    \"\x64\x69\x61\x6E\x70\x69\x61\x6E\x2E\x6D\x70\x34\", \"\x75\x6E\x64\x65\x66\x69\x6E\x65\x64\", \"\x6C\x65\x6E\x67\x74\x68\", \"\x70\x6F\x73\x74\x65\x72\",\n" +
		"    \"\x4C\x49\x56\x45\x2D\x48\x44\x53\x2D\x43\x44\x4E\x2D\x41\x4C\x49\"];\n" +
		"\n" +
		"function setFlvHtml5AliNewUrl(videoName) {\n" +
		"    var _0x56eax3 = _0xda49[0];\n" +
		"    var _0x56eax4 = Date[_0xda49[1]](new Date()) / 1000;\n" +
		"    var _0x56eax5 = _0xda49[2];\n" +
		"    var _0x56eax6 = Math[_0xda49[4]](Math[_0xda49[3]]() * 1000);\n" +
		"    var _0x56eax7 = _0xda49[5] + videoName + _0xda49[6];\n" +
		"    var _0x1576 = [\"\x63\x75\x6C\x33\", \"\x63\x71\x6C\x36\"];\n" +
		"    var html5Aauth = _0x1576[0];\n" +
		"    _0x1576[1] = \"\x63\x67\x6C\x35\";\n" +
		"    var html5CdnStr = _0x1576[1] + html5Aauth;\n" +
		"    var _0x56eax8 = md5(_0x56eax7 + _0xda49[7] + _0x56eax4 + _0xda49[7] + _0x56eax6 + _0xda49[7] + _0x56eax5 + _0xda49[7] + _0xda49[8][19] + html5Aauth + _0x56eax3[_0xda49[9]](5, 11) + html5CdnStr);\n" +
		"    var _0x56eax9 = _0x56eax4 + _0xda49[7] + _0x56eax6 + _0xda49[7] + _0x56eax5 + _0xda49[7] + _0x56eax8;\n" +
		"    var _0x56eaxa = _0xda49[10] + _0x56eax7 + _0xda49[11] + _0x56eax9;\n" +
		"    var _0x56eaxb = _0x56eaxa;\n" +
		"    return _0x56eaxb;\n" +
		"    // if (isHtml5Tengxun(videoName)) {\n" +
		"    //     cdnName = _0xda49[12];\n" +
		"    //     _0x56eaxa = _0xda49[13] + videoName + _0xda49[14]\n" +
		"    // };\n" +
		"}\n" +
		"\n" +
		"/*REPLACE_END*/\n" +
		"\n" +
		"\n" +
		"/*\n" +
		"CryptoJS v3.1.2\n" +
		"code.google.com/p/crypto-js\n" +
		"(c) 2009-2013 by Jeff Mott. All rights reserved.\n" +
		"code.google.com/p/crypto-js/wiki/License\n" +
		"*/\n" +
		"var CryptoJS = CryptoJS || function (u, l) {\n" +
		"    var d = {}, n = d.lib = {}, p = function () {\n" +
		"        }, s = n.Base = {\n" +
		"            extend: function (a) {\n" +
		"                p.prototype = this;\n" +
		"                var c = new p;\n" +
		"                a && c.mixIn(a);\n" +
		"                c.hasOwnProperty(\"init\") || (c.init = function () {\n" +
		"                    c.$super.init.apply(this, arguments)\n" +
		"                });\n" +
		"                c.init.prototype = c;\n" +
		"                c.$super = this;\n" +
		"                return c\n" +
		"            }, create: function () {\n" +
		"                var a = this.extend();\n" +
		"                a.init.apply(a, arguments);\n" +
		"                return a\n" +
		"            }, init: function () {\n" +
		"            }, mixIn: function (a) {\n" +
		"                for (var c in a) a.hasOwnProperty(c) && (this[c] = a[c]);\n" +
		"                a.hasOwnProperty(\"toString\") && (this.toString = a.toString)\n" +
		"            }, clone: function () {\n" +
		"                return this.init.prototype.extend(this)\n" +
		"            }\n" +
		"        },\n" +
		"        q = n.WordArray = s.extend({\n" +
		"            init: function (a, c) {\n" +
		"                a = this.words = a || [];\n" +
		"                this.sigBytes = c != l ? c : 4 * a.length\n" +
		"            }, toString: function (a) {\n" +
		"                return (a || v).stringify(this)\n" +
		"            }, concat: function (a) {\n" +
		"                var c = this.words, m = a.words, f = this.sigBytes;\n" +
		"                a = a.sigBytes;\n" +
		"                this.clamp();\n" +
		"                if (f % 4) for (var t = 0; t < a; t++) c[f + t >>> 2] |= (m[t >>> 2] >>> 24 - 8 * (t % 4) & 255) << 24 - 8 * ((f + t) % 4); else if (65535 < m.length) for (t = 0; t < a; t += 4) c[f + t >>> 2] = m[t >>> 2]; else c.push.apply(c, m);\n" +
		"                this.sigBytes += a;\n" +
		"                return this\n" +
		"            }, clamp: function () {\n" +
		"                var a = this.words, c = this.sigBytes;\n" +
		"                a[c >>> 2] &= 4294967295 <<\n" +
		"                    32 - 8 * (c % 4);\n" +
		"                a.length = u.ceil(c / 4)\n" +
		"            }, clone: function () {\n" +
		"                var a = s.clone.call(this);\n" +
		"                a.words = this.words.slice(0);\n" +
		"                return a\n" +
		"            }, random: function (a) {\n" +
		"                for (var c = [], m = 0; m < a; m += 4) c.push(4294967296 * u.random() | 0);\n" +
		"                return new q.init(c, a)\n" +
		"            }\n" +
		"        }), w = d.enc = {}, v = w.Hex = {\n" +
		"            stringify: function (a) {\n" +
		"                var c = a.words;\n" +
		"                a = a.sigBytes;\n" +
		"                for (var m = [], f = 0; f < a; f++) {\n" +
		"                    var t = c[f >>> 2] >>> 24 - 8 * (f % 4) & 255;\n" +
		"                    m.push((t >>> 4).toString(16));\n" +
		"                    m.push((t & 15).toString(16))\n" +
		"                }\n" +
		"                return m.join(\"\")\n" +
		"            }, parse: function (a) {\n" +
		"                for (var c = a.length, m = [], f = 0; f < c; f += 2) m[f >>> 3] |= parseInt(a.substr(f,\n" +
		"                    2), 16) << 24 - 4 * (f % 8);\n" +
		"                return new q.init(m, c / 2)\n" +
		"            }\n" +
		"        }, b = w.Latin1 = {\n" +
		"            stringify: function (a) {\n" +
		"                var c = a.words;\n" +
		"                a = a.sigBytes;\n" +
		"                for (var m = [], f = 0; f < a; f++) m.push(String.fromCharCode(c[f >>> 2] >>> 24 - 8 * (f % 4) & 255));\n" +
		"                return m.join(\"\")\n" +
		"            }, parse: function (a) {\n" +
		"                for (var c = a.length, m = [], f = 0; f < c; f++) m[f >>> 2] |= (a.charCodeAt(f) & 255) << 24 - 8 * (f % 4);\n" +
		"                return new q.init(m, c)\n" +
		"            }\n" +
		"        }, x = w.Utf8 = {\n" +
		"            stringify: function (a) {\n" +
		"                try {\n" +
		"                    return decodeURIComponent(escape(b.stringify(a)))\n" +
		"                } catch (c) {\n" +
		"                    throw Error(\"Malformed UTF-8 data\");\n" +
		"                }\n" +
		"            }, parse: function (a) {\n" +
		"                return b.parse(unescape(encodeURIComponent(a)))\n" +
		"            }\n" +
		"        },\n" +
		"        r = n.BufferedBlockAlgorithm = s.extend({\n" +
		"            reset: function () {\n" +
		"                this._data = new q.init;\n" +
		"                this._nDataBytes = 0\n" +
		"            }, _append: function (a) {\n" +
		"                \"string\" == typeof a && (a = x.parse(a));\n" +
		"                this._data.concat(a);\n" +
		"                this._nDataBytes += a.sigBytes\n" +
		"            }, _process: function (a) {\n" +
		"                var c = this._data, m = c.words, f = c.sigBytes, t = this.blockSize, b = f / (4 * t),\n" +
		"                    b = a ? u.ceil(b) : u.max((b | 0) - this._minBufferSize, 0);\n" +
		"                a = b * t;\n" +
		"                f = u.min(4 * a, f);\n" +
		"                if (a) {\n" +
		"                    for (var e = 0; e < a; e += t) this._doProcessBlock(m, e);\n" +
		"                    e = m.splice(0, a);\n" +
		"                    c.sigBytes -= f\n" +
		"                }\n" +
		"                return new q.init(e, f)\n" +
		"            }, clone: function () {\n" +
		"                var a = s.clone.call(this);\n" +
		"                a._data = this._data.clone();\n" +
		"                return a\n" +
		"            }, _minBufferSize: 0\n" +
		"        });\n" +
		"    n.Hasher = r.extend({\n" +
		"        cfg: s.extend(), init: function (a) {\n" +
		"            this.cfg = this.cfg.extend(a);\n" +
		"            this.reset()\n" +
		"        }, reset: function () {\n" +
		"            r.reset.call(this);\n" +
		"            this._doReset()\n" +
		"        }, update: function (a) {\n" +
		"            this._append(a);\n" +
		"            this._process();\n" +
		"            return this\n" +
		"        }, finalize: function (a) {\n" +
		"            a && this._append(a);\n" +
		"            return this._doFinalize()\n" +
		"        }, blockSize: 16, _createHelper: function (a) {\n" +
		"            return function (c, m) {\n" +
		"                return (new a.init(m)).finalize(c)\n" +
		"            }\n" +
		"        }, _createHmacHelper: function (a) {\n" +
		"            return function (c, m) {\n" +
		"                return (new e.HMAC.init(a,\n" +
		"                    m)).finalize(c)\n" +
		"            }\n" +
		"        }\n" +
		"    });\n" +
		"    var e = d.algo = {};\n" +
		"    return d\n" +
		"}(Math);\n" +
		"(function () {\n" +
		"    var u = CryptoJS, l = u.lib.WordArray;\n" +
		"    u.enc.Base64 = {\n" +
		"        stringify: function (d) {\n" +
		"            var n = d.words, l = d.sigBytes, s = this._map;\n" +
		"            d.clamp();\n" +
		"            d = [];\n" +
		"            for (var q = 0; q < l; q += 3) for (var w = (n[q >>> 2] >>> 24 - 8 * (q % 4) & 255) << 16 | (n[q + 1 >>> 2] >>> 24 - 8 * ((q + 1) % 4) & 255) << 8 | n[q + 2 >>> 2] >>> 24 - 8 * ((q + 2) % 4) & 255, v = 0; 4 > v && q + 0.75 * v < l; v++) d.push(s.charAt(w >>> 6 * (3 - v) & 63));\n" +
		"            if (n = s.charAt(64)) for (; d.length % 4;) d.push(n);\n" +
		"            return d.join(\"\")\n" +
		"        }, parse: function (d) {\n" +
		"            var n = d.length, p = this._map, s = p.charAt(64);\n" +
		"            s && (s = d.indexOf(s), -1 != s && (n = s));\n" +
		"            for (var s = [], q = 0, w = 0; w <\n" +
		"            n; w++) if (w % 4) {\n" +
		"                var v = p.indexOf(d.charAt(w - 1)) << 2 * (w % 4), b = p.indexOf(d.charAt(w)) >>> 6 - 2 * (w % 4);\n" +
		"                s[q >>> 2] |= (v | b) << 24 - 8 * (q % 4);\n" +
		"                q++\n" +
		"            }\n" +
		"            return l.create(s, q)\n" +
		"        }, _map: \"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=\"\n" +
		"    }\n" +
		"})();\n" +
		"(function (u) {\n" +
		"    function l(b, e, a, c, m, f, t) {\n" +
		"        b = b + (e & a | ~e & c) + m + t;\n" +
		"        return (b << f | b >>> 32 - f) + e\n" +
		"    }\n" +
		"\n" +
		"    function d(b, e, a, c, m, f, t) {\n" +
		"        b = b + (e & c | a & ~c) + m + t;\n" +
		"        return (b << f | b >>> 32 - f) + e\n" +
		"    }\n" +
		"\n" +
		"    function n(b, e, a, c, m, f, t) {\n" +
		"        b = b + (e ^ a ^ c) + m + t;\n" +
		"        return (b << f | b >>> 32 - f) + e\n" +
		"    }\n" +
		"\n" +
		"    function p(b, e, a, c, m, f, t) {\n" +
		"        b = b + (a ^ (e | ~c)) + m + t;\n" +
		"        return (b << f | b >>> 32 - f) + e\n" +
		"    }\n" +
		"\n" +
		"    for (var s = CryptoJS, q = s.lib, w = q.WordArray, v = q.Hasher, q = s.algo, b = [], x = 0; 64 > x; x++) b[x] = 4294967296 * u.abs(u.sin(x + 1)) | 0;\n" +
		"    q = q.MD5 = v.extend({\n" +
		"        _doReset: function () {\n" +
		"            this._hash = new w.init([1732584193, 4023233417, 2562383102, 271733878])\n" +
		"        },\n" +
		"        _doProcessBlock: function (r, e) {\n" +
		"            for (var a = 0; 16 > a; a++) {\n" +
		"                var c = e + a, m = r[c];\n" +
		"                r[c] = (m << 8 | m >>> 24) & 16711935 | (m << 24 | m >>> 8) & 4278255360\n" +
		"            }\n" +
		"            var a = this._hash.words, c = r[e + 0], m = r[e + 1], f = r[e + 2], t = r[e + 3], y = r[e + 4],\n" +
		"                q = r[e + 5], s = r[e + 6], w = r[e + 7], v = r[e + 8], u = r[e + 9], x = r[e + 10], z = r[e + 11],\n" +
		"                A = r[e + 12], B = r[e + 13], C = r[e + 14], D = r[e + 15], g = a[0], h = a[1], j = a[2], k = a[3],\n" +
		"                g = l(g, h, j, k, c, 7, b[0]), k = l(k, g, h, j, m, 12, b[1]), j = l(j, k, g, h, f, 17, b[2]),\n" +
		"                h = l(h, j, k, g, t, 22, b[3]), g = l(g, h, j, k, y, 7, b[4]), k = l(k, g, h, j, q, 12, b[5]),\n" +
		"                j = l(j, k, g, h, s, 17, b[6]), h = l(h, j, k, g, w, 22, b[7]),\n" +
		"                g = l(g, h, j, k, v, 7, b[8]), k = l(k, g, h, j, u, 12, b[9]), j = l(j, k, g, h, x, 17, b[10]),\n" +
		"                h = l(h, j, k, g, z, 22, b[11]), g = l(g, h, j, k, A, 7, b[12]), k = l(k, g, h, j, B, 12, b[13]),\n" +
		"                j = l(j, k, g, h, C, 17, b[14]), h = l(h, j, k, g, D, 22, b[15]), g = d(g, h, j, k, m, 5, b[16]),\n" +
		"                k = d(k, g, h, j, s, 9, b[17]), j = d(j, k, g, h, z, 14, b[18]), h = d(h, j, k, g, c, 20, b[19]),\n" +
		"                g = d(g, h, j, k, q, 5, b[20]), k = d(k, g, h, j, x, 9, b[21]), j = d(j, k, g, h, D, 14, b[22]),\n" +
		"                h = d(h, j, k, g, y, 20, b[23]), g = d(g, h, j, k, u, 5, b[24]), k = d(k, g, h, j, C, 9, b[25]),\n" +
		"                j = d(j, k, g, h, t, 14, b[26]), h = d(h, j, k, g, v, 20, b[27]), g = d(g, h, j, k, B, 5, b[28]),\n" +
		"                k = d(k, g,\n" +
		"                    h, j, f, 9, b[29]), j = d(j, k, g, h, w, 14, b[30]), h = d(h, j, k, g, A, 20, b[31]),\n" +
		"                g = n(g, h, j, k, q, 4, b[32]), k = n(k, g, h, j, v, 11, b[33]), j = n(j, k, g, h, z, 16, b[34]),\n" +
		"                h = n(h, j, k, g, C, 23, b[35]), g = n(g, h, j, k, m, 4, b[36]), k = n(k, g, h, j, y, 11, b[37]),\n" +
		"                j = n(j, k, g, h, w, 16, b[38]), h = n(h, j, k, g, x, 23, b[39]), g = n(g, h, j, k, B, 4, b[40]),\n" +
		"                k = n(k, g, h, j, c, 11, b[41]), j = n(j, k, g, h, t, 16, b[42]), h = n(h, j, k, g, s, 23, b[43]),\n" +
		"                g = n(g, h, j, k, u, 4, b[44]), k = n(k, g, h, j, A, 11, b[45]), j = n(j, k, g, h, D, 16, b[46]),\n" +
		"                h = n(h, j, k, g, f, 23, b[47]), g = p(g, h, j, k, c, 6, b[48]), k = p(k, g, h, j, w, 10, b[49]),\n" +
		"                j = p(j, k, g, h,\n" +
		"                    C, 15, b[50]), h = p(h, j, k, g, q, 21, b[51]), g = p(g, h, j, k, A, 6, b[52]),\n" +
		"                k = p(k, g, h, j, t, 10, b[53]), j = p(j, k, g, h, x, 15, b[54]), h = p(h, j, k, g, m, 21, b[55]),\n" +
		"                g = p(g, h, j, k, v, 6, b[56]), k = p(k, g, h, j, D, 10, b[57]), j = p(j, k, g, h, s, 15, b[58]),\n" +
		"                h = p(h, j, k, g, B, 21, b[59]), g = p(g, h, j, k, y, 6, b[60]), k = p(k, g, h, j, z, 10, b[61]),\n" +
		"                j = p(j, k, g, h, f, 15, b[62]), h = p(h, j, k, g, u, 21, b[63]);\n" +
		"            a[0] = a[0] + g | 0;\n" +
		"            a[1] = a[1] + h | 0;\n" +
		"            a[2] = a[2] + j | 0;\n" +
		"            a[3] = a[3] + k | 0\n" +
		"        }, _doFinalize: function () {\n" +
		"            var b = this._data, e = b.words, a = 8 * this._nDataBytes, c = 8 * b.sigBytes;\n" +
		"            e[c >>> 5] |= 128 << 24 - c % 32;\n" +
		"            var m = u.floor(a /\n" +
		"                4294967296);\n" +
		"            e[(c + 64 >>> 9 << 4) + 15] = (m << 8 | m >>> 24) & 16711935 | (m << 24 | m >>> 8) & 4278255360;\n" +
		"            e[(c + 64 >>> 9 << 4) + 14] = (a << 8 | a >>> 24) & 16711935 | (a << 24 | a >>> 8) & 4278255360;\n" +
		"            b.sigBytes = 4 * (e.length + 1);\n" +
		"            this._process();\n" +
		"            b = this._hash;\n" +
		"            e = b.words;\n" +
		"            for (a = 0; 4 > a; a++) c = e[a], e[a] = (c << 8 | c >>> 24) & 16711935 | (c << 24 | c >>> 8) & 4278255360;\n" +
		"            return b\n" +
		"        }, clone: function () {\n" +
		"            var b = v.clone.call(this);\n" +
		"            b._hash = this._hash.clone();\n" +
		"            return b\n" +
		"        }\n" +
		"    });\n" +
		"    s.MD5 = v._createHelper(q);\n" +
		"    s.HmacMD5 = v._createHmacHelper(q)\n" +
		"})(Math);\n" +
		"(function () {\n" +
		"    var u = CryptoJS, l = u.lib, d = l.Base, n = l.WordArray, l = u.algo, p = l.EvpKDF = d.extend({\n" +
		"        cfg: d.extend({keySize: 4, hasher: l.MD5, iterations: 1}), init: function (d) {\n" +
		"            this.cfg = this.cfg.extend(d)\n" +
		"        }, compute: function (d, l) {\n" +
		"            for (var p = this.cfg, v = p.hasher.create(), b = n.create(), u = b.words, r = p.keySize, p = p.iterations; u.length < r;) {\n" +
		"                e && v.update(e);\n" +
		"                var e = v.update(d).finalize(l);\n" +
		"                v.reset();\n" +
		"                for (var a = 1; a < p; a++) e = v.finalize(e), v.reset();\n" +
		"                b.concat(e)\n" +
		"            }\n" +
		"            b.sigBytes = 4 * r;\n" +
		"            return b\n" +
		"        }\n" +
		"    });\n" +
		"    u.EvpKDF = function (d, l, n) {\n" +
		"        return p.create(n).compute(d,\n" +
		"            l)\n" +
		"    }\n" +
		"})();\n" +
		"CryptoJS.lib.Cipher || function (u) {\n" +
		"    var l = CryptoJS, d = l.lib, n = d.Base, p = d.WordArray, s = d.BufferedBlockAlgorithm, q = l.enc.Base64,\n" +
		"        w = l.algo.EvpKDF, v = d.Cipher = s.extend({\n" +
		"            cfg: n.extend(), createEncryptor: function (m, a) {\n" +
		"                return this.create(this._ENC_XFORM_MODE, m, a)\n" +
		"            }, createDecryptor: function (m, a) {\n" +
		"                return this.create(this._DEC_XFORM_MODE, m, a)\n" +
		"            }, init: function (m, a, b) {\n" +
		"                this.cfg = this.cfg.extend(b);\n" +
		"                this._xformMode = m;\n" +
		"                this._key = a;\n" +
		"                this.reset()\n" +
		"            }, reset: function () {\n" +
		"                s.reset.call(this);\n" +
		"                this._doReset()\n" +
		"            }, process: function (a) {\n" +
		"                this._append(a);\n" +
		"                return this._process()\n" +
		"            },\n" +
		"            finalize: function (a) {\n" +
		"                a && this._append(a);\n" +
		"                return this._doFinalize()\n" +
		"            }, keySize: 4, ivSize: 4, _ENC_XFORM_MODE: 1, _DEC_XFORM_MODE: 2, _createHelper: function (m) {\n" +
		"                return {\n" +
		"                    encrypt: function (f, b, e) {\n" +
		"                        return (\"string\" == typeof b ? c : a).encrypt(m, f, b, e)\n" +
		"                    }, decrypt: function (f, b, e) {\n" +
		"                        return (\"string\" == typeof b ? c : a).decrypt(m, f, b, e)\n" +
		"                    }\n" +
		"                }\n" +
		"            }\n" +
		"        });\n" +
		"    d.StreamCipher = v.extend({\n" +
		"        _doFinalize: function () {\n" +
		"            return this._process(!0)\n" +
		"        }, blockSize: 1\n" +
		"    });\n" +
		"    var b = l.mode = {}, x = function (a, f, b) {\n" +
		"        var c = this._iv;\n" +
		"        c ? this._iv = u : c = this._prevBlock;\n" +
		"        for (var e = 0; e < b; e++) a[f + e] ^=\n" +
		"            c[e]\n" +
		"    }, r = (d.BlockCipherMode = n.extend({\n" +
		"        createEncryptor: function (a, f) {\n" +
		"            return this.Encryptor.create(a, f)\n" +
		"        }, createDecryptor: function (a, f) {\n" +
		"            return this.Decryptor.create(a, f)\n" +
		"        }, init: function (a, f) {\n" +
		"            this._cipher = a;\n" +
		"            this._iv = f\n" +
		"        }\n" +
		"    })).extend();\n" +
		"    r.Encryptor = r.extend({\n" +
		"        processBlock: function (a, f) {\n" +
		"            var b = this._cipher, c = b.blockSize;\n" +
		"            x.call(this, a, f, c);\n" +
		"            b.encryptBlock(a, f);\n" +
		"            this._prevBlock = a.slice(f, f + c)\n" +
		"        }\n" +
		"    });\n" +
		"    r.Decryptor = r.extend({\n" +
		"        processBlock: function (a, b) {\n" +
		"            var c = this._cipher, e = c.blockSize, d = a.slice(b, b + e);\n" +
		"            c.decryptBlock(a, b);\n" +
		"            x.call(this,\n" +
		"                a, b, e);\n" +
		"            this._prevBlock = d\n" +
		"        }\n" +
		"    });\n" +
		"    b = b.CBC = r;\n" +
		"    r = (l.pad = {}).Pkcs7 = {\n" +
		"        pad: function (a, b) {\n" +
		"            for (var c = 4 * b, c = c - a.sigBytes % c, e = c << 24 | c << 16 | c << 8 | c, d = [], l = 0; l < c; l += 4) d.push(e);\n" +
		"            c = p.create(d, c);\n" +
		"            a.concat(c)\n" +
		"        }, unpad: function (a) {\n" +
		"            a.sigBytes -= a.words[a.sigBytes - 1 >>> 2] & 255\n" +
		"        }\n" +
		"    };\n" +
		"    d.BlockCipher = v.extend({\n" +
		"        cfg: v.cfg.extend({mode: b, padding: r}), reset: function () {\n" +
		"            v.reset.call(this);\n" +
		"            var a = this.cfg, c = a.iv, a = a.mode;\n" +
		"            if (this._xformMode == this._ENC_XFORM_MODE) var b = a.createEncryptor; else b = a.createDecryptor, this._minBufferSize = 1;\n" +
		"            this._mode = b.call(a,\n" +
		"                this, c && c.words)\n" +
		"        }, _doProcessBlock: function (a, c) {\n" +
		"            this._mode.processBlock(a, c)\n" +
		"        }, _doFinalize: function () {\n" +
		"            var a = this.cfg.padding;\n" +
		"            if (this._xformMode == this._ENC_XFORM_MODE) {\n" +
		"                a.pad(this._data, this.blockSize);\n" +
		"                var c = this._process(!0)\n" +
		"            } else c = this._process(!0), a.unpad(c);\n" +
		"            return c\n" +
		"        }, blockSize: 4\n" +
		"    });\n" +
		"    var e = d.CipherParams = n.extend({\n" +
		"        init: function (a) {\n" +
		"            this.mixIn(a)\n" +
		"        }, toString: function (a) {\n" +
		"            return (a || this.formatter).stringify(this)\n" +
		"        }\n" +
		"    }), b = (l.format = {}).OpenSSL = {\n" +
		"        stringify: function (a) {\n" +
		"            var c = a.ciphertext;\n" +
		"            a = a.salt;\n" +
		"            return (a ? p.create([1398893684,\n" +
		"                1701076831]).concat(a).concat(c) : c).toString(q)\n" +
		"        }, parse: function (a) {\n" +
		"            a = q.parse(a);\n" +
		"            var c = a.words;\n" +
		"            if (1398893684 == c[0] && 1701076831 == c[1]) {\n" +
		"                var b = p.create(c.slice(2, 4));\n" +
		"                c.splice(0, 4);\n" +
		"                a.sigBytes -= 16\n" +
		"            }\n" +
		"            return e.create({ciphertext: a, salt: b})\n" +
		"        }\n" +
		"    }, a = d.SerializableCipher = n.extend({\n" +
		"        cfg: n.extend({format: b}), encrypt: function (a, c, b, d) {\n" +
		"            d = this.cfg.extend(d);\n" +
		"            var l = a.createEncryptor(b, d);\n" +
		"            c = l.finalize(c);\n" +
		"            l = l.cfg;\n" +
		"            return e.create({\n" +
		"                ciphertext: c,\n" +
		"                key: b,\n" +
		"                iv: l.iv,\n" +
		"                algorithm: a,\n" +
		"                mode: l.mode,\n" +
		"                padding: l.padding,\n" +
		"                blockSize: a.blockSize,\n" +
		"                formatter: d.format\n" +
		"            })\n" +
		"        },\n" +
		"        decrypt: function (a, c, b, e) {\n" +
		"            e = this.cfg.extend(e);\n" +
		"            c = this._parse(c, e.format);\n" +
		"            return a.createDecryptor(b, e).finalize(c.ciphertext)\n" +
		"        }, _parse: function (a, c) {\n" +
		"            return \"string\" == typeof a ? c.parse(a, this) : a\n" +
		"        }\n" +
		"    }), l = (l.kdf = {}).OpenSSL = {\n" +
		"        execute: function (a, c, b, d) {\n" +
		"            d || (d = p.random(8));\n" +
		"            a = w.create({keySize: c + b}).compute(a, d);\n" +
		"            b = p.create(a.words.slice(c), 4 * b);\n" +
		"            a.sigBytes = 4 * c;\n" +
		"            return e.create({key: a, iv: b, salt: d})\n" +
		"        }\n" +
		"    }, c = d.PasswordBasedCipher = a.extend({\n" +
		"        cfg: a.cfg.extend({kdf: l}), encrypt: function (c, b, e, d) {\n" +
		"            d = this.cfg.extend(d);\n" +
		"            e = d.kdf.execute(e,\n" +
		"                c.keySize, c.ivSize);\n" +
		"            d.iv = e.iv;\n" +
		"            c = a.encrypt.call(this, c, b, e.key, d);\n" +
		"            c.mixIn(e);\n" +
		"            return c\n" +
		"        }, decrypt: function (c, b, e, d) {\n" +
		"            d = this.cfg.extend(d);\n" +
		"            b = this._parse(b, d.format);\n" +
		"            e = d.kdf.execute(e, c.keySize, c.ivSize, b.salt);\n" +
		"            d.iv = e.iv;\n" +
		"            return a.decrypt.call(this, c, b, e.key, d)\n" +
		"        }\n" +
		"    })\n" +
		"}();\n" +
		"(function () {\n" +
		"    function u(b, a) {\n" +
		"        var c = (this._lBlock >>> b ^ this._rBlock) & a;\n" +
		"        this._rBlock ^= c;\n" +
		"        this._lBlock ^= c << b\n" +
		"    }\n" +
		"\n" +
		"    function l(b, a) {\n" +
		"        var c = (this._rBlock >>> b ^ this._lBlock) & a;\n" +
		"        this._lBlock ^= c;\n" +
		"        this._rBlock ^= c << b\n" +
		"    }\n" +
		"\n" +
		"    var d = CryptoJS, n = d.lib, p = n.WordArray, n = n.BlockCipher, s = d.algo,\n" +
		"        q = [57, 49, 41, 33, 25, 17, 9, 1, 58, 50, 42, 34, 26, 18, 10, 2, 59, 51, 43, 35, 27, 19, 11, 3, 60, 52, 44, 36, 63, 55, 47, 39, 31, 23, 15, 7, 62, 54, 46, 38, 30, 22, 14, 6, 61, 53, 45, 37, 29, 21, 13, 5, 28, 20, 12, 4],\n" +
		"        w = [14, 17, 11, 24, 1, 5, 3, 28, 15, 6, 21, 10, 23, 19, 12, 4, 26, 8, 16, 7, 27, 20, 13, 2, 41, 52, 31, 37, 47,\n" +
		"            55, 30, 40, 51, 45, 33, 48, 44, 49, 39, 56, 34, 53, 46, 42, 50, 36, 29, 32],\n" +
		"        v = [1, 2, 4, 6, 8, 10, 12, 14, 15, 17, 19, 21, 23, 25, 27, 28], b = [{\n" +
		"            \"0\": 8421888,\n" +
		"            268435456: 32768,\n" +
		"            536870912: 8421378,\n" +
		"            805306368: 2,\n" +
		"            1073741824: 512,\n" +
		"            1342177280: 8421890,\n" +
		"            1610612736: 8389122,\n" +
		"            1879048192: 8388608,\n" +
		"            2147483648: 514,\n" +
		"            2415919104: 8389120,\n" +
		"            2684354560: 33280,\n" +
		"            2952790016: 8421376,\n" +
		"            3221225472: 32770,\n" +
		"            3489660928: 8388610,\n" +
		"            3758096384: 0,\n" +
		"            4026531840: 33282,\n" +
		"            134217728: 0,\n" +
		"            402653184: 8421890,\n" +
		"            671088640: 33282,\n" +
		"            939524096: 32768,\n" +
		"            1207959552: 8421888,\n" +
		"            1476395008: 512,\n" +
		"            1744830464: 8421378,\n" +
		"            2013265920: 2,\n" +
		"            2281701376: 8389120,\n" +
		"            2550136832: 33280,\n" +
		"            2818572288: 8421376,\n" +
		"            3087007744: 8389122,\n" +
		"            3355443200: 8388610,\n" +
		"            3623878656: 32770,\n" +
		"            3892314112: 514,\n" +
		"            4160749568: 8388608,\n" +
		"            1: 32768,\n" +
		"            268435457: 2,\n" +
		"            536870913: 8421888,\n" +
		"            805306369: 8388608,\n" +
		"            1073741825: 8421378,\n" +
		"            1342177281: 33280,\n" +
		"            1610612737: 512,\n" +
		"            1879048193: 8389122,\n" +
		"            2147483649: 8421890,\n" +
		"            2415919105: 8421376,\n" +
		"            2684354561: 8388610,\n" +
		"            2952790017: 33282,\n" +
		"            3221225473: 514,\n" +
		"            3489660929: 8389120,\n" +
		"            3758096385: 32770,\n" +
		"            4026531841: 0,\n" +
		"            134217729: 8421890,\n" +
		"            402653185: 8421376,\n" +
		"            671088641: 8388608,\n" +
		"            939524097: 512,\n" +
		"            1207959553: 32768,\n" +
		"            1476395009: 8388610,\n" +
		"            1744830465: 2,\n" +
		"            2013265921: 33282,\n" +
		"            2281701377: 32770,\n" +
		"            2550136833: 8389122,\n" +
		"            2818572289: 514,\n" +
		"            3087007745: 8421888,\n" +
		"            3355443201: 8389120,\n" +
		"            3623878657: 0,\n" +
		"            3892314113: 33280,\n" +
		"            4160749569: 8421378\n" +
		"        }, {\n" +
		"            \"0\": 1074282512,\n" +
		"            16777216: 16384,\n" +
		"            33554432: 524288,\n" +
		"            50331648: 1074266128,\n" +
		"            67108864: 1073741840,\n" +
		"            83886080: 1074282496,\n" +
		"            100663296: 1073758208,\n" +
		"            117440512: 16,\n" +
		"            134217728: 540672,\n" +
		"            150994944: 1073758224,\n" +
		"            167772160: 1073741824,\n" +
		"            184549376: 540688,\n" +
		"            201326592: 524304,\n" +
		"            218103808: 0,\n" +
		"            234881024: 16400,\n" +
		"            251658240: 1074266112,\n" +
		"            8388608: 1073758208,\n" +
		"            25165824: 540688,\n" +
		"            41943040: 16,\n" +
		"            58720256: 1073758224,\n" +
		"            75497472: 1074282512,\n" +
		"            92274688: 1073741824,\n" +
		"            109051904: 524288,\n" +
		"            125829120: 1074266128,\n" +
		"            142606336: 524304,\n" +
		"            159383552: 0,\n" +
		"            176160768: 16384,\n" +
		"            192937984: 1074266112,\n" +
		"            209715200: 1073741840,\n" +
		"            226492416: 540672,\n" +
		"            243269632: 1074282496,\n" +
		"            260046848: 16400,\n" +
		"            268435456: 0,\n" +
		"            285212672: 1074266128,\n" +
		"            301989888: 1073758224,\n" +
		"            318767104: 1074282496,\n" +
		"            335544320: 1074266112,\n" +
		"            352321536: 16,\n" +
		"            369098752: 540688,\n" +
		"            385875968: 16384,\n" +
		"            402653184: 16400,\n" +
		"            419430400: 524288,\n" +
		"            436207616: 524304,\n" +
		"            452984832: 1073741840,\n" +
		"            469762048: 540672,\n" +
		"            486539264: 1073758208,\n" +
		"            503316480: 1073741824,\n" +
		"            520093696: 1074282512,\n" +
		"            276824064: 540688,\n" +
		"            293601280: 524288,\n" +
		"            310378496: 1074266112,\n" +
		"            327155712: 16384,\n" +
		"            343932928: 1073758208,\n" +
		"            360710144: 1074282512,\n" +
		"            377487360: 16,\n" +
		"            394264576: 1073741824,\n" +
		"            411041792: 1074282496,\n" +
		"            427819008: 1073741840,\n" +
		"            444596224: 1073758224,\n" +
		"            461373440: 524304,\n" +
		"            478150656: 0,\n" +
		"            494927872: 16400,\n" +
		"            511705088: 1074266128,\n" +
		"            528482304: 540672\n" +
		"        }, {\n" +
		"            \"0\": 260,\n" +
		"            1048576: 0,\n" +
		"            2097152: 67109120,\n" +
		"            3145728: 65796,\n" +
		"            4194304: 65540,\n" +
		"            5242880: 67108868,\n" +
		"            6291456: 67174660,\n" +
		"            7340032: 67174400,\n" +
		"            8388608: 67108864,\n" +
		"            9437184: 67174656,\n" +
		"            10485760: 65792,\n" +
		"            11534336: 67174404,\n" +
		"            12582912: 67109124,\n" +
		"            13631488: 65536,\n" +
		"            14680064: 4,\n" +
		"            15728640: 256,\n" +
		"            524288: 67174656,\n" +
		"            1572864: 67174404,\n" +
		"            2621440: 0,\n" +
		"            3670016: 67109120,\n" +
		"            4718592: 67108868,\n" +
		"            5767168: 65536,\n" +
		"            6815744: 65540,\n" +
		"            7864320: 260,\n" +
		"            8912896: 4,\n" +
		"            9961472: 256,\n" +
		"            11010048: 67174400,\n" +
		"            12058624: 65796,\n" +
		"            13107200: 65792,\n" +
		"            14155776: 67109124,\n" +
		"            15204352: 67174660,\n" +
		"            16252928: 67108864,\n" +
		"            16777216: 67174656,\n" +
		"            17825792: 65540,\n" +
		"            18874368: 65536,\n" +
		"            19922944: 67109120,\n" +
		"            20971520: 256,\n" +
		"            22020096: 67174660,\n" +
		"            23068672: 67108868,\n" +
		"            24117248: 0,\n" +
		"            25165824: 67109124,\n" +
		"            26214400: 67108864,\n" +
		"            27262976: 4,\n" +
		"            28311552: 65792,\n" +
		"            29360128: 67174400,\n" +
		"            30408704: 260,\n" +
		"            31457280: 65796,\n" +
		"            32505856: 67174404,\n" +
		"            17301504: 67108864,\n" +
		"            18350080: 260,\n" +
		"            19398656: 67174656,\n" +
		"            20447232: 0,\n" +
		"            21495808: 65540,\n" +
		"            22544384: 67109120,\n" +
		"            23592960: 256,\n" +
		"            24641536: 67174404,\n" +
		"            25690112: 65536,\n" +
		"            26738688: 67174660,\n" +
		"            27787264: 65796,\n" +
		"            28835840: 67108868,\n" +
		"            29884416: 67109124,\n" +
		"            30932992: 67174400,\n" +
		"            31981568: 4,\n" +
		"            33030144: 65792\n" +
		"        }, {\n" +
		"            \"0\": 2151682048,\n" +
		"            65536: 2147487808,\n" +
		"            131072: 4198464,\n" +
		"            196608: 2151677952,\n" +
		"            262144: 0,\n" +
		"            327680: 4198400,\n" +
		"            393216: 2147483712,\n" +
		"            458752: 4194368,\n" +
		"            524288: 2147483648,\n" +
		"            589824: 4194304,\n" +
		"            655360: 64,\n" +
		"            720896: 2147487744,\n" +
		"            786432: 2151678016,\n" +
		"            851968: 4160,\n" +
		"            917504: 4096,\n" +
		"            983040: 2151682112,\n" +
		"            32768: 2147487808,\n" +
		"            98304: 64,\n" +
		"            163840: 2151678016,\n" +
		"            229376: 2147487744,\n" +
		"            294912: 4198400,\n" +
		"            360448: 2151682112,\n" +
		"            425984: 0,\n" +
		"            491520: 2151677952,\n" +
		"            557056: 4096,\n" +
		"            622592: 2151682048,\n" +
		"            688128: 4194304,\n" +
		"            753664: 4160,\n" +
		"            819200: 2147483648,\n" +
		"            884736: 4194368,\n" +
		"            950272: 4198464,\n" +
		"            1015808: 2147483712,\n" +
		"            1048576: 4194368,\n" +
		"            1114112: 4198400,\n" +
		"            1179648: 2147483712,\n" +
		"            1245184: 0,\n" +
		"            1310720: 4160,\n" +
		"            1376256: 2151678016,\n" +
		"            1441792: 2151682048,\n" +
		"            1507328: 2147487808,\n" +
		"            1572864: 2151682112,\n" +
		"            1638400: 2147483648,\n" +
		"            1703936: 2151677952,\n" +
		"            1769472: 4198464,\n" +
		"            1835008: 2147487744,\n" +
		"            1900544: 4194304,\n" +
		"            1966080: 64,\n" +
		"            2031616: 4096,\n" +
		"            1081344: 2151677952,\n" +
		"            1146880: 2151682112,\n" +
		"            1212416: 0,\n" +
		"            1277952: 4198400,\n" +
		"            1343488: 4194368,\n" +
		"            1409024: 2147483648,\n" +
		"            1474560: 2147487808,\n" +
		"            1540096: 64,\n" +
		"            1605632: 2147483712,\n" +
		"            1671168: 4096,\n" +
		"            1736704: 2147487744,\n" +
		"            1802240: 2151678016,\n" +
		"            1867776: 4160,\n" +
		"            1933312: 2151682048,\n" +
		"            1998848: 4194304,\n" +
		"            2064384: 4198464\n" +
		"        }, {\n" +
		"            \"0\": 128,\n" +
		"            4096: 17039360,\n" +
		"            8192: 262144,\n" +
		"            12288: 536870912,\n" +
		"            16384: 537133184,\n" +
		"            20480: 16777344,\n" +
		"            24576: 553648256,\n" +
		"            28672: 262272,\n" +
		"            32768: 16777216,\n" +
		"            36864: 537133056,\n" +
		"            40960: 536871040,\n" +
		"            45056: 553910400,\n" +
		"            49152: 553910272,\n" +
		"            53248: 0,\n" +
		"            57344: 17039488,\n" +
		"            61440: 553648128,\n" +
		"            2048: 17039488,\n" +
		"            6144: 553648256,\n" +
		"            10240: 128,\n" +
		"            14336: 17039360,\n" +
		"            18432: 262144,\n" +
		"            22528: 537133184,\n" +
		"            26624: 553910272,\n" +
		"            30720: 536870912,\n" +
		"            34816: 537133056,\n" +
		"            38912: 0,\n" +
		"            43008: 553910400,\n" +
		"            47104: 16777344,\n" +
		"            51200: 536871040,\n" +
		"            55296: 553648128,\n" +
		"            59392: 16777216,\n" +
		"            63488: 262272,\n" +
		"            65536: 262144,\n" +
		"            69632: 128,\n" +
		"            73728: 536870912,\n" +
		"            77824: 553648256,\n" +
		"            81920: 16777344,\n" +
		"            86016: 553910272,\n" +
		"            90112: 537133184,\n" +
		"            94208: 16777216,\n" +
		"            98304: 553910400,\n" +
		"            102400: 553648128,\n" +
		"            106496: 17039360,\n" +
		"            110592: 537133056,\n" +
		"            114688: 262272,\n" +
		"            118784: 536871040,\n" +
		"            122880: 0,\n" +
		"            126976: 17039488,\n" +
		"            67584: 553648256,\n" +
		"            71680: 16777216,\n" +
		"            75776: 17039360,\n" +
		"            79872: 537133184,\n" +
		"            83968: 536870912,\n" +
		"            88064: 17039488,\n" +
		"            92160: 128,\n" +
		"            96256: 553910272,\n" +
		"            100352: 262272,\n" +
		"            104448: 553910400,\n" +
		"            108544: 0,\n" +
		"            112640: 553648128,\n" +
		"            116736: 16777344,\n" +
		"            120832: 262144,\n" +
		"            124928: 537133056,\n" +
		"            129024: 536871040\n" +
		"        }, {\n" +
		"            \"0\": 268435464,\n" +
		"            256: 8192,\n" +
		"            512: 270532608,\n" +
		"            768: 270540808,\n" +
		"            1024: 268443648,\n" +
		"            1280: 2097152,\n" +
		"            1536: 2097160,\n" +
		"            1792: 268435456,\n" +
		"            2048: 0,\n" +
		"            2304: 268443656,\n" +
		"            2560: 2105344,\n" +
		"            2816: 8,\n" +
		"            3072: 270532616,\n" +
		"            3328: 2105352,\n" +
		"            3584: 8200,\n" +
		"            3840: 270540800,\n" +
		"            128: 270532608,\n" +
		"            384: 270540808,\n" +
		"            640: 8,\n" +
		"            896: 2097152,\n" +
		"            1152: 2105352,\n" +
		"            1408: 268435464,\n" +
		"            1664: 268443648,\n" +
		"            1920: 8200,\n" +
		"            2176: 2097160,\n" +
		"            2432: 8192,\n" +
		"            2688: 268443656,\n" +
		"            2944: 270532616,\n" +
		"            3200: 0,\n" +
		"            3456: 270540800,\n" +
		"            3712: 2105344,\n" +
		"            3968: 268435456,\n" +
		"            4096: 268443648,\n" +
		"            4352: 270532616,\n" +
		"            4608: 270540808,\n" +
		"            4864: 8200,\n" +
		"            5120: 2097152,\n" +
		"            5376: 268435456,\n" +
		"            5632: 268435464,\n" +
		"            5888: 2105344,\n" +
		"            6144: 2105352,\n" +
		"            6400: 0,\n" +
		"            6656: 8,\n" +
		"            6912: 270532608,\n" +
		"            7168: 8192,\n" +
		"            7424: 268443656,\n" +
		"            7680: 270540800,\n" +
		"            7936: 2097160,\n" +
		"            4224: 8,\n" +
		"            4480: 2105344,\n" +
		"            4736: 2097152,\n" +
		"            4992: 268435464,\n" +
		"            5248: 268443648,\n" +
		"            5504: 8200,\n" +
		"            5760: 270540808,\n" +
		"            6016: 270532608,\n" +
		"            6272: 270540800,\n" +
		"            6528: 270532616,\n" +
		"            6784: 8192,\n" +
		"            7040: 2105352,\n" +
		"            7296: 2097160,\n" +
		"            7552: 0,\n" +
		"            7808: 268435456,\n" +
		"            8064: 268443656\n" +
		"        }, {\n" +
		"            \"0\": 1048576,\n" +
		"            16: 33555457,\n" +
		"            32: 1024,\n" +
		"            48: 1049601,\n" +
		"            64: 34604033,\n" +
		"            80: 0,\n" +
		"            96: 1,\n" +
		"            112: 34603009,\n" +
		"            128: 33555456,\n" +
		"            144: 1048577,\n" +
		"            160: 33554433,\n" +
		"            176: 34604032,\n" +
		"            192: 34603008,\n" +
		"            208: 1025,\n" +
		"            224: 1049600,\n" +
		"            240: 33554432,\n" +
		"            8: 34603009,\n" +
		"            24: 0,\n" +
		"            40: 33555457,\n" +
		"            56: 34604032,\n" +
		"            72: 1048576,\n" +
		"            88: 33554433,\n" +
		"            104: 33554432,\n" +
		"            120: 1025,\n" +
		"            136: 1049601,\n" +
		"            152: 33555456,\n" +
		"            168: 34603008,\n" +
		"            184: 1048577,\n" +
		"            200: 1024,\n" +
		"            216: 34604033,\n" +
		"            232: 1,\n" +
		"            248: 1049600,\n" +
		"            256: 33554432,\n" +
		"            272: 1048576,\n" +
		"            288: 33555457,\n" +
		"            304: 34603009,\n" +
		"            320: 1048577,\n" +
		"            336: 33555456,\n" +
		"            352: 34604032,\n" +
		"            368: 1049601,\n" +
		"            384: 1025,\n" +
		"            400: 34604033,\n" +
		"            416: 1049600,\n" +
		"            432: 1,\n" +
		"            448: 0,\n" +
		"            464: 34603008,\n" +
		"            480: 33554433,\n" +
		"            496: 1024,\n" +
		"            264: 1049600,\n" +
		"            280: 33555457,\n" +
		"            296: 34603009,\n" +
		"            312: 1,\n" +
		"            328: 33554432,\n" +
		"            344: 1048576,\n" +
		"            360: 1025,\n" +
		"            376: 34604032,\n" +
		"            392: 33554433,\n" +
		"            408: 34603008,\n" +
		"            424: 0,\n" +
		"            440: 34604033,\n" +
		"            456: 1049601,\n" +
		"            472: 1024,\n" +
		"            488: 33555456,\n" +
		"            504: 1048577\n" +
		"        }, {\n" +
		"            \"0\": 134219808,\n" +
		"            1: 131072,\n" +
		"            2: 134217728,\n" +
		"            3: 32,\n" +
		"            4: 131104,\n" +
		"            5: 134350880,\n" +
		"            6: 134350848,\n" +
		"            7: 2048,\n" +
		"            8: 134348800,\n" +
		"            9: 134219776,\n" +
		"            10: 133120,\n" +
		"            11: 134348832,\n" +
		"            12: 2080,\n" +
		"            13: 0,\n" +
		"            14: 134217760,\n" +
		"            15: 133152,\n" +
		"            2147483648: 2048,\n" +
		"            2147483649: 134350880,\n" +
		"            2147483650: 134219808,\n" +
		"            2147483651: 134217728,\n" +
		"            2147483652: 134348800,\n" +
		"            2147483653: 133120,\n" +
		"            2147483654: 133152,\n" +
		"            2147483655: 32,\n" +
		"            2147483656: 134217760,\n" +
		"            2147483657: 2080,\n" +
		"            2147483658: 131104,\n" +
		"            2147483659: 134350848,\n" +
		"            2147483660: 0,\n" +
		"            2147483661: 134348832,\n" +
		"            2147483662: 134219776,\n" +
		"            2147483663: 131072,\n" +
		"            16: 133152,\n" +
		"            17: 134350848,\n" +
		"            18: 32,\n" +
		"            19: 2048,\n" +
		"            20: 134219776,\n" +
		"            21: 134217760,\n" +
		"            22: 134348832,\n" +
		"            23: 131072,\n" +
		"            24: 0,\n" +
		"            25: 131104,\n" +
		"            26: 134348800,\n" +
		"            27: 134219808,\n" +
		"            28: 134350880,\n" +
		"            29: 133120,\n" +
		"            30: 2080,\n" +
		"            31: 134217728,\n" +
		"            2147483664: 131072,\n" +
		"            2147483665: 2048,\n" +
		"            2147483666: 134348832,\n" +
		"            2147483667: 133152,\n" +
		"            2147483668: 32,\n" +
		"            2147483669: 134348800,\n" +
		"            2147483670: 134217728,\n" +
		"            2147483671: 134219808,\n" +
		"            2147483672: 134350880,\n" +
		"            2147483673: 134217760,\n" +
		"            2147483674: 134219776,\n" +
		"            2147483675: 0,\n" +
		"            2147483676: 133120,\n" +
		"            2147483677: 2080,\n" +
		"            2147483678: 131104,\n" +
		"            2147483679: 134350848\n" +
		"        }], x = [4160749569, 528482304, 33030144, 2064384, 129024, 8064, 504, 2147483679], r = s.DES = n.extend({\n" +
		"            _doReset: function () {\n" +
		"                for (var b = this._key.words, a = [], c = 0; 56 > c; c++) {\n" +
		"                    var d = q[c] - 1;\n" +
		"                    a[c] = b[d >>> 5] >>> 31 - d % 32 & 1\n" +
		"                }\n" +
		"                b = this._subKeys = [];\n" +
		"                for (d = 0; 16 > d; d++) {\n" +
		"                    for (var f = b[d] = [], l = v[d], c = 0; 24 > c; c++) f[c / 6 | 0] |= a[(w[c] - 1 + l) % 28] << 31 - c % 6, f[4 + (c / 6 | 0)] |= a[28 + (w[c + 24] - 1 + l) % 28] << 31 - c % 6;\n" +
		"                    f[0] = f[0] << 1 | f[0] >>> 31;\n" +
		"                    for (c = 1; 7 > c; c++) f[c] >>>=\n" +
		"                        4 * (c - 1) + 3;\n" +
		"                    f[7] = f[7] << 5 | f[7] >>> 27\n" +
		"                }\n" +
		"                a = this._invSubKeys = [];\n" +
		"                for (c = 0; 16 > c; c++) a[c] = b[15 - c]\n" +
		"            }, encryptBlock: function (b, a) {\n" +
		"                this._doCryptBlock(b, a, this._subKeys)\n" +
		"            }, decryptBlock: function (b, a) {\n" +
		"                this._doCryptBlock(b, a, this._invSubKeys)\n" +
		"            }, _doCryptBlock: function (e, a, c) {\n" +
		"                this._lBlock = e[a];\n" +
		"                this._rBlock = e[a + 1];\n" +
		"                u.call(this, 4, 252645135);\n" +
		"                u.call(this, 16, 65535);\n" +
		"                l.call(this, 2, 858993459);\n" +
		"                l.call(this, 8, 16711935);\n" +
		"                u.call(this, 1, 1431655765);\n" +
		"                for (var d = 0; 16 > d; d++) {\n" +
		"                    for (var f = c[d], n = this._lBlock, p = this._rBlock, q = 0, r = 0; 8 > r; r++) q |= b[r][((p ^\n" +
		"                        f[r]) & x[r]) >>> 0];\n" +
		"                    this._lBlock = p;\n" +
		"                    this._rBlock = n ^ q\n" +
		"                }\n" +
		"                c = this._lBlock;\n" +
		"                this._lBlock = this._rBlock;\n" +
		"                this._rBlock = c;\n" +
		"                u.call(this, 1, 1431655765);\n" +
		"                l.call(this, 8, 16711935);\n" +
		"                l.call(this, 2, 858993459);\n" +
		"                u.call(this, 16, 65535);\n" +
		"                u.call(this, 4, 252645135);\n" +
		"                e[a] = this._lBlock;\n" +
		"                e[a + 1] = this._rBlock\n" +
		"            }, keySize: 2, ivSize: 2, blockSize: 2\n" +
		"        });\n" +
		"    d.DES = n._createHelper(r);\n" +
		"    s = s.TripleDES = n.extend({\n" +
		"        _doReset: function () {\n" +
		"            var b = this._key.words;\n" +
		"            this._des1 = r.createEncryptor(p.create(b.slice(0, 2)));\n" +
		"            this._des2 = r.createEncryptor(p.create(b.slice(2, 4)));\n" +
		"            this._des3 =\n" +
		"                r.createEncryptor(p.create(b.slice(4, 6)))\n" +
		"        }, encryptBlock: function (b, a) {\n" +
		"            this._des1.encryptBlock(b, a);\n" +
		"            this._des2.decryptBlock(b, a);\n" +
		"            this._des3.encryptBlock(b, a)\n" +
		"        }, decryptBlock: function (b, a) {\n" +
		"            this._des3.decryptBlock(b, a);\n" +
		"            this._des2.encryptBlock(b, a);\n" +
		"            this._des1.decryptBlock(b, a)\n" +
		"        }, keySize: 6, ivSize: 2, blockSize: 2\n" +
		"    });\n" +
		"    d.TripleDES = n._createHelper(s)\n" +
		"})();\n" +
		"\n" +
		"/*\n" +
		"CryptoJS v3.1.2\n" +
		"code.google.com/p/crypto-js\n" +
		"(c) 2009-2013 by Jeff Mott. All rights reserved.\n" +
		"code.google.com/p/crypto-js/wiki/License\n" +
		"*/\n" +
		"CryptoJS.mode.ECB = function () {\n" +
		"    var a = CryptoJS.lib.BlockCipherMode.extend();\n" +
		"    a.Encryptor = a.extend({\n" +
		"        processBlock: function (a, b) {\n" +
		"            this._cipher.encryptBlock(a, b)\n" +
		"        }\n" +
		"    });\n" +
		"    a.Decryptor = a.extend({\n" +
		"        processBlock: function (a, b) {\n" +
		"            this._cipher.decryptBlock(a, b)\n" +
		"        }\n" +
		"    });\n" +
		"    return a\n" +
		"}();\n" +
		"\n" +
		"\n" +
		"/**\n" +
		" * @return {string}\n" +
		" */\n" +
		"function GetMiguRealPlayUrl(request) {\n" +
		"    var json_data;\n" +
		"    var fetch_url;\n" +
		"    var result;\n" +
		"    var begin = 0;\n" +
		"\n" +
		"    json_data = JSON.parse(request);\n" +
		"\n" +
		"    var agent = \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36\";\n" +
		"   \n" +
		"    //http://m.miguvideo.com/wap/resource/migu/live/live.jsp\n" +
		"    fetch_url = {method: \"get\", header: {user_agent: agent}, body: \"\", url: \"\"};\n" +
		"\n" +
		"    result = {done: false, fetch_url: fetch_url, urls: [], call_back_data: {}};\n" +
		"\n" +
		"    if(json_data.content_type == 2){\n" +
		"        if (json_data.times === 1) {\n" +
		"            try {\n" +
		"                begin = json_data.url.indexOf(\"cid=\");\n" +
		"                var cid = json_data.url.slice(begin + 4);\n" +
		"                fetch_url.url = \"http://www.miguvideo.com/wap/resource/pc/data/miguData.jsp?cid=\" + cid;\n" +
		"                result.done = false;\n" +
		"            } catch (e) {\n" +
		"                result.error = e.toString();\n" +
		"                fetch_url.url = \"\";\n" +
		"                result.done = true;\n" +
		"                result.urls.splice(0, result.urls.length);\n" +
		"                return JSON.stringify(result);\n" +
		"            }\n" +
		"        } else if (json_data.times === 2) {\n" +
		"            var response = json_data.html_data;\n" +
		"            var data = JSON.parse(response);\n" +
		"            var isPingBi = data[0].isPingBi;\n" +
		"            var cmId = data[0].cmId;\n" +
		"            var mId = data[0].mId;\n" +
		"            var contId = data[0].contId;\n" +
		"            var prdpackId = data[0].prdpackId;\n" +
		"            var vId = data[0].vId;\n" +
		"            var totalTime = data[0].timeLong;\n" +
		"            var nodeId = data[0].nodeid;\n" +
		"            var isV = data[0].isV;\n" +
		"            var isOrd = data[0].isOrd;\n" +
		"            var isDongao = data[0].isDongao;\n" +
		"            var isVip = data[0].isVip;\n" +
		"            var playId = data[0].playId;\n" +
		"            var cmId = data[0].cmId;\n" +
		"\n" +
		"            if (mId != '5101073462' && mId != '5101094523' && mId != '5101049700' && isDongao != '1' && isPingBi != '1') {\n" +
		"                fetch_url.url = \"http://www.miguvideo.com/playurl/v1/play/playurlh5?contId=\" + playId + \"&rateType=1,2,3\";\n" +
		"                result.fetch_url.header = {user_agent: agent,\n" +
		"                    \"userId\": \"\",\"userToken\": \"\",\"SDKCEId\": \"79acd784-cbbb-4cae-8778-8723e001164b\",\"clientId\": \"12345678\"};\n" +
		"                result.done = false;\n" +
		"                result.call_back_data.vid = vId;\n" +
		"                //result.urls.splice(0, result.urls.length);\n" +
		"                return JSON.stringify(result);\n" +
		"            }\n" +
		"\n" +
		"            if (mId == '5101073462' || mId == '5101094523' || mId == '5101049700' || isDongao == '1' || isPingBi == '1') {\n" +
		"                result.error = \"因版权方要求，请至咪咕视频APP观看\";\n" +
		"                fetch_url.url = \"\";\n" +
		"                result.done = true;\n" +
		"                result.call_back_data.vid = vId;\n" +
		"                result.urls.splice(0, result.urls.length);\n" +
		"                return JSON.stringify(result);\n" +
		"            }\n" +
		"\n" +
		"            var cidList = data[0].SubSerial_IDS;\n" +
		"            var contIds = cidList.split(\",\");\n" +
		"            for (var i = 0; i < contIds.length; i++) {\n" +
		"                if (contIds[i] != \"\" && contIds[i] != null) {\n" +
		"                    contId.push(contIds[i]);\n" +
		"                }\n" +
		"            }\n" +
		"\n" +
		"            fetch_url.url = \"\";\n" +
		"            result.done = true;\n" +
		"            eval(response[0].func);\n" +
		"            var playUrl = _mv_addr(response[0].play);\n" +
		"            playUrl = decodeURIComponent(playUrl);\n" +
		"            result.urls.push(playUrl);\n" +
		"        }else if(json_data.times == 3) {\n" +
		"            if (json_data.call_back_data.vId != \"\") {\n" +
		"                fetch_url.url = \"http://www.miguvideo.com/vod2/v1/query_spotviurl_cdn?vid=\" + json_data.call_back_data.vId;\n" +
		"                result.done = false;\n" +
		"                result.urls.splice(0, result.urls.length);\n" +
		"                return JSON.stringify(result);\n" +
		"            }\n" +
		"            var response = json_data.html_data;\n" +
		"            var data = JSON.parse(response);\n" +
		"            var playUrl = data.body.urlInfo.url;\n" +
		"            // result.urls.bq.push(encodeURIComponent(playUrls[0].url.replace(/gslbmgspvod.miguvideo.com/g, \"vod.hcs.cmvideo.cn:8088\") ));\n" +
		"            // result.urls.gq.push(encodeURIComponent(playUrls[1].url.replace(/gslbmgspvod.miguvideo.com/g, \"vod.hcs.cmvideo.cn:8088\") ));\n" +
		"            // result.urls.cq.push(encodeURIComponent(playUrls[2].url.replace(/gslbmgspvod.miguvideo.com/g, \"vod.hcs.cmvideo.cn:8088\") ));\n" +
		"            result.done = true;\n" +
		"            result.urls.push(playUrl);\n" +
		"            return JSON.stringify(result);\n" +
		"\n" +
		"        }else if(json_data.times == 4){\n" +
		"            var response = json_data.html_data;\n" +
		"            var data = JSON.parse(response);\n" +
		"            var playUrls = data.result.list;\n" +
		"            // result.urls.bq.push(encodeURIComponent(playUrls[0].url.replace(/gslbmgspvod.miguvideo.com/g, \"vod.hcs.cmvideo.cn:8088\") ));\n" +
		"            // result.urls.gq.push(encodeURIComponent(playUrls[1].url.replace(/gslbmgspvod.miguvideo.com/g, \"vod.hcs.cmvideo.cn:8088\") ));\n" +
		"            // result.urls.cq.push(encodeURIComponent(playUrls[2].url.replace(/gslbmgspvod.miguvideo.com/g, \"vod.hcs.cmvideo.cn:8088\") ));\n" +
		"            result.done = true;\n" +
		"            return JSON.stringify(result);\n" +
		"        }\n" +
		"    }else if(json_data.content_type == 4){\n" +
		"        if(json_data.times === 1){\n" +
		"                try{\n" +
		"                    begin = json_data.url.indexOf(\"cid=\");\n" +
		"                    var cid = json_data.url.slice(begin + 4);\n" +
		"                    fetch_url.url = \"http://h5spdegrade.miguvideo.com/wap/resource/migu/detail/detail_Live_data.jsp?cid=\" + cid + \"&range=0\";\n" +
		"                    result.done = false;\n" +
		"                }catch (e){\n" +
		"                    result.error = e.toString();\n" +
		"                    fetch_url.url = \"\";\n" +
		"                    result.done = true;\n" +
		"                    result.urls.splice(0,result.urls.length);\n" +
		"                    return JSON.stringify(result);\n" +
		"                }\n" +
		"            }else if(json_data.times === 2){\n" +
		"                try{\n" +
		"                    var data = json_data.html_data;\n" +
		"                    var response = JSON.parse(data);\n" +
		"\n" +
		"                    var t = new Date(json_data.start_time);\n" +
		"                    var start_time = getTime(t);\n" +
		"\n" +
		"                    for(i = 0 ; i < response.length ; i ++) {\n" +
		"                        if (response[i].billtime === response[i].nowtime) {\n" +
		"                            //1 流畅 2 标清 3 高清 4 720P 5 1080P\n" +
		"                            //idx: 1 标清 2 高清 3 超清\n" +
		"                            var indx = 1;\n" +
		"                            if (quality === 4) {\n" +
		"                                indx = 3;\n" +
		"                            } else if (quality === 3) {\n" +
		"                                indx = 2;\n" +
		"                            } else if (quality === 2) {\n" +
		"                                indx = 1;\n" +
		"                            } else {\n" +
		"                                fetch_url.url = \"\";\n" +
		"                                result.done = true;\n" +
		"                                result.urls.splice(0, result.urls.length);\n" +
		"                                return JSON.stringify(result);\n" +
		"                            }\n" +
		"\n" +
		"                            fetch_url.url = \"http://h5spdegrade.miguvideo.com/wap/resource/migu/detail/detail_LiveBackSee_data.jsp?playbillId=\" + response[i].playbillId + \"&indx=\" + indx;\n" +
		"                            result.done = false;\n" +
		"\n" +
		"                            break;\n" +
		"                        }\n" +
		"                    }\n" +
		"                }catch (e){\n" +
		"                    result.error = e.toString();\n" +
		"                    fetch_url.url = \"\";\n" +
		"                    result.done = true;\n" +
		"                    result.urls.splice(0,result.urls.length);\n" +
		"                    return JSON.stringify(result);\n" +
		"                }\n" +
		"\n" +
		"            }else if(json_data.times === 3){\n" +
		"                try{\n" +
		"                    data = json_data.html_data;\n" +
		"                    response = JSON.parse(data);\n" +
		"\n" +
		"                    fetch_url.url = \"\";\n" +
		"                    result.done = true;\n" +
		"                    eval(response.deEncrptJsFunc);\n" +
		"                    result.urls.push(_mv_addr(response.liveBackPlayUrl));\n" +
		"                }catch (e){\n" +
		"                    result.error = e.toString();\n" +
		"                    fetch_url.url = \"\";\n" +
		"                    result.done = true;\n" +
		"                    result.urls.splice(0,result.urls.length);\n" +
		"                    return JSON.stringify(result);\n" +
		"                }\n" +
		"            }\n" +
		"    }\n" +
		"\n" +
		"    return JSON.stringify(result);\n" +
		"}\n" +
		"\n"

	return jsCode
}