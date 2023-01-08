package services

import (
	"os"
	"time"
)

func CreateChannels(channel_id string, preset string) {
	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	parameters := []map[string]interface{}{
		map[string]interface{}{
			"uri":    "/channels",
			"method": "POST",
			"body":   map[string]interface{}{"id": channel_id, "preset": preset, "initial_segment_index": time.Now().Unix()},
		},
		map[string]interface{}{
			"uri":    "/channels/" + channel_id + "/timelines",
			"method": "POST",
			"body":   map[string]interface{}{"id": "main", "active": true, "max_segments": 20, "max_manifest_segments": 10},
		},
	}
	PostWithRetries(segmenter_api_url+"/multi", parameters)
}

func CreateTracks(channel_id string, variant_id string, track_id string, media_type string) {
	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	parameters := []map[string]interface{}{
		{
			"uri":    "/channels/" + channel_id + "/variants",
			"method": "POST",
			"body": map[string]interface{}{
				"id":        variant_id,
				"track_ids": []string{},
				"role":      nil,
				"label":     nil,
				"lang":      nil,
			},
		},
		{
			"uri":    "/channels/" + channel_id + "/tracks",
			"method": "POST",
			"body": map[string]interface{}{
				"id":         track_id,
				"media_type": media_type,
			},
		},
		{
			"uri":    "/channels/" + channel_id + "/variants/" + variant_id + "/tracks",
			"method": "POST",
			"body":   map[string]interface{}{"id": track_id},
		},
	}
	PostWithRetries(segmenter_api_url+"/multi", parameters)
}
