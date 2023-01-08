package handlers

type RequestMessage struct {
	EventType string     `json:"event_type"`
	Id        string     `json:"id"`
	InputId   string     `json:"input_id"`
	InputType string     `json:"input_type"`
	RTMP      *RTMP      `json:"rtmp"`
	Mpegts    *Mpegts    `json:"mpegts"`
	MediaInfo *MediaInfo `json:"media_info"`
}

type ResponseMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}
