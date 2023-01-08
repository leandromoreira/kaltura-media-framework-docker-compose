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
		log.Printf("Path=%+v Header=%+v\n Query=%+v\n", c.Request.URL.Path, c.Request.Header, c.Request.URL.Query())
		log.Printf("Body=%+v\n", string(body))
		log.Println("<<<<<<<<<< RequestLoggerMiddleware [ END ] >>>>>>>>>>")
		c.Next()
	}
}
