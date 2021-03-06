package testing

import (
	"net/http/httptest"

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

	h := map[string]interface{}{
		"username": "testonly",
		"email":    "testonly@fakemail.co",
	}
	_, _ = r.HMSet(module.SessionTokenPrefix(token), h).Result()
	return token
}

func NewTestEntities(router *gin.Engine, apis ...utils.ServiceAPI) (*httptest.Server, *db.DB, *redis.Client, map[string]string) {
	server := NewTestRouter(router, apis...)
	d, r := NewTestDB()
	tk := NewTestSessionToken(r)

	hdr := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

	return server, d, r, hdr
}
