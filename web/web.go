package web

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leandromoreira/kaltura-media-framework-docker-compose/handlers"
)

func unimplemented_logging(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("<<NOT IMPLEMENTED> Req Path=%+v Query=%+v\n Headers=%+v\n Body=%+v\n Context=%+v \n", c.Request.URL.Path, c.Request.Header, c.Request.URL.Query(), body, c)
}

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
		c.BindJSON(&requestMessage)

		if requestMessage.EventType == "" || handlers.All[requestMessage.EventType] == nil {
			unimplemented_logging(c)
			c.JSON(http.StatusBadRequest, gin.H{"code": "UNKNOWN_EVENT_TYPE", "message": "Unknown event type"})
			return
		}
		handlers.All[requestMessage.EventType](c)
	})

	r.NoRoute(func(c *gin.Context) {
		unimplemented_logging(c)

		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r
}
