package main

const (
	_                = iota
	FLV_FILE_TYPE     //flv file
	Flv_HTTP_TYPE  //http-flv
	Flv_RTMP_TYPE  //rtmp
)

const (
	AUDIO_TAG  = byte(0x08)
	VIDEO_TAG  = byte(0x09)
	SCRIPT_TAG = byte(0x12)
)

//options for http flv
type HttpOptions struct {
	UserAgent string
	Referer   string
}
