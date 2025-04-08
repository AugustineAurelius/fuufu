package wal

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"os"
	"time"
)

type Wal struct {
	f *os.File
}

type WalOption func(w *Wal)

// |record len (4 bytes) | timestamp (8 bytes) | hashed key (32 bytes) | compressed value len (4 bytes) | value... | crc32 (4 bytes)
func OpenWAL(opts ...WalOption) (*Wal, error) {
	f, err := os.OpenFile(time.Now().UTC().Format("20060102150405")+".txt", os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}
	w := Wal{
		f: f,
	}

	for _, opt := range opts {
		opt(&w)
	}

	return &w, nil
}

func (w *Wal) Close() error {
	return w.f.Close()
}

func (w *Wal) Add(key, value []byte) error {
	var buf bytes.Buffer

	crc := crc32.NewIEEE()

	valueLen := len(value)
	walLen := 4 + 32 + 8 + 4 + uint32(valueLen) + 4

	if err := binary.Write(&buf, binary.LittleEndian, walLen); err != nil {
		return err
	}

	if err := binary.Write(&buf, binary.LittleEndian, time.Now().UTC().UnixMicro()); err != nil {
		return err
	}

	buf.Write(key)

	if err := binary.Write(&buf, binary.LittleEndian, uint32(valueLen)); err != nil {
		return err
	}
	buf.Write(value)

	crc.Write(buf.Bytes())
	buf.Write(crc.Sum(nil))

	if _, err := w.f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (w *Wal) FileName() string {
	return w.f.Name()
}
