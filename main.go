package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Map map[string]interface{}

func (m Map) M(s string) Map {
	return m[s].(map[string]interface{})
}

func (m Map) S(s string) string {
	return m[s].(string)
}

func MapS(m map[string]interface{}, k string) string {
	return m[k].(string)
}

func MapM(m map[string]interface{}, k string) map[string]interface{} {
	return m[k].(map[string]interface{})
}

func getCCDecodeUpstream() Map {
	return Map{
		"id":          "cc-vid",
		"url":         os.Getenv("CC_DECODER_URL"),
		"resume_from": "last_sent",
	}
}

func publish(map_request_payload Map) Map {
	_, is_rtmp := map_request_payload["rtmp"]
	_, is_mpegts := map_request_payload["mpegts"]
	var stream_name string

	if is_rtmp {
		stream_name = map_request_payload.M("rtmp").S("name")
	} else if is_mpegts {
		stream_name = map_request_payload.M("mpegts").S("stream_id")
	} else {
		fmt.Printf("trying to publish an unknownn payload %+v\n", map_request_payload)
		return Map{}
	}
	und_pos := strings.LastIndex(stream_name, "_")
	channel_id := stream_name[:und_pos]
	variant_id := stream_name[und_pos:]
	media_type := map_request_payload.M("media_info").S("media_type")
	preset := "main"

	if channel_id[0:3] == "ll_" {
		preset = "ll"
		channel_id = channel_id[3:]
	}

	segmenterChannelCreate(channel_id, preset)

	// passthrough
	track_id := string(media_type[0]) + variant_id
	setupSegmenterTrack(channel_id, variant_id, track_id, media_type)

	upstreams := []Map{
		{"id": "main", "url": os.Getenv("SEGMENTER_KMP_URL")},
	}

	// return the publish response
	return Map{
		"channel_id": channel_id,
		"track_id":   track_id,
		"upstreams":  upstreams,
	}
}

func postMulti(url string, requests []Map) {
	retries := 3
	for i := 0; i < retries; i++ {
		j, _ := json.Marshal(requests)
		r, err := http.NewRequest("POST", url+"/multi", bytes.NewBuffer(j))
		if err != nil {
			continue
		}
		defer r.Body.Close()
		var map_request_payload Map
		body, err := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &map_request_payload); err != nil {
			continue
		}

		i = 3
	}
}

func segmenterChannelCreate(channel_id string, preset string) {

	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	parameters := []Map{
		Map{
			"uri":    "/channels",
			"method": "POST",
			"body":   Map{"id": channel_id, "preset": preset, "initial_segment_index": time.Now().Unix()},
		},
		Map{
			"uri":    "/channels/$channelId/timelines",
			"method": "POST",
			"body":   Map{"id": "main", "active": true, "max_segments": 20, "max_manifest_segments": 10},
		},
	}
	postMulti(segmenter_api_url, parameters)
}

func setupSegmenterTrack(channel_id string, variant_id string, track_id string, media_type string) {
	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	postMulti(segmenter_api_url, []Map{
		segmenterVariantCreate(channel_id, variant_id, []string{}, "", "", ""),
		segmenterTrackCreate(channel_id, track_id, media_type),
		segmenterVariantAddTrack(channel_id, variant_id, track_id),
	})

}
func segmenterVariantCreate(channel_id string, variant_id string, track_ids []string, role string, label string, lang string) Map {
	return Map{
		"uri":    "/channels/$channelId/variants",
		"method": "POST",
		"body": Map{
			"id":        variant_id,
			"track_ids": track_ids,
			"role":      role,
			"label":     label,
			"lang":      lang,
		},
	}
}

func segmenterTrackCreate(channel_id string, track_id string, media_type string) Map {
	return Map{
		"uri":    "/channels/$channelId/tracks",
		"method": "POST",
		"body": Map{
			"id":         track_id,
			"media_type": media_type,
		},
	}
}
func segmenterVariantAddTrack(channel_id string, variant_id string, track_id string) Map {
	return Map{
		"uri":    "/channels/$channelId/variants/$variantId/tracks",
		"method": "POST",
		"body":   Map{"id": track_id},
	}
}

func kalmedia_controller(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("start controller request %+v  %+v\n", r, r.URL.Query())

	switch r.Method {
	case "GET", "POST":
		var map_request_payload_raw interface{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading the request body %+v\n", err)
			http.Error(w, "error reading the request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &map_request_payload_raw)
		if err != nil {
			fmt.Printf("error unmarshaling the request body %+v\n", err)
			http.Error(w, "error unmarshaling the request body", http.StatusBadRequest)
			return
		}
		map_request_payload := map_request_payload_raw.(map[string]interface{})

		event_type := MapS(map_request_payload, "event_type")
		if event_type == "" {
			http.Error(w, "event_type is required", http.StatusBadRequest)
		}
		switch event_type {
		case "connect":
			log.Println("kalmedia_controller -> connect")
		case "unpublish":
			log.Println("kalmedia_controller -> unpublish")
			j, _ := json.Marshal(Map{
				"code":    "ok",
				"message": "",
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		case "republish":
			log.Println("kalmedia_controller -> republish")
			upstream_id := MapS(map_request_payload, "id")
			w.Header().Set("Content-Type", "application/json")

			switch upstream_id {
			case "cc":
				j, _ := json.Marshal(getCCDecodeUpstream())
				w.Write(j)
			default:
				j, _ := json.Marshal(Map{
					"url": os.Getenv("SEGMENTER_KMP_URL"),
				})
				w.Write(j)
			}
		case "publish":
			log.Println("kalmedia_controller -> publish")
			j, _ := json.Marshal(publish(Map{
				"url": os.Getenv("SEGMENTER_KMP_URL"),
			}))
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		default:
			http.Error(w, "event_type=["+event_type+"] is invalid", http.StatusBadRequest)
		}
	default:
		fmt.Printf("unexpected verb %+v\n", r)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	http.HandleFunc("/control", kalmedia_controller)

	log.Println("Starting kalmedia_controller!")
	http.ListenAndServe(":80", nil)
}
