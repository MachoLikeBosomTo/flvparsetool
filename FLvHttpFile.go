package main

import (
	"errors"
	"encoding/binary"
	"log"
	"fmt"
	"net/http"
)

type FlvHttpFileParse struct {
	ctx     FlvContext
	httpReq *http.Response
}

func (parser *FlvHttpFileParse) OpenFlv(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}

	parser.httpReq, err = http.Get(name)
	if err != nil {
		return err
	}

	parser.ctx.fileName = name

	return err
}

func (parser *FlvHttpFileParse) CloseFlv() (err error) {
	if parser.httpReq != nil {
		err = parser.httpReq.Body.Close()
	}

	return err
}

func (parser *FlvHttpFileParse) ReadFlvHeader() (header *FlvHeader, err error) {
	//body, err := ioutil.Read(parser.httpReq.Body)

	err = binary.Read(parser.httpReq.Body, binary.BigEndian, &parser.ctx.flvHeader)
	header = &parser.ctx.flvHeader
	log.Printf("%#v\n", parser.ctx.flvHeader)

	return header, err
}

func (parser *FlvHttpFileParse) ReadFlvTag() (flvTag *FlvTag, err error) {
	var flvTagHeaderData FlvTagHeaderData

	err = binary.Read(parser.httpReq.Body, binary.BigEndian, &flvTagHeaderData)
	if err != nil {
		return flvTag, err
	}
	ParseFlvTagHeaderData(&flvTagHeaderData, &parser.ctx.flvTag.TagHeader)
	log.Printf("%+v\n", parser.ctx.flvTag.TagHeader)

	parser.ctx.flvTag.TagData = make([]byte, parser.ctx.flvTag.TagHeader.DataSize)
	err = binary.Read(parser.httpReq.Body, binary.BigEndian, &parser.ctx.flvTag.TagData)
	if err != nil {
		return flvTag, err
	}
	err = binary.Read(parser.httpReq.Body, binary.BigEndian, &parser.ctx.flvTag.PreTagSize)
	if err != nil {
		return flvTag, err
	}

	if parser.ctx.flvTag.TagHeader.DataSize+11 != parser.ctx.flvTag.PreTagSize {
		err = fmt.Errorf("PreTagSize:%d mismatch DataSize:%d add 11", parser.ctx.flvTag.PreTagSize, parser.ctx.flvTag.TagHeader.DataSize)
		return flvTag, err
	}
	//log.Printf("%+v\n", parser.ctx.flvTag.PreTagSize)
	flvTag = &parser.ctx.flvTag

	return flvTag, err
}
