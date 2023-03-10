package handlers

type MediaInfo struct {
	Bitrate   *int64  `json:"bitrate"`
	CodecID   *int64  `json:"codec_id"`
	ExtraData *string `json:"extra_data"`
	MediaType *string `json:"media_type"`
	MediaInfoVideo
	MediaInfoAudio
}

type MediaInfoVideo struct {
	FrameRate   *float64 `json:"frame_rate"`
	Height      *int64   `json:"height"`
	Width       *int64   `json:"width"`
	CeaCaptions bool     `json:"cea_captions"`
}

type MediaInfoAudio struct {
	Channels      *int64  `json:"channels"`
	ChannelLayout *string `json:"channel_layout"`
	BitsPerSample *int64  `json:"bits_per_sample"`
	SampleRate    *int64  `json:"sample_rate"`
}
