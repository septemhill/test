package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/septemhill/test/api"
)

func serviceInit() {
	router := gin.Default()

	router.Use(api.SetLogger())
	router.Use(api.SetPostgreSqlDB())
	router.Use(api.SetRedisDB())

	api.LoadAllServices(router)

	if len(os.Args) <= 1 {
		log.Error("port please")
		return
	}

	router.Run(":" + os.Args[1])
}

func main() {
	serviceInit()
}
