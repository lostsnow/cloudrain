package charset

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// Converts given Reader to a UTF-8 bytes
func DecodeReader(s io.Reader, enc string) ([]byte, error) {
	reader, err := charset.NewReaderLabel(enc, s)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

// Converts given bytes to a UTF-8 bytes
func Decode(s []byte, enc string) ([]byte, error) {
	return DecodeReader(bytes.NewReader(s), enc)
}

// Converts a Reader to bytes encoded with given encoding
func EncodeReader(s io.Reader, enc string) ([]byte, error) {
	e, _ := charset.Lookup(enc)
	if e == nil {
		return nil, fmt.Errorf("unsupported charset: %q", enc)
	}
	var buf bytes.Buffer
	writer := transform.NewWriter(&buf, e.NewEncoder())
	_, err := io.Copy(writer, s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Converts a bytes to a given encoding bytes
func Encode(s []byte, enc string) ([]byte, error) {
	return EncodeReader(bytes.NewReader(s), enc)
}
