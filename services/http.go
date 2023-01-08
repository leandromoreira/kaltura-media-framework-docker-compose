package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func log(msg string) {
	markers := " >>>>>>>> HTTP Post <<<<<<< \n"
	fmt.Printf("%+v %+v %+v\n", markers, msg, markers)
}

func PostWithRetries(url string, requests []map[string]interface{}) {
	retries := 3
	for i := 0; i < retries; i++ {
		j, err := json.Marshal(requests)
		if err != nil {
			msg := fmt.Sprintf("[[EXTERNAL POST]] <<ERROR>> while marshalling try #%+v url=%+v err=%+v\n", i+1, url, err)
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}
		//resp, err := http.Post(url+"/multi", "application/json", bytes.NewBuffer(j))
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(j))
		if err != nil {
			msg := fmt.Sprintf("[[EXTERNAL POST]] <<ERROR>> while posting try #%+v url=%+v err=%+v json=%+v\n", i+1, url, err, string(j[:]))
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			msg := fmt.Sprintf("[[EXTERNAL POST]] <<FAIL STATUS CODE>> %+v try #%+v url=%+v err=%+v body=%+v\n", resp.Status, i+1, url, err, resp.Body)
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}

		var map_request_payload_raw interface{}
		err = json.NewDecoder(resp.Body).Decode(&map_request_payload_raw)

		if err != nil {
			msg := fmt.Sprintf("[[EXTERNAL POST]] <<ERROR>> while unmarshalling try #%+v url=%+v err=%+v body=%+v\n", i+1, url, err, resp.Body)
			log(msg)
			time.Sleep(4 * time.Second)
			continue
		}
		failed := false

		map_request_payload := map_request_payload_raw.([]interface{})
		for _, response_map := range map_request_payload {
			response_map_typed := response_map.(map[string]interface{})
			response_code := GetFloat64FromMap(response_map_typed, "code")
			if response_code < 200 || response_code > 299 {
				msg := fmt.Sprintf("[[EXTERNAL POST]] <<FAIL STATUS CODE>> %+v try #%+v url=%+v body=%+v request=%+v\n", response_code, i+1, url, resp.Body, requests)
				log(msg)
				failed = true
				time.Sleep(4 * time.Second)
				continue
			}

		}
		if !failed {
			msg := fmt.Sprintf("[[EXTERNAL POST]] <<SUCCESS POST>> try #%+v url=%+v response=%+v requests=%+v\n", i+1, url, resp, requests)
			log(msg)
			i = 3
		}
	}
}

func GetFloat64FromMap(m map[string]interface{}, k string) float64 {
	return m[k].(float64)
}
