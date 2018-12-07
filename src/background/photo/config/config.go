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
	StaticRoot    string `json:"static_root"`
	EnableOrmLog  bool   `json:"enable_orm_log"`
	EnableHttpLog bool   `json:"enable_http_log"`
}

var c config

func init() {
	c.ProductionEnv = false
	c.StorageRoot = "/home/lyric/data/stock/"
	c.LogRoot = "../log/"
	c.DBName = "mysql"
	c.DBSource = "root:hahawap@tcp(47.106.111.101:3306)/photo?charset=utf8&parseTime=True&loc=Local"
	c.StaticRoot = "c:/work/code/e94/src/background/photo/html"
	c.LoggerLevel = 0
	c.EnableOrmLog = true
	c.EnableHttpLog = true
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

func GetStaticRoot() string {
	return c.StaticRoot
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
