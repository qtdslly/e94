package main

import (
	"background/newmovie/config"
	"background/common/logger"
	"background/newmovie/model"
	"background/common/util"
	"github.com/PuerkitoBio/goquery"

	"strings"
	"background/common/constant"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"flag"
)

func main(){
	logger.SetLevel(config.GetLoggerLevel())

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		return
	}

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}
	db.LogMode(true)
	getIviChannelList(db)
}


func getIviChannelList(db *gorm.DB){
	url := "http://ivi.bupt.edu.cn/"
	query, err := goquery.NewDocument(url)
	if err != nil {
		logger.Debug(url)
		logger.Error(err)
		return
	}
	//fmt.Println(query.Html())

	movieDoc := query.Find(".container").Find(".row")

	movieDoc.Each(func(i int, s *goquery.Selection) {

		rows := s.Find("div")
		if strings.Contains(rows.Text(),"移动端"){
			rows.Each(func(i int, s *goquery.Selection) {
				a := s.Find("a").Eq(1)
				playUrl ,_ := a.Attr("href")
				if !strings.Contains(playUrl,"http"){
					//http://ivi.bupt.edu.cn/hls/cctv1hd.m3u8
					playUrl = "http://ivi.bupt.edu.cn" + playUrl
				}

				title := s.Find("p").Eq(0).Text()

				//logger.Debug(title,"	",playUrl)
				var stream model.Stream
				stream.Title = title
				stream.Title = strings.Replace(stream.Title,"高清","",-1)
				stream.Title = strings.Replace(stream.Title,"-","",-1)
				stream.Title = TrimChinese(stream.Title)
				stream.Pinyin = util.TitleToPinyin(stream.Title)
				stream.Title = strings.Trim(stream.Title," ")
				logger.Debug(stream.Title)

				tx := db.Begin()
				if err := tx.Where("title = ?",stream.Title).First(&stream).Error ; err == gorm.ErrRecordNotFound{
					if strings.Contains(stream.Title,"CCTV"){
						stream.Category = "央视"
					}else if strings.Contains(stream.Title,"卫视"){
						stream.Category = "卫视"
					}else{
						stream.Category = "地方"
					}

					stream.OnLine = constant.MediaStatusOnLine
					stream.Sort = 0

					if err = tx.Create(&stream).Error ; err != nil{
						tx.Rollback()
						logger.Error(err)
						return
					}
				}

				var play model.PlayUrl
				play.Url = playUrl
				play.Provider = uint32(constant.ContentProviderIvibupt)
				if err := tx.Where("provider = ? and url = ?",play.Provider,play.Url).First(&play).Error ; err == gorm.ErrRecordNotFound{
					play.Title = title
					play.OnLine = constant.MediaStatusOnLine
					play.Sort = 0
					play.ContentType = uint8(constant.MediaTypeStream)
					play.ContentId = stream.Id
					play.Quality = uint8(constant.VideoQuality720p)

					if err = tx.Create(&play).Error ; err != nil{
						tx.Rollback()
						logger.Error(err)
						return
					}
				}
				tx.Commit()
			})
		}
	})

}


func TrimChinese(title string)(string){
	if !strings.Contains(title,"CCTV"){
		return title
	}
	//48-57 45 43 65-90 97-122
	rTitle := ([]rune)(title)
	result := ""
	for _, m := range rTitle {
		if m == 43 || m == 45 || (m >= 48 && m <= 57) || (m >= 65 && m <=90) || (m >= 97 && m <= 122){
			result += string(m)
		}
	}
	return result
}