package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leandromoreira/kaltura-media-framework-docker-compose/handlers"
	"github.com/stretchr/testify/assert"
)

func TestUnknownEventType(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	var jsonData = []byte(`{
		"event_type": "foobar"
	}`)
	req, _ := http.NewRequest("POST", "/control", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	expected_message := handlers.ResponseMessage{Code: "UNKNOWN_EVENT_TYPE", Message: "Unknown event type"}
	var response handlers.ResponseMessage
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, expected_message, response)
}

func TestConnectRoute(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	var jsonData = []byte(`{
		"event_type": "connect"
	}`)
	req, _ := http.NewRequest("POST", "/control", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	expected_message := handlers.ResponseMessage{Code: "ok", Message: "success"}
	var response handlers.ResponseMessage
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected_message, response)
}
