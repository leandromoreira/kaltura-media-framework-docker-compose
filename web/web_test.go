package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
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

func TestUnpublishRoute(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	var jsonData = []byte(`{
		"event_type": "unpublish"
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

func TestRepublishNonCCRoute(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	var jsonData = []byte(`{
		"event_type": "republish",
		"id": "main"
	}`)
	req, _ := http.NewRequest("POST", "/control", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	expected_message := handlers.ResponseMessage{Code: "ok", Message: "success", URL: "http://segmenter_kmp_url"}
	var response handlers.ResponseMessage
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected_message, response)
}

func TestRepublishCCRoute(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	var jsonData = []byte(`{
		"event_type": "republish",
		"id": "cc"
	}`)
	req, _ := http.NewRequest("POST", "/control", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	expected_message := handlers.ResponseMessage{Code: "NOT_IMPLEMENTED", Message: "cc republish is not implemented"}
	var response handlers.ResponseMessage
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, expected_message, response)
}

// TODO: create table test for all inputs combination (rtmp, video, audio, mpegts, cc, etc)
func TestPublishVideoRTMPRoute(t *testing.T) {
	router := Setup()

	w := httptest.NewRecorder()
	var jsonData = []byte(`{
	  "event_type": "publish",
	  "input_id": "rtmp://testserver:1935/live?arg=value/streamname_1/video",
	  "input_type": "rtmp",
	  "rtmp": {
	      "app": "live",
	      "flashver": "FMLE/3.0 (compatible; FMSc/1.0)",
	      "swf_url": "rtmp://testserver:1935/live?arg=value",
	      "tc_url": "rtmp://testserver:1935/live?arg=value",
	      "page_url": "",
	      "addr": "5.6.7.8",
	      "connection": 983,
	      "name": "channelid_variantid",
	      "type": "live",
	      "args": "videoKeyframeFrequency=5&totalDatarate=200"
	  },
	  "media_info": {
	      "media_type": "video",
	      "bitrate": 82000,
	      "codec_id": 7,
	      "extra_data": "0164000bffe100196764000bacd942847e5c0440000003004000000783c50a658001000468efbcb0",
	      "width": 160,
	      "height": 120,
	      "frame_rate": 15.00
	  }
	}`)
	defer gock.Off()

	channels_and_timele_creation := []map[string]interface{}{
		{
			"uri": "/channels", "method": "POST",
			"body": map[string]interface{}{
				"id":                    "channelid",
				"preset":                "normal",
				"initial_segment_index": 1,
			},
		},
		{
			"uri": "/channels/channelid/timelines", "method": "POST",
			"body": map[string]interface{}{
				"id":                    "main",
				"active":                true,
				"max_segments":          20,
				"max_manifest_segments": 10,
			},
		},
	}
	variantes_and_tracks_creation := []map[string]interface{}{
		{
			"uri": "/channels/channelid/variants", "method": "POST",
			"body": map[string]interface{}{
				"id":        "variantid",
				"tracks_id": []string{},
				"role":      nil,
				"label":     nil,
				"lang":      nil,
			},
		},
		{
			"uri": "/channels/channelid/tracks", "method": "POST",
			"body": map[string]interface{}{
				"id":         "video_variantid",
				"media_type": "video",
			},
		},
		{
			"uri": "/channels/channelid/variants/variantid/tracks", "method": "POST",
			"body": map[string]interface{}{
				"id": "video_variantid",
			},
		},
	}

	gock.New("http://segmenter_api_url").
		Post("/multi").
		MatchType("json").
		JSON(channels_and_timele_creation).
		Reply(200).
		JSON([]map[string]interface{}{
			{"code": 200, "message": "success"},
			{"code": 200, "message": "success"},
		})

	gock.New("http://segmenter_api_url").
		Post("/multi").
		MatchType("json").
		JSON(variantes_and_tracks_creation).
		Reply(200).
		JSON([]map[string]interface{}{
			{"code": 200, "message": "success"},
			{"code": 200, "message": "success"},
		})

	req, _ := http.NewRequest("POST", "/control", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	expected_message := handlers.PublishResponse{
		ChannelId: "channelid",
		TrackId:   "v_variantid",
		Upstreams: []handlers.PublishResponseUpstream{
			handlers.PublishResponseUpstream{Id: "main", URL: "http://segmenter_kmp_url"},
		},
	}
	var response handlers.PublishResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected_message, response)
}
