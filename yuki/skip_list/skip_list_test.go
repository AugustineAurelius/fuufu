package skiplist_test

import (
	"strconv"
	"testing"

	skiplist "github.com/AugustineAurelius/fuufu/yuki/skip_list"
	"github.com/stretchr/testify/assert"
)

func Test_Put(t *testing.T) {
	sl := skiplist.New()
	for i := range 101 {
		sl.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)))
	}

	for i := range 101 {
		value, ok := sl.Get([]byte(strconv.Itoa(i)))
		assert.True(t, ok)
		assert.Equal(t, []byte(strconv.Itoa(i)), value)
	}

	i := 0
	iterator := sl.NewIterator()
	for iterator.Valid() {
		iterator.Next()
		i++
	}

	assert.Equal(t, 101, i)
}
