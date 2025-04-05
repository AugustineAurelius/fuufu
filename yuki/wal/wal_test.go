package wal_test

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"
	"os"
	"testing"
	"time"

	"github.com/AugustineAurelius/fuufu/yuki/wal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WAL(t *testing.T) {
	w, err := wal.OpenWAL()
	require.NoError(t, err)

	value := []byte("123123123")
	hashedKey := sha256.Sum256([]byte("123"))
	w.Add(hashedKey[:], value)
	expectedKey := hex.EncodeToString(hashedKey[:])

	time.Sleep(time.Second)

	require.NoError(t, w.Close())

	info, err := os.ReadFile(w.FileName())
	require.NoError(t, err)

	var fullLen uint32
	n, err := binary.Decode(info[:4], binary.LittleEndian, &fullLen)
	assert.Equal(t, 4, n)
	require.NoError(t, err)

	assert.Equal(t, uint32(77), fullLen)
	assert.Equal(t, expectedKey, hex.EncodeToString(info[4:36]))

	var valueLen uint32
	n, err = binary.Decode(info[36:40], binary.LittleEndian, &valueLen)
	assert.Equal(t, 4, n)
	require.NoError(t, err)
	assert.Equal(t, uint32(9), valueLen)

	assert.Equal(t, value, info[40:49])

	crc := crc32.NewIEEE()
	crc.Write(info[:49])
	assert.Equal(t, crc.Sum(nil), info[49:53])

	os.Remove(w.FileName())
}
