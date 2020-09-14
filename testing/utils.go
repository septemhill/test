package testing

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/utils"
)

func NewTestRouter(r *gin.Engine, apis ...utils.ServiceAPI) *httptest.Server {
	r.Use(middleware.SetTestPostgresDB())
	r.Use(middleware.SetRedisDB())

	for _, api := range apis {
		api(r)
	}

	return httptest.NewServer(r)
}

func NewTestPostgresDB() *db.DB {
	return db.ConnectToTest()
}

func NewTestRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{})
}

func NewTestDB() (*db.DB, *redis.Client) {
	return NewTestPostgresDB(), NewTestRedisDB()
}
