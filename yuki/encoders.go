package yuki

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
)

var ErrZeroBytes = errors.New("zero bytes writed")

func flateEncode(data []byte, buf *bytes.Buffer) error {
	if len(data) == 0 {
		return nil
	}

	w, err := flate.NewWriter(buf, flate.DefaultCompression)
	if err != nil {
		return err
	}
	defer w.Close()

	n, err := w.Write(data)
	switch {
	case n == 0:
		return ErrZeroBytes
	case err != nil:
		return err
	default:
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
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
