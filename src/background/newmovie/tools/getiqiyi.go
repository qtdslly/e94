package main

import (
	"background/newmovie/config"
	"background/newmovie/model"
	"background/common/logger"
	"background/common/constant"
	"background/common/util"

	"strings"
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"flag"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/PuerkitoBio/goquery"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	i := 1
	for {
		url := "http://list.iqiyi.com/www/1/----------0---11-" + fmt.Sprint(i) +  "-1-iqiyi--.html"
		GetAiQiYiMovie(url,db)
		i = i + 1
		if i == 31{
			break
		}
	}

	//GetYouKuPageInfo("http://v.youku.com/v_show/id_XNzA4ODY0NzQ0.html")
}


func GetAiQiYiMovie(url string,db *gorm.DB){
	query := GetAiQiYiPageInfo(url)

	if query != nil{
		FilterAiQiYiMovieInfo(query,db)
	}
}


func FilterAiQiYiMovieInfo(document *goquery.Document,db *gorm.DB)(){
	movieDoc := document.Find("body").Find(".page-list").Find(".wrapper-content").Find(".site-main").Find(".wrapper-cols").Find(".wrapper-piclist").Find("ul").Find("li")

	movieDoc.Each(func(i int, s *goquery.Selection) {
		title,_ := s.Find("a").Eq(0).Attr("title")
		url,_ := s.Find("a").Eq(0).Attr("href")

		thumb_y,_ := s.Find("img").Eq(0).Attr("src")
		if !strings.Contains(thumb_y,"http"){
			thumb_y = "http:" + thumb_y
		}

		score := s.Find(".score").Eq(0).Text()

		dd := s.Find(".role_info").Eq(0).Find("em")

		actors := ""
		dd.Each(func(i int, t *goquery.Selection) {
			actors += t.Text() + ","
		})
		actors = strings.Replace(actors," ","",-1)
		actors = strings.Replace(actors,"\n","",-1)
		if len(actors) > 0{
			actors = actors[0:len(actors) - 1]
		}
		actors = strings.Replace(actors,"主演:,","",-1)

		logger.Debug(title)
		logger.Debug(url)
		logger.Debug(thumb_y)
		logger.Debug(actors)
		directors,description,tags,publishDate,language := GetAiQiYiInfoByUrl(url)
		logger.Debug(publishDate)
		logger.Debug(language)
		logger.Debug(score)
		logger.Debug(directors)
		logger.Debug(description)
		logger.Debug(tags)

		var video  model.Video
		video.Actors = actors
		video.Title = title
		video.Description = description
		video.Directors = directors
		video.Language = language
		video.PublishDate = publishDate
		video.Actors = actors
		video.ThumbY = thumb_y
		score1,_ := strconv.ParseFloat(score,10)
		video.Score = score1
		video.Pinyin = util.TitleToPinyin(video.Title)
		video.Status = constant.MediaStatusReleased

		now := time.Now()
		video.CreatedAt = now
		video.UpdatedAt = now
		if err := db.Where("title = ?",video.Title).First(&video).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&video)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["description"] = description
			}
			if len(thumb_y) > 0{
				updateMap["thumb_y"] = thumb_y
			}
			if len(language) > 0{
				updateMap["language"] = language
			}
			if len(actors) > 0{
				updateMap["actors"] = actors
			}
			if len(tags) > 0{
				updateMap["tags"] = tags
			}
			if len(directors) > 0{
				updateMap["directors"] = directors
			}
			if len(publishDate) > 0{
				updateMap["publishDate"] = directors
			}
			if len(score) > 0{
				updateMap["score"] = score1
			}

			if err = db.Model(model.Video{}).Where("id=?", video.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}

		var episode model.Episode
		episode.Title = title
		episode.VideoId = video.Id
		episode.Description = description
		episode.Score = score1
		episode.Pinyin = util.TitleToPinyin(video.Title)

		episode.CreatedAt = now
		episode.UpdatedAt = now
		if err := db.Where("video_id = ?",video.Id).First(&episode).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&episode)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["description"] = description
			}
			if len(score) > 0{
				updateMap["score"] = score
			}

			if err = db.Model(model.Episode{}).Where("id = ?", episode.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}

		var playUrl model.PlayUrl
		playUrl.Title = episode.Title
		playUrl.ContentType = constant.MediaTypeEpisode
		playUrl.ContentId = episode.Id
		playUrl.Provider = constant.ContentProviderIqiyi
		playUrl.Url = url
		playUrl.Disabled = true

		playUrl.CreatedAt = now
		playUrl.UpdatedAt = now
		if err := db.Where("content_id = ? and content_type = ? and provider = ?",episode.Id,playUrl.ContentType,playUrl.Provider).First(&playUrl).Error ; err == gorm.ErrRecordNotFound{
			db.Create(&playUrl)
		}else{
			updateMap := make(map[string]interface{})
			if len(description) > 0{
				updateMap["url"] = url
			}

			if err = db.Model(model.PlayUrl{}).Where("id = ?", playUrl.Id).Update(updateMap).Error; err != nil {
				logger.Error(err)
				return
			}
		}
	})
}


