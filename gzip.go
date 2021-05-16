package kaban

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// gzip
func compressGZIP(blob []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	if _, err := zw.Write(blob); err != nil {
		return nil, fmt.Errorf("gzip.Writer.Write() %s", err.Error())
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("gzip.Writer.Close() %s", err.Error())
	}
	return buf.Bytes(), nil
}

// gunzip
func decompressGZIP(blob []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	gzipBuf := bytes.NewReader(blob)
	zr, err := gzip.NewReader(gzipBuf)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewReader() %s", err.Error())
	}
	if _, err := io.Copy(buf, zr); err != nil {
		return nil, fmt.Errorf("io.Copy() %s", err.Error())
	}
	if err := zr.Close(); err != nil {
		return nil, fmt.Errorf("gzip.Reader.Close() %s", err.Error())
	}
	return buf.Bytes(), nil
}
