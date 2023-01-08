package handlers

type RequestMessage struct {
	EventType string `json:"event_type"`
	Id        string `json:"id"`
}

type ResponseMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}
