package compression

import (
	"bytes"
	"compress/gzip"
	"hash"
	"hash/crc32"
	"os"
)

type FastGzipOutputStream struct {
	writer  *gzip.Writer
	crc     hash.Hash32
	file    *os.File
	buffer  *bytes.Buffer
	tempOut []byte
}

func NewFastGzipOutputStream(filePath string, compressionLevel int, tempOutputLength int) (*FastGzipOutputStream, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	writer, err := gzip.NewWriterLevel(buffer, compressionLevel)
	if err != nil {
		return nil, err
	}

	tempOut := make([]byte, tempOutputLength)

	return &FastGzipOutputStream{
		writer:  writer,
		crc:     crc32.NewIEEE(),
		file:    file,
		buffer:  buffer,
		tempOut: tempOut,
	}, nil
}
