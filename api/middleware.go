package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
	"github.com/septemhill/test/utils"
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

func Mailer(c *gin.Context) *utils.Mailer {
	return c.MustGet(module.RESOURCE_MAILER).(*utils.Mailer)
}

func SetLogger(logger *logrus.Logger) gin.HandlerFunc {
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

func SetTestPostgreSqlDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := db.ConnectToTest()
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

func SetMailer(mailer *utils.Mailer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(module.RESOURCE_MAILER, mailer)
		c.Next()
	}
}
