package main

type FlvHeader struct {
	Signature  [3]byte // always "FLV" (0x46 0x4C 0x56)
	Version    byte    // 0x01
	AVTag      byte    //
	Offset     uint32
	PreTagSize uint32
}

type FlvTagHeaderData struct {
	TagType      byte
	DataSize     [3]byte
	Timestamp    [3]byte
	TimestampExt byte
	StreamId     [3]byte
}

type FlvTagHeader struct {
	TagType   byte
	DataSize  uint32
	Timestamp uint32
	StreamId  uint32
}

type FlvTag struct {
	TagHeader  FlvTagHeader
	TagData    []byte
	PreTagSize uint32
}

/*
func (flvTag * FlvTag)GetFlvTagType()(byte){
	return flvTag.TagHeader.TagType
}

func (flvTag * FlvTag)GetFlvTagDataSize()(byte){
	return flvTag.TagHeader.TagType
}
*/

type FlvContext struct {
	fileName       string
	fileType       int
	currentAudioTs int64
	currentVideoTs int64
	preAudioTs     int64
	preVideoTs     int64
	parsedSize     uint64
	parsedTags     uint64
	flvHeader      FlvHeader
	flvTag         FlvTag
}

type FlvParse interface {
	OpenFlv(name string) (err error)
	CloseFlv() (err error)
	ReadFlvHeader() (header *FlvHeader, err error)
	ReadFlvTag() (tagHeader *FlvTag, err error)
}

type FlvPack interface {
	OpenFlv(name string) (err error)
	CloseFlv() (err error)
	WriteFlvHeader(header *FlvHeader) (err error)
	WriteFlvTag(flvTag *FlvTag) (err error)
}

var (
	_ FlvParse = &FlvFileParse{}
	_ FlvParse = &FlvHttpFileParse{}
	//_ FlvParse = FlvStreamRtmpParse{}
)

var (
	_ FlvPack = &FlvFilePack{}
)

func ParseFlvTagHeaderData(data *FlvTagHeaderData, header *FlvTagHeader) {
	header.TagType = data.TagType
	header.DataSize = uint32(data.DataSize[0])<<16 | uint32(data.DataSize[1])<<8 | uint32(data.DataSize[2])
	header.Timestamp = uint32(data.TimestampExt)<<32 | uint32(data.Timestamp[0])<<16 | uint32(data.Timestamp[1])<<8 | uint32(data.Timestamp[2])
	header.StreamId = uint32(data.StreamId[0])<<16 | uint32(data.StreamId[1])<<8 | uint32(data.StreamId[2])
}

func PackFlvTagHeaderData(header *FlvTagHeader, data *FlvTagHeaderData) {
	data.TagType = header.TagType
	data.DataSize[0] = byte(header.DataSize >> 16 & 0xff)
	data.DataSize[1] = byte(header.DataSize >> 8 & 0xff)
	data.DataSize[2] = byte(header.DataSize & 0xff)

	data.TimestampExt = byte(header.Timestamp >> 32 & 0xff)
	data.Timestamp[0] = byte(header.Timestamp >> 16 & 0xff)
	data.Timestamp[1] = byte(header.Timestamp >> 8 & 0xff)
	data.Timestamp[2] = byte(header.Timestamp & 0xff)

	data.StreamId[0] = byte(header.StreamId >> 16 & 0xff)
	data.StreamId[1] = byte(header.StreamId >> 16 & 0xff)
	data.StreamId[2] = byte(header.StreamId >> 16 & 0xff)
}
