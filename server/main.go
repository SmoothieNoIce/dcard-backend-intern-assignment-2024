package main

import (
	"fmt"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/database/cache"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/route"
	"log"
	_ "net/http/pprof"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// @title           Meow Meow API
// @version         v1
// @description     dcard intern assignment.

// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	kernal()
}

func kernal() {
	rotateLog()
	config.Setup("config.json")
	models.Setup(false, false)
	cache.SetUpDefaultDB()
	r := route.Setup()
	err := r.Run(":8000")
	if err != nil {
		panic(err.Error())
	}
}

func rotateLog() {
	t := time.Now()
	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("storage/logs/server.%s.log", t.Format("2006-01-02")),
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     10,   //days
		Compress:   true, // disabled by default
	})
}
