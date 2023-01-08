package handlers

type RTMP struct {
	Addr       *string `json:"addr"`
	App        *string `json:"app"`
	Args       *string `json:"args"`
	Connection *int64  `json:"connection"`
	Flashver   *string `json:"flashver"`
	Name       *string `json:"name"`
	PageURL    *string `json:"page_url"`
	SwfURL     *string `json:"swf_url"`
	TcURL      *string `json:"tc_url"`
	Type       *string `json:"type"`
}
