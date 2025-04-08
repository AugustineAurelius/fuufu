package bloom

import (
	"encoding/binary"
	"hash/fnv"
)

type filter struct {
	bits *bitSet
	m    uint64
	k    int
}

func NewFilter(bits *bitSet, m uint64, k int) *filter {
	return &filter{bits, m, k}
}

func (f *filter) Add(data []byte) {
	h := fnv.New64()
	buf := make([]byte, 8)

	for i := range f.k {
		h.Reset()
		binary.LittleEndian.PutUint64(buf, uint64(i))
		h.Write(buf)
		h.Write(data)
		f.bits.SetOneOn(int(h.Sum64() % f.m))
	}
}

func (f *filter) Test(data []byte) bool {
	h := fnv.New64()

	buf := make([]byte, 8)
	for i := range f.k {
		h.Reset()
		binary.LittleEndian.PutUint64(buf, uint64(i))
		h.Write(buf)
		h.Write(data)
		if !f.bits.IsOne(int(h.Sum64() % f.m)) {
			return false
		}
	}
	return true
}
