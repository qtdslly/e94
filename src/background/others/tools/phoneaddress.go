package main

import (

	"background/others/config"
	"background/common/logger"
	"background/others/model"

	"flag"
	"log"
	"strings"
	"fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/imroc/req"
	"github.com/jinzhu/gorm"
	"github.com/axgle/mahonia"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	logger.SetLevel(config.GetLoggerLevel())

	db, err := gorm.Open(config.GetDBName(), config.GetDBSource())
	if err != nil {
		logger.Fatal("Open db Failed!!!!", err)
		return
	}

	db.LogMode(true)

	model.InitModel(db)

	SyncPhoneAddress(db)
}

func GetPhoneAddress(phone string)(string){
	url := "http://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel=" + phone
	resp, err := req.Post(url)
	if err != nil {
		logger.Error("post error" , err)
		return ""
	}

	resbody := mahonia.NewDecoder("gbk").NewReader(resp.Response().Body)
	result, err := ioutil.ReadAll(resbody)
	if err != nil{
		return ""
	}
	return string(result)
}
func SyncPhoneAddress(db *gorm.DB){
	var statusConfig model.StatusConfig

	if err := db.Where("`key` = 'phone_address_key'").First(&statusConfig).Error ; err != nil{
		logger.Error("query status_config error:",err)
		return
	}
	type Status struct{
		PhoneHead  string  `gorm:"phone_head" json:"phone_head"`
		Start      int     `gorm:"start" json:"start"`
	}

	var p Status
	if err := json.Unmarshal([]byte(statusConfig.Value),&p) ; err != nil{
		logger.Error(err)
		return
	}

	logger.Print("phonehead:" + p.PhoneHead)
	logger.Print("start:" + fmt.Sprint(p.Start))


	heads := []string{"152","153","180","181","189","177","173","149","130","131","132",
		"155","156","145","185","186","176","175","134","135","136","137",
		"138","139","150","151","157","158","159","182","183","184",
		"187","188","147","178"}

	start := p.Start
	found := false
	for _,head := range heads{
		if found == true{
			start = 0
		}
		if head != p.PhoneHead && found == false{
			continue
		}
		found = true
		for start <= 9999{
			phoneNum := head + fmt.Sprintf("%04d0000",start)
			fmt.Println(phoneNum)
			phoneInfo := GetPhoneAddress(phoneNum)
			if len(phoneInfo) < 60{
				fmt.Println("手机号码" + phoneNum + "不存在")
			}else{
				fmt.Println(phoneInfo)
				var phoneAddress model.PhoneAddress
				phoneAddress = GetPhoneModel(phoneInfo)
				log.Println(phoneAddress.Mts)
				log.Println(phoneAddress.Province)
				log.Println(phoneAddress.CatName)
				log.Println(phoneAddress.AreaVid)
				log.Println(phoneAddress.IspVid)
				log.Println(phoneAddress.Carrier)
				if err := db.Save(&phoneAddress).Error ; err != nil{
					logger.Error("Save phone_address error:" , err)
					return
				}

			}
			start++
			//time.Sleep(time.Second * 10)
		}
	}
}


func GetPhoneModel(phoneInfo string)(phoneAddress model.PhoneAddress){
	phoneInfo = phoneInfo[strings.Index(phoneInfo,"{") + 1 : len(phoneInfo) - 2 ]
	infos := strings.Split(phoneInfo,",")
	phoneAddress.Mts = infos[0][10:len(infos[0]) - 1]
	phoneAddress.Province = infos[1][15:len(infos[1]) - 1]
	phoneAddress.CatName = infos[2][14:len(infos[2]) - 1]
	phoneAddress.AreaVid = infos[4][11:len(infos[4]) - 1]
	phoneAddress.IspVid = infos[5][11:len(infos[5]) - 1]
	phoneAddress.Carrier = infos[6][11:len(infos[6]) - 2]
	return phoneAddress
}