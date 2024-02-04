package config

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
)

type Conf struct {
	AppEnv   string   `json:"app_env"`
	Database Database `json:"database"`
	Host     string   `json:"host"`
	Redis    Redis    `json:"redis"`
	TimeZone int      `json:"time_zone"`
	Detail   Detail   `json:"detail"`
}

type Database struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	DBName   string `json:"db_name"`
}

type Redis struct {
	Host     string `json:"host"`
	Password string `json:"password"`
}

type Detail struct {
	MaxAdEveryday int `json:"max_ad_everyday"`
	MaxActiveAd   int `json:"max_active_ad"`
}

var AppConfig Conf
var BotRefreshTime string
var BotGuildCountRefreshTime string

func Setup(fileLocation string) {
	// load config.json
	raw, err := os.ReadFile(fileLocation)
	if err != nil {
		panic("Error occured while reading config:" + err.Error())
	}
	json.Unmarshal(raw, &AppConfig)

	if AppConfig.AppEnv == "local" {
		BotRefreshTime = "*/50 * * * *"
		BotGuildCountRefreshTime = "0 16 * * *"
	} else if AppConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		BotRefreshTime = "*/50 * * * *"
		BotGuildCountRefreshTime = "0 16 * * *"
	} else {
		BotRefreshTime = "*/50 * * * *"
		BotGuildCountRefreshTime = "0 16 * * *"
	}
}
