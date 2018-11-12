package main

import (
	"errors"
	"os"
	"encoding/binary"
	"bufio"
	"log"
	"fmt"
)

type FlvFileParse struct {
	ctx    FlvContext
	file   *os.File
	reader *bufio.Reader
}

type FlvFilePack struct {
	file   *os.File
	writer *bufio.Writer
}

func (parser *FlvFileParse) OpenFlv(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}

	parser.file, err = os.Open(name)
	if err != nil {
		return err
	}

	parser.reader = bufio.NewReader(parser.file)
	parser.ctx.fileName = name

	return err
}

func (parser *FlvFileParse) CloseFlv() (err error) {
	if parser.file != nil {
		err = parser.file.Close()
	}

	return err
}

func (parser *FlvFileParse) ReadFlvHeader() (header *FlvHeader, err error) {
	err = binary.Read(parser.reader, binary.BigEndian, &parser.ctx.flvHeader)
	header = &parser.ctx.flvHeader
	log.Printf("%#v\n", parser.ctx.flvHeader)

	return header, err
}

func (parser *FlvFileParse) ReadFlvTag() (flvTag *FlvTag, err error) {
	var flvTagHeaderData FlvTagHeaderData

	err = binary.Read(parser.reader, binary.BigEndian, &flvTagHeaderData)
	if err != nil {
		return flvTag, err
	}
	ParseFlvTagHeaderData(&flvTagHeaderData, &parser.ctx.flvTag.TagHeader)
	log.Printf("%+v\n", parser.ctx.flvTag.TagHeader)

	parser.ctx.flvTag.TagData = make([]byte, parser.ctx.flvTag.TagHeader.DataSize)
	err = binary.Read(parser.reader, binary.BigEndian, &parser.ctx.flvTag.TagData)
	if err != nil {
		return flvTag, err
	}
	err = binary.Read(parser.reader, binary.BigEndian, &parser.ctx.flvTag.PreTagSize)
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

func (pack *FlvFilePack) OpenFlv(name string) (err error) {
	if name == "" {
		return errors.New("empty name")
	}

	pack.file, err = os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	pack.writer = bufio.NewWriter(pack.file)

	return err
}

func (pack *FlvFilePack) CloseFlv() (err error) {
	if pack.file != nil {
		err = pack.file.Close()
	}

	return err
}

func (pack *FlvFilePack) WriteFlvHeader(header *FlvHeader) (err error) {
	err = binary.Write(pack.writer, binary.BigEndian, header)

	return err
}

func (pack *FlvFilePack) WriteFlvTag(flvTag *FlvTag) (err error) {
	var flvTagHeaderData FlvTagHeaderData

	PackFlvTagHeaderData(&flvTag.TagHeader, &flvTagHeaderData)
	err = binary.Write(pack.writer, binary.BigEndian, &flvTagHeaderData)
	if err != nil {
		return err
	}

	err = binary.Write(pack.writer, binary.BigEndian, flvTag.TagData)
	if err != nil {
		return err
	}

	err = binary.Write(pack.writer, binary.BigEndian, flvTag.PreTagSize)
	if err != nil {
		return err
	}

	return err
}
