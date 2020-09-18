package testing

import (
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
	"github.com/septemhill/test/utils"
)

func NewTestRouter(r *gin.Engine, apis ...utils.ServiceAPI) *httptest.Server {
	r.Use(middleware.SetTestPostgresDB())
	r.Use(middleware.SetTestRedisDB())

	for _, api := range apis {
		api(r)
	}

	return httptest.NewServer(r)
}

func NewTestPostgresDB() *db.DB {
	return db.ConnectToTest()
}

func NewTestRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func NewTestDB() (*db.DB, *redis.Client) {
	return NewTestPostgresDB(), NewTestRedisDB()
}

func NewTestSessionToken(r *redis.Client) string {
	token := utils.GenerateRandomString(utils.RANDOM_HEX_ONLY, module.SESSION_TOKEN_LEN)
	r.Set(module.SessionTokenPrefix(token), "testonly@fakemail.co", time.Hour*1).Result()
	return token
}
