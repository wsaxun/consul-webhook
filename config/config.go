package config

import (
	"flag"
	"log"
	"os"
	"path"
	"sync"

	"github.com/spf13/viper"
)

var (
	config     *Config
	configOnce sync.Once
)

func InitConfig() {
	configOnce.Do(func() {
		var cfgPath string
		flag.StringVar(&cfgPath, "config", "config.toml", "")
		flag.Parse()
		InitConfigByViper(cfgPath)
	})
}

func InitConfigByViper(configPath string) {
	// 获取此文件所在的目录
	dir := path.Dir(configPath)
	if !path.IsAbs(dir) {
		wd, err := os.Getwd()
		if err != nil {
			panic("非绝对目录")
		}
		dir = path.Join(wd, dir)
	}
	filePath := path.Join(dir, path.Base(configPath))
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic("读取文件出错")
	}
	if err = viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	log.Print("初始化配置文件成功")
}

func GetConfig() *Config {
	return config
}

type Config struct {
	App struct {
		Port string `toml:"port"`
		Name string `toml:"name"`
		Env  string `toml:"env"`
	} `toml:"app"`
	Consul struct {
		Address string `toml:"address"`
		Scheme  string `toml:"scheme"`
		Token   string `toml:"token"`
	} `toml:"consul"`
	Repo *Repo `toml:"repo"`
}

type Repo struct {
	Url      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Path     string `toml:"path"`
}
