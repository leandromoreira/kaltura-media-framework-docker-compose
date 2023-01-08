package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Republish(r RequestMessage, c *gin.Context) {
	upstream_id := r.Id

	switch upstream_id {
	case "cc":
		// TODO: implement cc
		c.JSON(http.StatusBadRequest, gin.H{"code": "NOT_IMPLEMENTED", "message": "cc republish is not implemented"})
		return
	default:
		success := ResponseMessage{Code: "ok", Message: "success", URL: os.Getenv("SEGMENTER_KMP_URL")}
		c.JSON(http.StatusOK, success)
	}
}
