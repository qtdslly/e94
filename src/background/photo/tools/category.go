package main
import (
	"fmt"
	"io/ioutil"
	"background/common/logger"
	"os"
	"runtime"
	"strings"
	"io"
	"time"
	"background/photo/model"
	"background/photo/config"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"

	"flag"
	"log"
)
func main() {
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

	runtime.GOMAXPROCS(runtime.NumCPU())
	Category("C:/work/photo/result/",db)
	fmt.Println("OK!")
}
// 执行操作
func Category(path string,db *gorm.DB) {
	resultPath := "C:/work/photo/category/"
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		origin := path + file.Name()
		result := resultPath
		if strings.Contains(strings.ToLower(file.Name()), "small") {
			result += "small/"
			if strings.Contains(strings.ToLower(file.Name()), "-v-"){
				result += "v/"
			}else{
				result += "h/"
			}

		}else{
			result += "big/"
			if strings.Contains(strings.ToLower(file.Name()), "-v-"){
				result += "v/"
			}else{
				result += "h/"
			}
		}

		now := time.Now()
		dir := fmt.Sprintf("%04d%02d%02d%02d%02d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute())
		result += dir
		exists,err := PathExists(result)
		if err != nil{
			logger.Error(err)
			return
		}
		if !exists{
			err := os.Mkdir(result, os.ModePerm)
			if err != nil {
				fmt.Printf("mkdir failed![%v]\n", err)
			} else {
				fmt.Printf("mkdir success!\n")
			}
		}
		result += "/" + file.Name()
		fmt.Println("正在处理" + origin + ">>>" + result)
		if _ ,err := CopyFile(result,origin) ; err != nil{
			logger.Debug(err)
			return
		}

		var photo model.Photo
		photo.Url = strings.Replace(result,"C:/work/photo/category","",-1)
		photo.Count = 0
		photo.State = model.PhotoStateOnLine
		photo.Title = "老婆和果果"
		if err = db.Create(&photo).Error ; err != nil{
			logger.Error(err)
			return
		}
		time.Sleep(time.Second * 4)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}