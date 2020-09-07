package api

import (
	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/utils"
)

func loadServices(r gin.IRouter, apis ...utils.ServiceAPI) {
	for _, f := range apis {
		f(r)
	}
}

func LoadRootService(r gin.IRouter) {
	loadServices(r, RootService)
}

func LoadRootAndAccountService(r gin.IRouter) {
	loadServices(r, RootService, AccountService)
}

func LoadAllServices(r gin.IRouter) {
	loadServices(r, RootService, AccountService, BlogService)
}
