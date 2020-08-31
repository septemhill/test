package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
	"github.com/sirupsen/logrus"
)

func validateSessionToken(c *gin.Context) {
	token := c.GetHeader(HEADER_SESSION_TOKEN)
	rdb := RedisDB(c)
	if _, err := rdb.Get(token).Result(); err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errMessage": "Invalid session token"})
		return
	}
	c.Next()
}

func Logger(c *gin.Context) *logrus.Logger {
	return c.MustGet(module.RESOURCE_LOG).(*logrus.Logger)
}

func RedisDB(c *gin.Context) *redis.Client {
	return c.MustGet(module.RESOURCE_MDB).(*redis.Client)
}

func PostgresDB(c *gin.Context) *db.DB {
	return c.MustGet(module.RESOURCE_RDB).(*db.DB)
}

func SetLogger() gin.HandlerFunc {
	f, err := os.OpenFile("api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0775)
	if err != nil {
		panic("logger initialize failed: " + err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(f)

	return func(c *gin.Context) {
		c.Set(module.RESOURCE_LOG, logger)
		c.Next()
	}
}

func SetPostgreSqlDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := db.Connect()
		c.Set(module.RESOURCE_RDB, d)
		c.Next()
		d.Close()
	}
}

func SetRedisDB() gin.HandlerFunc {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return func(c *gin.Context) {
		c.Set(module.RESOURCE_MDB, r)
		c.Next()
		r.Close()
	}
}
