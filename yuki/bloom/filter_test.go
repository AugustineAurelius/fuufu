package bloom_test

import (
	"testing"

	"github.com/AugustineAurelius/fuufu/yuki/bloom"
	"github.com/stretchr/testify/require"
)

func Test_Filter(t *testing.T) {
	filter := bloom.NewFilter(bloom.NewBitSet(), 1_000_000, 64)
	data := []byte("test data")
	filter.Add(data)

	require.True(t, filter.Test(data))
	require.False(t, filter.Test([]byte("test dato")))
}
