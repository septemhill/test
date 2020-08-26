package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
	"github.com/septemhill/test/api"
)

func serviceInit() {
	router := gin.Default()

	router.Use(api.SetLogger())
	router.Use(api.SetPostgreSqlDB())
	router.Use(api.SetRedisDB())

	api.LoadRootAndAccountService(router)

	if len(os.Args) <= 1 {
		log.Error("port please")
		return
	}

	router.Run(":" + os.Args[1])
}

func main() {
	//serviceInit()
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	s, err := r.Get("ksksksks").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("SSS", s)
}
