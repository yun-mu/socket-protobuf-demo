/*
Package config default config
*/
package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// Config 配置
type Config struct {
	AppInfo   appInfo   `json:"AppInfo"`
	Log       logConf   `json:"Log"`
	DB        db        `json:"DB"`
	Redis     redis     `json:"Redis"`
	Security  security  `json:"Security"`
	Wechat    wechat    `json:"Wechat"`
	EmailInfo emailInfo `json:"EmailInfo"`
	Qiniu     qiniu     `json:"Qiniu"`
}

type appInfo struct {
	Env        string `json:"Env"` // example: local, dev, prod
	Slogan     string `json:"Slogan"`
	HTTPAddr   string `json:"HTTPAddr"`
	SocketAddr string `json:"SocketAddr"`
}

type logConf struct {
	LogBasePath string `json:"LogBasePath"`
	LogFileName string `json:"LogFileName"`
}

type security struct {
	Secret string `json:"Secret"`
}

type db struct {
	DriverName  string `json:"DriverName"`
	Host        string `json:"Host"`
	Port        string `json:"Port"`
	DBName      string `json:"DBName"`
	User        string `json:"User"`
	PW          string `json:"PW"`
	AdminDBName string `json:"AdminDBName"`
}

type redis struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
	PW   string `json:"PW"`
}

type wechat struct {
	AppID        string `json:"AppID"`
	AppSecret    string `json:"AppSecret"`
	RedirectURI  string `json:"RedirectURI"`
	ResponseType string `json:"ResponseType"`
}

type emailInfo struct {
	From     string   `json:"From"`
	To       []string `json:"To"`
	UserName string   `json:"UserName"`
	AuthCode string   `json:"AuthCode"`
	Host     string   `json:"Host"`
}

type qiniu struct {
	AccessKey string `json:"AccessKey"`
	SecretKey string `json:"SecretKey"`
	Bucket    string `json:"Bucket"` // 空间
}

// Conf 配置
var Conf *Config

var filePrefix = "/app/config/"

func init() {
	log.Println("begin init all configs")
	initConf()
	log.Println("over init all configs")
}

func initConf() {
	log.Println("begin init default config")

	Conf = &Config{}
	fileName := "default.json"

	if v, ok := os.LookupEnv("CONFIG_PATH_PREFIX"); ok {
		filePrefix = v
	}
	// read default config
	data, err := ioutil.ReadFile(filePrefix + fileName)
	if err != nil {
		log.Println("config-initConf: read default.json error")
		log.Panic(err)
		return
	}
	err = jsoniter.Unmarshal(data, Conf)
	if err != nil {
		log.Println("config-initConf: unmarshal default.json error")
		log.Panic(err)
		return
	}

	if v, ok := os.LookupEnv("Slogan"); ok {
		Conf.AppInfo.Slogan = v
	}

	if v, ok := os.LookupEnv("MONGO_INITDB_ROOT_USERNAME"); ok {
		Conf.DB.User = v
	}
	if v, ok := os.LookupEnv("MONGO_INITDB_ROOT_PASSWORD"); ok {
		Conf.DB.PW = v
	}
	if v, ok := os.LookupEnv("MONGO_INITDB_DATABASE"); ok {
		Conf.DB.DBName = v
	}
	if v, ok := os.LookupEnv("RedisPass"); ok {
		Conf.Redis.PW = v
	}

	if v, ok := os.LookupEnv("WeixinAppID"); ok {
		Conf.Wechat.AppID = v
	}
	if v, ok := os.LookupEnv("WeixinAppSecret"); ok {
		Conf.Wechat.AppSecret = v
	}

	if v, ok := os.LookupEnv("QINIU_ACCESS_KEY"); ok {
		Conf.Qiniu.AccessKey = v
	}
	if v, ok := os.LookupEnv("QINIU_SECRET_KEY"); ok {
		Conf.Qiniu.SecretKey = v
	}
	if v, ok := os.LookupEnv("QINIU_BUCKET"); ok {
		Conf.Qiniu.Bucket = v
	}

	if v, ok := os.LookupEnv("FromEmail"); ok {
		Conf.EmailInfo.To = []string{v}
		Conf.EmailInfo.From = v
	}
	if v, ok := os.LookupEnv("EmailAuthCode"); ok {
		Conf.EmailInfo.AuthCode = v
	}
	if v, ok := os.LookupEnv("ToEmail"); ok {
		Conf.EmailInfo.To = strings.Fields(v)
	}
	if v, ok := os.LookupEnv("EmailHost"); ok {
		Conf.EmailInfo.Host = v
	}

	log.Println("over init default config")
}
