package web

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leandromoreira/kaltura-media-framework-docker-compose/handlers"
)

func Start() {
	r := Setup()
	controller_port := "9191"
	if cp := os.Getenv("CONTROLLER_PORT"); cp != "" {
		controller_port = cp
	}
	r.Run(":" + controller_port)
}

func Setup() *gin.Engine {
	r := gin.Default()

	if debug := os.Getenv("CONTROLLER_DEBUG"); debug != "" {
		r.Use(RequestLoggerMiddleware())
	}

	r.POST("/control", func(c *gin.Context) {
		var requestMessage handlers.RequestMessage
		err := c.BindJSON(&requestMessage)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UNKNOWN_PAYLOAD", "message": err})
			return
		}

		if requestMessage.EventType == "" || handlers.All[requestMessage.EventType] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UNKNOWN_EVENT_TYPE", "message": "Unknown event type"})
			return
		}
		handlers.All[requestMessage.EventType](requestMessage, c)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r
}
