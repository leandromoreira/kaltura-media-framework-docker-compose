package web

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Println("<<<<<<<<<< RequestLoggerMiddleware [START] >>>>>>>>>>")
		log.Println(string(body))
		log.Println(c.Request.Header)
		log.Println("<<<<<<<<<< RequestLoggerMiddleware [ END ] >>>>>>>>>>")
		c.Next()
	}
}
