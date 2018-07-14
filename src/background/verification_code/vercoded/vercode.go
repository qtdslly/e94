package main

import (
	"background/verification_code/captcha"
	"background/verification_code/config"
	"background/common/logger"
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"log"
	"html/template"
	"flag"
)
const (
	dx = 60
	dy = 30
)

func main() {

	configPath := flag.String("conf", "../config/config.json", "Config file path")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Config Failed!!!!", err)
		return
	}

	fontFils,err := ListDir(config.GetFontsDir(),".ttf");
	if(err != nil){
		logger.Error(err)
		return ;
	}

	captcha.SetFontFamily(fontFils...);

	http.HandleFunc("/", Index)
	http.HandleFunc("/get/", Get)
	fmt.Println("服务已启动...");
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(config.GetIndexHtml())
	if err != nil {
		logger.Error(err)
	}
	t.Execute(w, nil)
}
func Get(w http.ResponseWriter, r *http.Request) {
	captchaImage,err := captcha.NewCaptchaImage(dx,dy,captcha.RandLightColor());
	captchaImage.DrawNoise(captcha.CaptchaComplexHigh);
	captchaImage.DrawTextNoise(captcha.CaptchaComplexHigh);
	captchaImage.DrawText(captcha.RandText(4));
	//captchaImage.Drawline(3);
	captchaImage.DrawBorder(captcha.ColorToRGB(0x17A7A7A));
	captchaImage.DrawHollowLine();
	if err != nil {
		logger.Error(err)
	}
	captchaImage.SaveImage(w,captcha.ImageFormatJpeg);
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		logger.Error(err,dirPth)
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}