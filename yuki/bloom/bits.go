package bloom

import "math/big"

type bitSet struct {
	bits big.Int
}

func NewBitSet() *bitSet {
	return &bitSet{}
}

func (bs *bitSet) SetOneOn(index int) {
	bs.bits.SetBit(&bs.bits, index, 1)
}

func (bs *bitSet) IsOne(index int) bool {
	return bs.bits.Bit(index) != 0
}
