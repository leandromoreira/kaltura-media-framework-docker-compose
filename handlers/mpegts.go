package handlers

type Mpegts struct {
	Addr       *string `json:"addr"`
	Connection *int64  `json:"connection"`
	Index      *int64  `json:"index"`
	Pid        *int64  `json:"pid"`
	ProgNum    *int64  `json:"prog_num"`
	StreamID   *string `json:"stream_id"`
}
