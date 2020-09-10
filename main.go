package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/septemhill/test/api"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/utils"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	f, err := os.OpenFile("api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0775)
	if err != nil {
		panic("failed on logger initialize: " + err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(f)

	return logger
}

func NewMailer() *utils.Mailer {
	mailer := new(utils.Mailer)
	if err := env.Parse(mailer); err != nil {
		panic("failed on initialize mailer: " + err.Error())
	}

	return mailer
}

func serviceInit() {
	router := gin.Default()

	router.Use(middleware.SetLogger(NewLogger()))
	router.Use(middleware.SetPostgresDB())
	router.Use(middleware.SetRedisDB())
	router.Use(middleware.SetMailer(NewMailer()))

	api.LoadAllServices(router)

	if len(os.Args) <= 1 {
		log.Print("port please")
		return
	}

	router.Run(":" + os.Args[1])
}

func main() {
	serviceInit()
}
