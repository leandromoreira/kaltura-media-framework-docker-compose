package handlers

import (
	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context)

var All = map[string]Handler{
	"connect": Connect,
}
