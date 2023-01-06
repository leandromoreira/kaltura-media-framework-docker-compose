package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func log(msg string) {
	fmt.Printf(" >>>>>>>>> Kalmedia Controller <<<<<<<<<<<<  %+v\n", msg)
}

func GetStringFromMap(m map[string]interface{}, k string) string {
	return m[k].(string)
}

func GetMapFromMap(m map[string]interface{}, k string) map[string]interface{} {
	return m[k].(map[string]interface{})
}

func CheckExistenceOfKey(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

func getCCDecodeUpstream() map[string]interface{} {
	return map[string]interface{}{
		"id":          "cc-vid",
		"url":         os.Getenv("CC_DECODER_URL"),
		"resume_from": "last_sent",
	}
}

func publish(map_request_payload map[string]interface{}) map[string]interface{} {
	var stream_name string
	is_rtmp := CheckExistenceOfKey(map_request_payload, "rtmp")
	is_mpegts := CheckExistenceOfKey(map_request_payload, "mpegts")

	if is_rtmp {
		stream_name = GetStringFromMap(GetMapFromMap(map_request_payload, "rtmp"), "name")
	} else if is_mpegts {
		stream_name = GetStringFromMap(GetMapFromMap(map_request_payload, "mpegts"), "stream_id")
	} else {
		msg := fmt.Sprintf("trying to publish an unknownn payload %+v\n", map_request_payload)
		log(msg)
		return map[string]interface{}{}
	}
	und_pos := strings.LastIndex(stream_name, "_")
	channel_id := stream_name[:und_pos]
	variant_id := stream_name[und_pos:]
	media_type := GetStringFromMap(GetMapFromMap(map_request_payload, "media_info"), "media_type")
	preset := "main"

	if channel_id[0:3] == "ll_" {
		preset = "ll"
		channel_id = channel_id[3:]
	}

	segmenterChannelCreate(channel_id, preset)

	// passthrough
	track_id := string(media_type[0]) + variant_id
	setupSegmenterTrack(channel_id, variant_id, track_id, media_type)

	upstreams := []map[string]interface{}{
		{"id": "main", "url": os.Getenv("SEGMENTER_KMP_URL")},
	}

	// return the publish response
	publish_payload := map[string]interface{}{
		"channel_id": channel_id,
		"track_id":   track_id,
		"upstreams":  upstreams,
	}
	msg := fmt.Sprintf("trying to publish the payload %+v\n", publish_payload)
	log(msg)
	return publish_payload
}

func postMulti(url string, requests []map[string]interface{}) {
	msg := fmt.Sprintf("[[EXTERNAL POST]] url=%+v\n", url)
	log(msg)

	retries := 3
	for i := 0; i < retries; i++ {

		//var res map[string]interface{}
		//json.NewDecoder(resp.Body).Decode(&res)
		//fmt.Println(res["json"])

		j, err := json.Marshal(requests)
		if err != nil {
			msg := fmt.Sprintf("[[EXTERNAL POST]] error while marshalling try #%+v url=%+v err=%+v\n", i+1, url, err)
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}
		resp, err := http.Post(url+"/multi", "application/json", bytes.NewBuffer(j))
		if err != nil {
			msg := fmt.Sprintf("[[EXTERNAL POST]] error while posting try #%+v url=%+v err=%+v json=%+v\n", i+1, url, err, string(j[:]))
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			msg := fmt.Sprintf("[[EXTERNAL POST]] FAIL STATUS CODE %+v try #%+v url=%+v err=%+v body=%+v\n", resp.Status, i+1, url, err, resp.Body)
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}

		msg := fmt.Sprintf("[[EXTERNAL POST]] SUCCESS POST try #%+v url=%+v response=%+v requests=%+v\n", i+1, url, resp, requests)
		log(msg)
		var map_request_payload interface{}
		err = json.NewDecoder(resp.Body).Decode(&map_request_payload)
		if err != nil {
			msg := fmt.Sprintf("[[EXTERNAL POST]] error while unmarshalling try #%+v url=%+v err=%+v body=%+v\n", i+1, url, err, resp.Body)
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}

		i = 3
	}
}

func segmenterChannelCreate(channel_id string, preset string) {

	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	parameters := []map[string]interface{}{
		map[string]interface{}{
			"uri":    "/channels",
			"method": "POST",
			"body":   map[string]interface{}{"id": channel_id, "preset": preset, "initial_segment_index": time.Now().Unix()},
		},
		map[string]interface{}{
			"uri":    "/channels/$channelId/timelines",
			"method": "POST",
			"body":   map[string]interface{}{"id": "main", "active": true, "max_segments": 20, "max_manifest_segments": 10},
		},
	}
	postMulti(segmenter_api_url, parameters)
}

func setupSegmenterTrack(channel_id string, variant_id string, track_id string, media_type string) {
	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	postMulti(segmenter_api_url, []map[string]interface{}{
		segmenterVariantCreate(channel_id, variant_id, []string{}, "", "", ""),
		segmenterTrackCreate(channel_id, track_id, media_type),
		segmenterVariantAddTrack(channel_id, variant_id, track_id),
	})

}
func segmenterVariantCreate(channel_id string, variant_id string, track_ids []string, role string, label string, lang string) map[string]interface{} {
	return map[string]interface{}{
		"uri":    "/channels/$channelId/variants",
		"method": "POST",
		"body": map[string]interface{}{
			"id":        variant_id,
			"track_ids": track_ids,
			"role":      role,
			"label":     label,
			"lang":      lang,
		},
	}
}

func segmenterTrackCreate(channel_id string, track_id string, media_type string) map[string]interface{} {
	return map[string]interface{}{
		"uri":    "/channels/$channelId/tracks",
		"method": "POST",
		"body": map[string]interface{}{
			"id":         track_id,
			"media_type": media_type,
		},
	}
}
func segmenterVariantAddTrack(channel_id string, variant_id string, track_id string) map[string]interface{} {
	return map[string]interface{}{
		"uri":    "/channels/$channelId/variants/$variantId/tracks",
		"method": "POST",
		"body":   map[string]interface{}{"id": track_id},
	}
}

func kalmedia_controller(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("[[START]] controller request (%+v) %+v\n", r.Method, r)
	log(msg)

	switch r.Method {
	case "GET", "POST":
		var map_request_payload_raw interface{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := fmt.Sprintf("[[ERROR]] reading the request body %+v\n", err)
			log(msg)
			http.Error(w, "error reading the request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &map_request_payload_raw)
		if err != nil {
			msg := fmt.Sprintf("[[ERROR]] unmarshaling the request body %+v\n", err)
			log(msg)
			http.Error(w, "error unmarshaling the request body", http.StatusBadRequest)
			return
		}
		map_request_payload := map_request_payload_raw.(map[string]interface{})

		event_type := GetStringFromMap(map_request_payload, "event_type")
		if event_type == "" {
			http.Error(w, "event_type is required", http.StatusBadRequest)
		}

		switch event_type {
		case "connect", "unpublish":
			msg := fmt.Sprintf("[[%+v]]\n", event_type)
			log(msg)

			j, _ := json.Marshal(map[string]interface{}{
				"code":    "ok",
				"message": "",
			})

			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		case "republish":
			upstream_id := GetStringFromMap(map_request_payload, "id")

			msg := fmt.Sprintf("[[%+v]] upstream_id %+v\n", event_type, upstream_id)
			log(msg)

			switch upstream_id {
			case "cc":
				j, _ := json.Marshal(getCCDecodeUpstream())
				msg := fmt.Sprintf("[[%+v]] upstream_id %+v j %+v \n", event_type, upstream_id, j)
				log(msg)
				w.Header().Set("Content-Type", "application/json")
				w.Write(j)
			default:
				j, _ := json.Marshal(map[string]interface{}{
					"url": os.Getenv("SEGMENTER_KMP_URL"),
				})
				//msg := fmt.Sprintf("[[%+v]] upstream_id %+v j %+v \n", event_type, upstream_id, string(j[:]))
				//log(msg)
				w.Header().Set("Content-Type", "application/json")
				w.Write(j)
			}
		case "publish":
			msg := fmt.Sprintf("[[%+v]]\n", event_type)
			log(msg)
			j, _ := json.Marshal(publish(map_request_payload))
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		default:
			msg := fmt.Sprintf("[[%+v]] unknown event type\n", event_type)
			log(msg)
			http.Error(w, "event_type=["+event_type+"] is invalid", http.StatusBadRequest)
		}
	default:
		msg := fmt.Sprintf("[[HTTP VERB]] unexpected verb for request %+v\n", r)
		log(msg)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	http.HandleFunc("/control", kalmedia_controller)

	msg := fmt.Sprintf("[[STARTING]] kalmedia_controller server \n")
	log(msg)
	http.ListenAndServe(":80", nil)
}