func GetAiQiYiInfoByUrl(url string)(string,string,string,string,string){
	query := GetAiQiYiPageInfo(url)

	base := query.Find("#block-B").Find("div").Eq(0)
	tvid,_ := base.Attr("data-player-tvid")
	videoId,_ := base.Attr("data-player-videoid")
	apiUrl := "http://cache.video.iqiyi.com/jp/vi/" + tvid + "/" + videoId + "/?status=1&callback=window.Q.__callbacks__.cbjy2ray"
	//url := "http://mixer.video.iqiyi.com/jp/mixin/videos/" + videoId + "?callback=window.Q.__callbacks__.cbp0q4rh&status=1"
	requ, err := http.NewRequest("GET", apiUrl, nil)

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		fmt.Println("get server jscode error!!!",err)
		return "","","","",""
	}

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return "","","","",""
	}

	data := string(recv)
	//try{window.Q.__callbacks__.cber3tx2({"code":"A00000","data":{"shortTitle":"扫毒","editorInfo":"","videoQipuId":166477800,"rewardAllowed":0,"nurl":"","supId":392525502,"onlineStatus":1,"dtype":3,"pvu":"","issueTime":20131213,"writer":["陈木胜","文隽","凌志民","黄进","谭惠贞"],"payMarkUrl":"","subKey":"764778","up":"2018-06-11 19:37:03","ipLimit":0,"un":"","followerCount":0,"exclusive":0,"vn":"扫毒","pturl":"","sm":0,"startTime":-1,"st":200,"povu":"","subType":7,"showChannelId":1,"sc":0,"cc":0,"etm":"20190122","supName":"扫毒 其他配音","albumAlias":"","mdown":0,"pano":{"type":1},"qiyiPlayStrategy":"","endTime":-1,"isDisplayCircle":0,"nvid":"","ar":"华语","es":1,"bmsg":{"t":"20180611193703","f":"kafka","mid":"b733be75558644b6b19f6bd86db0e968","sp":"891410101,"},"stl":{"d":"http:\/\/meta.video.qiyi.com","stl":[]},"au":"http:\/\/www.iqiyi.com\/v_19rrifvgg6.html","upOrder":1,"ty":20131213,"payMark":0,"plc":{"4":{"downAld":1,"coa":1},"10":{"downAld":1,"coa":1},"5":{"downAld":1,"coa":1},"22":{"downAld":1,"coa":1},"11":{"downAld":1,"coa":1},"12":{"downAld":1,"coa":1},"7":{"downAld":1,"coa":1},"18":{"downAld":1,"coa":1},"13":{"downAld":1,"coa":1},"14":{"downAld":1,"coa":1},"8":{"downAld":1,"coa":1},"1":{"downAld":1,"coa":1},"9":{"downAld":1,"coa":1},"6":{"downAld":1,"coa":1},"16":{"downAld":1,"coa":1},"21":{"downAld":1,"coa":1},"3":{"downAld":1,"coa":1},"17":{"downAld":1,"coa":1},"2":{"downAld":1,"coa":1},"15":{"downAld":1,"coa":1},"19":{"downAld":1,"coa":1},"20":{"downAld":1,"coa":1}},"aid":624246,"tvFocuse":"重回老派港片暴力美学","cType":1,"idl":1,"votes":[],"albumQipuId":166477800,"categoryKeywords":"电影,624246,0,http:\/\/list.iqiyi.com\/www\/1\/------------------.html 华语,1,1,http:\/\/list.iqiyi.com\/www\/1\/1------------------.html 战争,7,2,http:\/\/list.iqiyi.com\/www\/1\/-7-----------------.html 动作,11,2,http:\/\/list.iqiyi.com\/www\/1\/-11-----------------.html 悬疑,289,2,http:\/\/list.iqiyi.com\/www\/1\/-289-----------------.html 犯罪,291,2,http:\/\/list.iqiyi.com\/www\/1\/-291-----------------.html 普通话,20001,0,http:\/\/list.iqiyi.com\/www\/1\/------------------.html","qiyiProduced":0,"stm":"20170610","plg":8054,"asubt":"","coa":0,"keyword":"","authors":"","sid":0,"ntvd":0,"uid":0,"apic":"http:\/\/pic5.qiyipic.com\/image\/20180212\/2c\/fa\/v_50764778_m_601_m5.jpg","d":"陈木胜","vid":"359b3e8558ed43549a4e016ccd7fcc4a","tvPhase":0,"e":"","an":"扫毒","bossStatus":0,"info":"从小一起长大的好兄弟马昊天、张子伟和 苏建秋共同效力于警队扫毒科。在一次临时改变计划的行动后，建秋因为卧底身份不能与妻子过正常生活而心生退意，但在阿天与子伟的劝说下，三人决定进行最后一搏。建秋跟着毒贩老大黑柴前往泰国与毒贩Bobby进行对接，目的是见到行动的最大目标“八面佛”，阿天则和子伟以及同事阿益进行跟进。由于泰国警方的配合不力导致建秋身份暴露，虽然建秋答应阿天继续把交易完成，但失败的导火索已经埋下......","is":1,"tpl":[],"flpps":[],"fl":[],"supType":"DIFF_DUB","qtId":1151316,"s":"","a":"袁泉|卢海鹏|宝儿|卢惠光|林国斌|吴廷烨|马浴柯|吴岱融|罗兰|江若琳|释延能|李兆基|维他亚·潘斯林加姆","c":1,"ppsuploadid":0,"tvid":764778,"actors":["袁泉","卢海鹏","宝儿","卢惠光","林国斌","吴廷烨","马浴柯","吴岱融","罗兰","江若琳","释延能","李兆基","维他亚·潘斯林加姆"],"userVideocount":0,"ma":"古天乐|刘青云|张家辉|方力申","rewardMessage":"","vType":0,"vu":"http:\/\/www.iqiyi.com\/v_19rrifvgg6.html","vpic":"http:\/\/pic5.qiyipic.com\/image\/20180212\/2c\/fa\/v_50764778_m_601_m5.jpg","tvSeason":0,"previewImageUrl":"http:\/\/preimage0.qiyipic.com\/preimage\/20170518\/84\/56\/v_50764778_m_611_m1.jpg","producers":"","ppsInfo":{"name":"扫毒","shortTitle":"扫毒"},"mainActorRoles":["苏建秋","马昊天","张子伟"],"allowEditVVIqiyi":0,"actorRoles":["八面佛","缅娜","Booby","黑柴","益哥","白毛段坤","警署高层(行动)","子伟母亲","机场售票员","通缉犯","警署高层(管理)","八面佛手下"],"tags":[],"isPopup":1,"pd":1,"cnPremTime":0,"lg":0,"is3D":0,"tvEname":"The White Storm","producer":"","presentor":[],"pubTime":"1464003844000","ifs":1,"circle":{"type":0,"id":0},"tg":"华语 战争 动作 剧情 悬疑 犯罪 普通话 友谊 卧底 毒品 伤感 暴力 都市 当代 冒险 中成本 杀戮","commentAllowed":1,"platforms":["PHONE","PAD","PC","PC_APP","TV","PHONE_WEB_IQIYI","PAD_WEB_IQIYI"],"isVip":"","isTopChart":0,"cpa":1,"subt":"","dts":"20180611193703"}});}catch(e){};

	logger.Debug(data)
	data = strings.Replace(data,"try{window.Q.__callbacks__.cbjy2ray(","",-1)
	data = strings.Replace(data,");}catch(e){};","",-1)
	logger.Debug(data)
	data = strings.Replace(data,"\n","",-1)
	data = strings.Replace(data,"\r","",-1)


	type AqiYiData struct {
		PublishDate    uint32       `json:"issueTime"`
		Title    string       `json:"vn"`
		Language    string       `json:"ar"`
		Director    string       `json:"d"`
		Description    string       `json:"info"`
		Tags    string       `json:"tg"`
	}
	type AqiYi struct {
		Code    string       `json:"code"`
		Data    AqiYiData       `json:"data"`
	}

	var aiqiyi AqiYi
	if err = json.Unmarshal([]byte(data), &aiqiyi); err != nil {
		logger.Error(err)
		return "","","","",""
	}

	directors := aiqiyi.Data.Director
	description := aiqiyi.Data.Description
	tags := aiqiyi.Data.Tags
	publishDate := aiqiyi.Data.PublishDate
	language := aiqiyi.Data.Language

	apiUrl = "http://mixer.video.iqiyi.com/jp/mixin/videos/" + tvid + "?callback=window.Q.__callbacks__.cbp0q4rh&status=1"

	requ, err = http.NewRequest("GET", apiUrl, nil)

	resp, err = http.DefaultClient.Do(requ)
	if err != nil {
		fmt.Println("get server jscode error!!!",err)
		return "","","","",""
	}

	recv, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return "","","","",""
	}

	data = string(recv)
	data = strings.Replace(data,"try{window.Q.__callbacks__.cbp0q4rh(","",-1)
	data = strings.Replace(data,");}catch(e){};","",-1)

	type AqiYiData1 struct {
		Score    float32       `json:"score"`
	}

	type AqiYi1 struct {
		Code    string       `json:"code"`
		Data    AqiYiData1       `json:"data"`
	}

	logger.Debug(data)
	var aiqiyi1 AqiYi1
	if err = json.Unmarshal([]byte(data), &aiqiyi1); err != nil {
		logger.Error(err)
		return "","","","",""
	}

	//score := aiqiyi1.Data.Score

	return directors,description,tags,fmt.Sprint(publishDate),language
}

