package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Connect(_ RequestMessage, c *gin.Context) (int, map[string]interface{}) {
	return http.StatusOK, map[string]interface{}{"code": "ok", "message": "success"}
}
