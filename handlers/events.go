package handlers

type RequestMessage struct {
	EventType string `json:"event_type"`
}

type ResponseMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
