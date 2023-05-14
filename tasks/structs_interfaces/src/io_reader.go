package src

import (
	"io"
	"strings"
)

type Reader interface {
	Read(p []byte) (int, error)
	ReadAll(bufSize int) (string, error)
	BytesRead() int64
}

type CountingToLowerReaderImpl struct {
	Reader         io.Reader
	TotalBytesRead int64
}

func NewCountingReader(r io.Reader) *CountingToLowerReaderImpl {
	return &CountingToLowerReaderImpl{
		Reader: r,
	}
}

func (cr *CountingToLowerReaderImpl) Read(p []byte) (int, error) {
	readBytes, err := cr.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	cr.TotalBytesRead += int64(readBytes)
	convertBytesToLowerCase(p)
	return readBytes, nil
}

func (cr *CountingToLowerReaderImpl) ReadAll(bufSize int) (string, error) {
	stringBuffer := make([]byte, bufSize)
	readBytes, err := cr.Reader.Read(stringBuffer)
	builder := strings.Builder{}
	for ; err == nil; readBytes, err = cr.Reader.Read(stringBuffer) {
		builder.Write(stringBuffer[:readBytes])
	}
	transformedString := strings.ToLower(builder.String())
	return transformedString, nil
}

func (cr *CountingToLowerReaderImpl) BytesRead() int64 {
	return cr.TotalBytesRead
}

func convertBytesToLowerCase(bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		if bytes[i] >= 'A' && bytes[i] <= 'Z' {
			bytes[i] = bytes[i] + 32
		}
	}
}
