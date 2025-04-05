package wal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"time"
)

type Wal struct {
	f            *os.File
	consumer     chan []byte
	workerAmount int
}

var defaultWorkerAmount = 3

type WalOption func(w *Wal)

// |record len (4 bytes)| hashed key (32 bytes) | compressed value len (4 bytes) | value... | crc32 (32 bytes)
func OpenWAL(opts ...WalOption) (*Wal, error) {
	//TODO: augustine_aurelius may be os.O_SYNC
	f, err := os.OpenFile(time.Now().UTC().Format("20060102150405")+".txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	w := Wal{
		f: f,

		consumer:     make(chan []byte, 1),
		workerAmount: defaultWorkerAmount,
	}

	for _, opt := range opts {
		opt(&w)
	}

	for range w.workerAmount {
		go w.write()
	}

	return &w, nil
}

func (w *Wal) Close() error {
	return w.f.Close()
}

func (w *Wal) write() {
	for data := range w.consumer {
		_, err := w.f.Write(data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (w *Wal) Add(key, value []byte) {
	var buf bytes.Buffer

	crc := crc32.NewIEEE()

	valueLen := len(value)
	walLen := 32 + 32 + uint32(valueLen) + 4

	if err := binary.Write(&buf, binary.LittleEndian, walLen); err != nil {
		fmt.Println(err)
	}

	buf.Write(key)
	if err := binary.Write(&buf, binary.LittleEndian, uint32(valueLen)); err != nil {
		fmt.Println(err)
	}
	buf.Write(value)

	crc.Write(buf.Bytes())
	buf.Write(crc.Sum(nil))

	w.consumer <- buf.Bytes()
}

func (w *Wal) FileName() string {
	return w.f.Name()
}

func WithWorkerAmount(workerAmount int) WalOption {
	if workerAmount <= 0 {
		workerAmount = defaultWorkerAmount
	}
	return func(w *Wal) {
		w.workerAmount = workerAmount
	}
}
