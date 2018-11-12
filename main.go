package main

import (
	"log"
	"flag"
	"strings"
)

var (
	input          = flag.String("i", "", "set input")
	output         = flag.String("o", "", "set output")
	audioTimeDelta = flag.Uint("atd", 0, "set audio time delta")
	videoTimeDelta = flag.Uint("vtd", 0, "set video time delta")
	tagCountLimit  = flag.Uint("tc", 0, "set tag count")
)

func init() {
	flag.Parse()
}

func main() {
	if *input == "" {
		flag.Usage()
		return
	}

	var flvType = FLV_FILE_TYPE
	if strings.Contains(*input, "http") {
		flvType = Flv_HTTP_TYPE
	} else if strings.Contains(*input, "rtmp") {
		flvType = Flv_RTMP_TYPE
	}

	var parser FlvParse

	switch flvType {
	case FLV_FILE_TYPE:
		parser = new(FlvFileParse)
	case Flv_HTTP_TYPE:
		parser = new(FlvHttpFileParse)
	case Flv_RTMP_TYPE:
		//parser = new(FlvHttpFileParse)
		return
	}

	err := parser.OpenFlv(*input)
	if err != nil {
		log.Println(err)
		return
	}

	var outFlag = false
	var packer FlvPack
	if *output != "" && (*audioTimeDelta != 0 || *videoTimeDelta != 0) {
		outFlag = true
		packer = new(FlvFilePack)
		packer.OpenFlv(*output)
		if err != nil {
			log.Println(err)
			return
		}
	}

	header, err := parser.ReadFlvHeader()
	if err != nil {
		log.Println(err)
		return
	}

	if outFlag {
		err = packer.WriteFlvHeader(header)
		if err != nil {
			log.Println(err)
			return
		}
	}

	var tagCount uint = 0
	for {
		tag, err := parser.ReadFlvTag()
		if err != nil {
			log.Println(err)
			break
		}
		if outFlag {
			if *audioTimeDelta != 0 && tag.TagHeader.TagType == AUDIO_TAG {
				tag.TagHeader.Timestamp = tag.TagHeader.Timestamp + uint32(*audioTimeDelta)
			}
			if *videoTimeDelta != 0 && tag.TagHeader.TagType == VIDEO_TAG {
				tag.TagHeader.Timestamp = tag.TagHeader.Timestamp + uint32(*videoTimeDelta)
			}

			err = packer.WriteFlvTag(tag)
			if err != nil {
				log.Println(err)
				return
			}
		}

		tagCount++
		if tagCount >= *tagCountLimit && *tagCountLimit != 0{
			log.Printf("tagCount:%d ge tagCountLimit:%d", tagCount, *tagCountLimit)
			break
		}
		//break
	}

	parser.CloseFlv()

	if outFlag {
		packer.CloseFlv()
	}

	return
}
