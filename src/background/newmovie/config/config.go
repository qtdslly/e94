package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	ProductionEnv bool   `json:"production_env"`
	StorageRoot   string `json:"storage_root"`
	StorageTmp    string `json:"storage_tmp"`
	LogRoot       string `json:"log_root"`
	DBName        string `json:"db_name"`
	DBSource      string `json:"db_source"`
	LoggerLevel   uint8  `json:"logger_level"`
	EnableOrmLog  bool   `json:"enable_orm_log"`
	EnableHttpLog bool   `json:"enable_http_log"`
	TmplRoot      string `json:"tmpl_root"`
	StaticRoot    string `json:"static_root"`
	CmsRoot       string `json:"cms_root"`

}

var c config

func init() {
	c.ProductionEnv = false
	c.StorageRoot = "/home/lyric/data/stock/"
	c.LogRoot = "../log/"
	c.DBName = "mysql"
	c.DBSource = "imsuser:imspass@tcp(localhost:3306)/movie?charset=utf8&parseTime=True&loc=Local"
	c.LoggerLevel = 0
	c.EnableOrmLog = true
	c.EnableHttpLog = true
	c.CmsRoot = "/root/Git/e94/src/background/newmovie/"
	c.TmplRoot = "/root/Git/e94/src/background/newmovie/tmpl/"
	c.StaticRoot = "/root/Git/e94/src/background/newmovie/static/"
}

func LoadConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var ctmp config
	err = json.Unmarshal(b, &ctmp)
	if err != nil {
		return err
	}

	c = ctmp
	return nil
}

func IsProductionEnv() bool {
	return c.ProductionEnv
}

func GetLogRoot() string {
	return c.LogRoot
}

func GetDBName() string {
	return c.DBName
}

func GetDBSource() string {
	return c.DBSource
}

func GetLoggerLevel() uint8 {
	return c.LoggerLevel
}

func IsOrmLogEnabled() bool {
	return c.EnableOrmLog
}

func IsHttpLogEnabled() bool {
	return c.EnableHttpLog
}

func GetStorageRoot() string {
	return c.StorageRoot
}

func GetStorageTmp() string {
	return c.StorageTmp
}

func GetTmplRoot() string {
	return c.TmplRoot
}

func GetStaticRoot() string {
	return c.StaticRoot
}

func GetCmsRoot() string {
	return c.CmsRoot
}