package bloom_test

import (
	"testing"

	"github.com/AugustineAurelius/fuufu/yuki/bloom"
	"github.com/stretchr/testify/require"
)

func Test_Bits(t *testing.T) {

	bits := bloom.NewBitSet()
	for i := 1000; i < 2000; i++ {
		bits.SetOneOn(i)

	}
	for i := 0; i < 10000; i++ {
		if i > 999 && i < 2000 {
			require.True(t, bits.IsOne(i))
		}
	}
}
