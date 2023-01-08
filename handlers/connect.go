package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Connect(_ RequestMessage, c *gin.Context) {
	success := ResponseMessage{Code: "ok", Message: "success"}
	c.JSON(http.StatusOK, success)
}
