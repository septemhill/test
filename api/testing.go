package api

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func newTestRouter(r *gin.Engine, apis ...ServiceAPI) *httptest.Server {
	r.Use(SetTestPostgreSqlDB())
	r.Use(SetRedisDB())

	for _, api := range apis {
		api(r)
	}

	return httptest.NewServer(r)
}
