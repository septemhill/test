package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/api"
	"github.com/septemhill/test/utils"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	f, err := os.OpenFile("api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0775)
	if err != nil {
		panic("logger initialize failed: " + err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(f)

	return logger
}

func NewMailer() *utils.Mailer {
	//mailer := new(utils.Mailer)
	mailer := &utils.Mailer{
		Host:     "smtp.gmail.com",
		Port:     587,
		User:     "septemhill@gmail.com",
		Password: "ntiofslamjztwjjm",
	}
	if err := env.Parse(mailer); err != nil {
		panic("failed on initialize mail info" + err.Error())
	}

	return mailer
}

func serviceInit() {
	router := gin.Default()

	router.Use(api.SetLogger(NewLogger()))
	router.Use(api.SetPostgreSqlDB())
	router.Use(api.SetRedisDB())
	router.Use(api.SetMailer(NewMailer()))

	api.LoadAllServices(router)

	if len(os.Args) <= 1 {
		log.Print("port please")
		return
	}

	router.Run(":" + os.Args[1])
}

func main() {
	//serviceInit()
}
