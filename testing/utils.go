package testing

import (
	"context"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
	"github.com/septemhill/test/utils"
	"gopkg.in/guregu/null.v4"
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

func NewAccount(ctx context.Context, db *db.DB) *module.Account {
	name := utils.GenerateRandomString(utils.RANDOM_ALL, 7)
	pass := utils.GenerateRandomString(utils.RANDOM_ALL, 12)
	phone := utils.GenerateRandomString(utils.RANDOM_DIGIT_ONLY, 10)

	acc := &module.Account{
		Username: name,
		Password: pass,
		Email:    name + "@balabababa.com",
		Phone:    null.StringFrom(phone),
	}

	_ = module.CreateAccount(ctx, db, *acc)

	return acc
}

func DeleteAccounts(ctx context.Context, db *db.DB, accs ...*module.Account) {
	for _, acc := range accs {
		_ = module.DeleteAccount(ctx, db, *acc)
	}
}
