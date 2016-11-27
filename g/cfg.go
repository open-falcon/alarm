package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type QueueConfig struct {
	Sms  string `json:"sms"`
	Mail string `json:"mail"`
}

type RedisConfig struct {
	Addr          string   `json:"addr"`
	MaxIdle       int      `json:"maxIdle"`
	HighQueues    []string `json:"highQueues"`
	LowQueues     []string `json:"lowQueues"`
	UserSmsQueue  string   `json:"userSmsQueue"`
	UserMailQueue string   `json:"userMailQueue"`
}

type ApiConfig struct {
	Portal string `json:"portal"`
	Uic    string `json:"uic"`
	Links  string `json:"links"`
}

type GlobalConfig struct {
	Debug    bool         `json:"debug"`
	UicToken string       `json:"uicToken"`
	Http     *HttpConfig  `json:"http"`
	Queue    *QueueConfig `json:"queue"`
	Redis    *RedisConfig `json:"redis"`
	Api      *ApiConfig   `json:"api"`
	Database string       `json:"database"`
	MaxIdle  int          `json:"maxIdle"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}
