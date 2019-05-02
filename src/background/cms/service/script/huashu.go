package script

import (
	"background/common/logger"
	"io/ioutil"
	"net/http"
	"bytes"
	"background/newmovie/model"
)

func GetHuashuRealPlayUrl(playUrl model.PlayUrl)(string){
	//jsCode := GetHuashuJsCode()

	var qualityCode string
	if playUrl.Quality == 4{
		qualityCode = "1000996"
	}else{
		qualityCode = "1000995"
	}

	huaShuSrc := "<message module=\"CATALOG_SERVICE\" version=\"1.0\">\n" +
		"\t<header action=\"REQUEST\" command=\"CONTENT_QUERY\" sequence=\"20121030212732_103861\" component-id=\"SYSTEM2\" component-type=\"THIRD_PARTY_SYSTEM\" />\n" +
		"\t<body>\n" +
		"\t\t<contents>\n" +
		"\t\t\t<content>\n" +
		"\t\t\t\t<code>" + playUrl.Url + "</code>\n" +
		"\t\t\t\t<site-code>1000889</site-code>\n" +
		"\t\t\t\t<items-index>-1</items-index>\n" +
		"\t\t\t\t<folder-code>" + qualityCode + "</folder-code>\n" +
		"\t\t\t\t<format>-1</format>\n" +
		"\t\t\t</content>\n" +
		"\t\t</contents>\n" +
		"\t</body>\n" +
		"</message>"

	logger.Debug(huaShuSrc)
	requ, err := http.NewRequest("POST", "http://101.71.69.172:8080/wasu_catalog/catalog", bytes.NewBuffer([]byte(huaShuSrc)))
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

	data := string(recv)



	return data
}

//func GetHuashuJsCode()(string){
//
//	jsCode := "function GetHuashuRealPlayUrl(request) {\n" +
//		"    var json_data;\n" +
//		"    var fetch_url;\n" +
//		"    var result;\n" +
//		"\n" +
//		"    json_data = JSON.parse(request);\n" +
//		"    var referer = json_data.url;\n" +
//		"\n" +
//		"    var agent = \"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36\";\n" +
//		"    fetch_url = {method: \"POST\", header: {user_agent: agent,referer: referer}, body: \"\", url: \"\"};\n" +
//		"    result = {done: false, fetch_url: fetch_url, urls: [], call_back_data: {}};\n" +
//		"    var quality;  //quality  : 1 流畅 2 标清 3 高清 4 720P 5 1080P\n" +
//		"    if(typeof(json_data.quality) === \"undefined\" ){\n" +
//		"        quality = 3;\n" +
//		"    }else{\n" +
//		"        quality = json_data.quality;\n" +
//		"    }\n" +
//		"\n" +
//		"    if(json_data.times === 1){\n" +
//		"            try{\n" +
//		"                var quality_code;\n" +
//		"                if(quality === 3){\n" +
//		"                    quality_code = 1000995; /*标清 720 * 576*/\n" +
//		"                }else if(quality === 4){\n" +
//		"                    /*高清 1280* 720*/\n" +
//		"                    quality_code = 1000996;\n" +
//		"                }else{\n" +
//		"                    fetch_url.url = \"\";\n" +
//		"                    result.done = true;\n" +
//		"                    result.urls.splice(0,result.urls.length);\n" +
//		"                    return JSON.stringify(result);\n" +
//		"                }\n" +
//		"                fetch_url.header = {\"Content-Type\":\"application/xml\"};\n" +
//		"                fetch_url.body = \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n\" +\n" +
//		"                    \"<message module=\\\"CATALOG_SERVICE\\\" version=\\\"1.0\\\">\\n\" +\n" +
//		"\t\t\t\t\t\"\\t<header action=\\\"REQUEST\\\" command=\\\"CONTENT_QUERY\\\" sequence=\\\"20121030212732_103861\\\" component-id=\\\"SYSTEM2\\\" component-type=\\\"THIRD_PARTY_SYSTEM\\\" />\\n\" +\n" +
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
//		"                fetch_url.url = \"http://101.71.69.172:8080/wasu_catalog/catalog\";\n" +
//		"                result.done = false;\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString();\n" +
//		"                fetch_url.url = \"\";\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }else if(json_data.times === 2){\n" +
//		"            try{\n" +
//		"                var re = \"<playUrl>.*</playUrl>\";\n" +
//		"                var players = json_data.html_data.match(re)[0].toString();\n" +
//		"                players = players.replace(\"<![CDATA[\",\"\").replace(\"]]>\",\"\").replace(/\\s+/g,\"\").replace(\"<playUrl>\",\"\").replace(\"</playUrl>\",\"\");\n" +
//		"\n" +
//		"                fetch_url.url = \"\"\n" +
//		"                result.done = true;\n" +
//		"                result.urls.push(players);\n" +
//		"            }catch (e){\n" +
//		"                result.error = e.toString()\n" +
//		"                fetch_url.url = \"\";\n" +
//		"                result.done = true;\n" +
//		"                result.urls.splice(0,result.urls.length);\n" +
//		"                return JSON.stringify(result);\n" +
//		"            }\n" +
//		"        }\n" +
//		"\n" +
//		"    return JSON.stringify(result);\n" +
//		"}\n" +
//		"\n"
//
//	return jsCode
//}