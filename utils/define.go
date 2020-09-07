package utils

import "github.com/gin-gonic/gin"

type ServiceAPI func(gin.IRouter) gin.IRouter
