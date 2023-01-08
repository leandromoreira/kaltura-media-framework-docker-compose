package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leandromoreira/kaltura-media-framework-docker-compose/services"
)

type PublishResponseUpstream struct {
	Id  string `json:"id"`
	URL string `json:"url"`
}
type PublishResponse struct {
	ChannelId string                    `json:"channel_id"`
	TrackId   string                    `json:"track_id"`
	Upstreams []PublishResponseUpstream `json:"upstreams"`
}

func Publish(r RequestMessage, c *gin.Context) (int, map[string]interface{}) {
	is_rtmp := r.RTMP != nil
	is_mpegts := r.Mpegts != nil
	var stream_name *string

	if is_rtmp {
		stream_name = r.RTMP.Name
	} else if is_mpegts {
		stream_name = r.Mpegts.StreamID
	} else {
		return http.StatusBadRequest, map[string]interface{}{"code": "NOT_IMPLEMENTED", "message": "it is only allowed to publish rtpm and mpegts"}
	}

	variant_separator_index := strings.LastIndex(*stream_name, "_")
	stream_name_string := *stream_name

	channel_id := stream_name_string[:variant_separator_index]
	variant_id := stream_name_string[variant_separator_index:]
	media_type := r.MediaInfo.MediaType
	preset := "main"
	is_low_latency := channel_id[0:3] == "ll_"

	if is_low_latency {
		preset = "ll"
		channel_id = channel_id[3:]
	}

	//TODO: it should stop if it was not possible to create a channel
	services.CreateChannels(channel_id, preset)

	// TODO: Check if it's media_type[0]
	media_type_string := *media_type
	track_id := string(media_type_string[0]) + variant_id
	services.CreateTracks(channel_id, variant_id, track_id, media_type_string)

	upstreams := []map[string]interface{}{
		{"id": "main", "url": os.Getenv("SEGMENTER_KMP_URL")},
	}
	publish_payload := map[string]interface{}{
		"channel_id": channel_id,
		"track_id":   track_id,
		"upstreams":  upstreams,
	}
	return http.StatusOK, publish_payload
}
