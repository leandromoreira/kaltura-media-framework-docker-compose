package handlers

import (
	"github.com/gin-gonic/gin"
)

type Handler func(r RequestMessage, c *gin.Context)

var All = map[string]Handler{
	"connect":   Connect,
	"unpublish": Connect,
	"republish": Republish,
}
