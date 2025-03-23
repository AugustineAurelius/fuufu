package cache_test

import (
	"runtime"
	"testing"

	"github.com/AugustineAurelius/fuufu/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func Test_WeakCache(t *testing.T) {
	cache := cache.NewWeak[string, string]()

	data := "cached data"
	cache.Set("key1", data)

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.NotNil(t, val)

	runtime.GC()

	val, ok = cache.Get("key1")
	assert.False(t, ok)
	assert.Nil(t, val)

}
