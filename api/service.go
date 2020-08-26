package api

import "github.com/gin-gonic/gin"

type ServiceAPI func(*gin.Engine) *gin.Engine

func LoadServices(r *gin.Engine, apis ...ServiceAPI) {
	for _, f := range apis {
		f(r)
	}
}

func LoadRootService(r *gin.Engine) {
	LoadServices(r, RootService)
}

func LoadRootAndAccountService(r *gin.Engine) {
	LoadServices(r, RootService, AccountService)
}