func GetAiQiYiMovieOtherInfo(url string)(string,string,string,string,string,string,string){
	if url == ""{
		return "","","","","","",""
	}

	query := GetAiQiYiPageInfo(url)
	newUrl,_ := query.Find("#bpmodule-playpage-righttitle-code").Find(".tvinfo").Eq(0).Find("h2").Eq(0).Find("a").Eq(0).Attr("href")
	//publish_date := query.Find("#bpmodule-playpage-lefttitle").First(".player-title").First(".title-wrap").First(".desc").First(".video-status").First(".bold").First("span").Text()
	//description := query.Find("#module_basic_intro").First(".mod").First(".c").First(".tab-c").First(".summary-wrap").First(".summary").Text()

	if newUrl == ""{
		return "","","","","","",""
	}
	if !strings.Contains(newUrl,"http"){
		newUrl = "http:" + newUrl
	}
	logger.Debug(newUrl)
	q := GetAiQiYiPageInfo(newUrl)
	base := q.Find("body").Find(".s-body").Find(".yk-content").Find(".mod").Find(".p-base").Find("ul").Eq(0)

	year := base.Find(".p-row").Eq(0).Find("span").Eq(0).Text()
	publish_date := base.Children().Eq(2).Eq(0).Find("span").Eq(0).Text()
	score := base.Find(".p-score").Eq(0).Find("span").Eq(1).Text()
	directors := base.Children().Eq(6).Find("a").Eq(0).Text()
	country := base.Children().Eq(7).Find("a").Eq(0).Text()
	description := base.Children().Eq(12).Find("span").Eq(0).Text()

	tagDoc := base.Children().Eq(8).Find("a")
	var tags string
	tagDoc.Each(func(i int, s *goquery.Selection) {
		tags += s.Text() + ","

	})

	if len(tags) > 0{
		tags = tags[0:len(tags) - 1]
	}

	publish_date = strings.Replace(publish_date,"上映：","",-1)
	description = strings.Replace(description,"简介：","",-1)

	return year,publish_date,score,directors,country,description,tags
}
func GetAiQiYiPageInfo(url string)( *goquery.Document){
	logger.Debug(url)
	query, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return nil
	}
	return query
}