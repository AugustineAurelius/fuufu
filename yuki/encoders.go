package yuki

import (
	"bytes"
	"compress/flate"
	"encoding/gob"
	"errors"
	"io"
)

var ErrZeroBytes = errors.New("zero bytes writed")

func gobEncode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecode(data []byte, e any) error {
	if len(data) == 0 {
		return nil
	}

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(&e)
}

func flateEncode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(data)
	if n == 0 {
		return nil, ErrZeroBytes
	}
	if err != nil {
		return nil, err
	}
	defer w.Close()

	w.Flush()

	return buf.Bytes(), nil
}

func flateDecode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	decoded, err := io.ReadAll(r)
	if err != nil && errors.Is(err, io.EOF) {
		return nil, err
	}
	return decoded, nil
}
