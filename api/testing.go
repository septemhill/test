package api

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
)

func newTestRouter(r *gin.Engine, apis ...ServiceAPI) *httptest.Server {
	r.Use(SetTestPostgreSqlDB())
	r.Use(SetRedisDB())

	for _, api := range apis {
		api(r)
	}

	return httptest.NewServer(r)
}

func newTestDB() (*db.DB, *redis.Client) {
	d := db.ConnectToTest()
	r := redis.NewClient(&redis.Options{})

	return d, r
}
