package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getCCDecodeUpstream() map[string]interface{} {
	return map[string]interface{}{
		"id":          "cc-vid",
		"url":         os.Getenv("CC_DECODER_URL"),
		"resume_from": "last_sent",
	}
}

func publish(map_request_payload map[string]interface{}) {
	_, is_rtmp := map_request_payload["rtmp"]
	_, is_mpegts := map_request_payload["mpegts"]
	//_, is_cc := map_request_payload["cc"]
	var stream_name string
	//  var channel_id string
	//  var track_id string

	if is_rtmp {
		stream_name = map_request_payload["rtmp"]["name"]
	} else if is_mpegts {
		stream_name = map_request_payload["mpegts"]["stream_id"]
		//	} else if is_cc {
		//	  channel_id = map_request_payload["cc"]["channel_id"]
		//	  track_id = map_request_payload["cc"]["service_id"]
	} else {
		http.Error(w, "trying to publish an unknownn payload", http.StatusBadRequest)
	}
	und_pos := Strings.LastPos(stream_name, "_")
	channel_id := stream_name[:und_pos]
	variant_id := stream_name[und_pos:]
	media_type, _ = map_request_payload["media_info"]["media_type"]
	preset := "main"

	if channel_id[0:3] == "ll_" {
		preset = "ll"
		channel_id = channel_id[3:]
	}

	segmenterChannelCreate(channel_id, preset)

}

func postMulti(url string, requests []map[string]interface{}) {
	retries := 3
	for i := 0; i < retries; i++ {
		j, _ := json.Marshal(requests)
		r, err := http.NewRequest("POST", url+"/multi", bytes.NewBuffer(j))
		if err != nil {
			continue
		}
		defer r.Body.Close()
		var map_request_payload map[string]interface{}
		if err := json.Unmarshal(r.Body, &map_request_payload); err != nil {
			continue
		}

		i = 3
	}
}

func segmenterChannelCreate(channel_id string, preset string) {

	segmenter_api_url := os.Getenv("SEGMENTER_API_URL")
	parameters := []map[string]interface{}{
		{},
		{},
	}
	postMulti(segmenter_api_url, parameters)
}

func kalmedia_controller(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var map_request_payload map[string]interface{}
		map_request_payload = json.Unmarshal(r.Body, &map_request_payload)

		event_type, ok := map_request_payload["event_type"]
		if !ok {
			http.Error(w, "event_type is required", http.StatusBadRequest)
		}
		switch event_type {
		case "connect":
			log.Println("kalmedia_controller -> connect")
		case "unpublish":
			j, _ := json.Marshal(map[string]interface{}{
				"code":    "ok",
				"message": "",
			})
			w.Write(j)
		case "republish":
			upstream_id, ok := map_request_payload["id"]
			if !ok {
				http.Error(w, "id is required for [republish]", http.StatusBadRequest)
			}
			//channel_id, ok := map_request_payload["id"]
			//if !ok {
			//	http.Error(w, "channel_id is required for [republish]", http.StatusBadRequest)
			//}
			switch upstream_id {
			case "cc":
				j, _ := json.Marshal(getCCDecodeUpstream())
				w.Write(j)
			default:
				j, _ := json.Marshal(map[string]interface{}{
					"url": os.Getenv("SEGMENTER_KMP_URL"),
				})
				w.Write(j)
			}
		case "publish":
		default:
			http.Error(w, "event_type=["+event_type+"] is invalid", http.StatusBadRequest)

		}

		// if only one expected
		param1 := r.URL.Query().Get("param1")
		if param1 != "" {
			// ... process it, will be the first (only) if multiple were given
			// note: if they pass in like ?param1=&param2= param1 will also be "" :|
		}
		// Just send out the JSON version of 'tom'
		j, _ := json.Marshal(tom)
		w.Write(j)
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := &person{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tom = p
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	http.HandleFunc("/controller", kalmedia_controller)

	log.Println("Starting kalmedia_controller!")
	http.ListenAndServe(":80", nil)
}
