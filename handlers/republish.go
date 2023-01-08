package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Republish(r RequestMessage, c *gin.Context) (int, map[string]interface{}) {
	upstream_id := r.Id

	switch upstream_id {
	case "cc":
		// TODO: implement cc
		return http.StatusBadRequest, map[string]interface{}{"code": "NOT_IMPLEMENTED", "message": "cc republish is not implemented"}
	default:
		return http.StatusOK, map[string]interface{}{"code": "ok", "message": "success", "url": os.Getenv("SEGMENTER_KMP_URL")}
	}
}
