// Package config provides ...
package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Conf config instance
	Conf Config
	// ModeDev run mode as development
	ModeDev = "dev"
	// ModeProd run mode as production
	ModeProd = "prod"
	// WorkDir workspace dir
	WorkDir string
)

// Mode run mode
type Mode struct {
	Name       string `yaml:"name"`
	EnableHTTP bool   `yaml:"enablehttp"`
	HTTPPort   int    `yaml:"httpport"`
	EnableGRPC bool   `yaml:"enablegrpc"`
	GRPCPort   int    `yaml:"grpcport"`
	Host       string `yaml:"host"`
}

// Database sql database
type Database struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

// General common
type General struct {
	PageNum    int    `yaml:"pagenum"`    // 前台每页文章数量
	PageSize   int    `yaml:"pagesize"`   // 后台每页文章数量
	StartID    int    `yaml:"startid"`    // 文章启始ID
	DescPrefix string `yaml:"descprefix"` // 文章描述前缀
	Identifier string `yaml:"identifier"` // 文章截取标识
	Length     int    `yaml:"length"`     // 文章预览长度
	Timezone   string `yaml:"timezone"`   // 时区
}

// Disqus comments
type Disqus struct {
	ShortName   string `yaml:"shortname"`
	PublicKey   string `yaml:"publickey"`
	AccessToken string `yaml:"accesstoken"`
}

// Twitter card
type Twitter struct {
	Card    string `yaml:"card"`
	Site    string `yaml:"site"`
	Image   string `yaml:"image"`
	Address string `yaml:"address"`
}

// Google analytics
type Google struct {
	URL     string `yaml:"url"`
	Tid     string `yaml:"tid"`
	V       string `yaml:"v"`
	T       string `yaml:"t"`
	AdSense string `yaml:"adsense"`
}

// Qiniu oss
type Qiniu struct {
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
	AccessKey string `yaml:"accesskey"`
	SecretKey string `yaml:"secretkey"`
}

// FeedRPC feedr
type FeedRPC struct {
	FeedrURL string   `yaml:"feedrurl"`
	PingRPC  []string `yaml:"pingrpc"`
}

// Account info
type Account struct {
	Username    string `yaml:"username"` // *
	Password    string `yaml:"password"` // *
	Email       string `yaml:"email"`
	PhoneNumber string `yaml:"phonenumber"`
	Address     string `yaml:"address"`
}

// Blogger info
type Blogger struct {
	BlogName  string `yaml:"blogname"`
	SubTitle  string `yaml:"subtitle"`
	BeiAn     string `yaml:"beian"`
	BTitle    string `yaml:"btitle"`
	Copyright string `yaml:"copyright"`
}

// WiBlogApp config
type WiBlogApp struct {
	Mode

	StaticVersion int      `yaml:"staticversion"`
	HotWords      []string `yaml:"hotwords"`
	General       General  `yaml:"general"`
	Disqus        Disqus   `yaml:"disqus"`
	Google        Google   `yaml:"google"`
	Qiniu         Qiniu    `yaml:"qiniu"`
	Twitter       Twitter  `yaml:"twitter"`
	FeedRPC       FeedRPC  `yaml:"feedrpc"`
	Account       Account  `yaml:"account"`
	Blogger       Blogger  `yaml:"blogger"`
}

// BackupApp config
type BackupApp struct {
	Mode

	BackupTo string `yaml:"backupto"`
	Interval string `yaml:"interval"` // circle backup, default: 7d
	Validity string `yaml:"validity"` // storage days, default: 60d
	Qiniu    Qiniu  `yaml:"qiniu"`    // qiniu config
}

// Config app config
type Config struct {
	RunMode   string    `yaml:"runmode"`
	AppName   string    `yaml:"appname"`
	Database  Database  `yaml:"database"`
	ESHost    string    `yaml:"eshost"`
	WiBlogApp WiBlogApp `yaml:"wiblogapp"`
	BackupApp BackupApp `yaml:"backupapp"`
}

func init() {
	var err error
	WorkDir = workDir()
	filename := filepath.Join(WorkDir, "conf", "app.yml")

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		panic(err)
	}

	// read run mode from env
	Conf.RunMode = ModeDev
	if runmode := os.Getenv("RUN_MODE"); runmode == ModeProd {
		Conf.RunMode = ModeProd
	}

	readDBEnv()
}

func readDBEnv() {
	key := strings.ToUpper(Conf.AppName) + "_DB_DRIVER"
	if d := os.Getenv(key); d != "" {
		Conf.Database.Driver = d
	}
	key = strings.ToUpper(Conf.AppName) + "_DB_SOURCE"
	if s := os.Getenv(key); s != "" {
		Conf.Database.Source = s
	}
}
