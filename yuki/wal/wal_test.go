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

	assert.Equal(t, uint32(61), fullLen)

	var timestamp int64
	n, err = binary.Decode(info[4:12], binary.LittleEndian, &timestamp)
	assert.Equal(t, 8, n)
	require.NoError(t, err)

	assert.Equal(t, expectedKey, hex.EncodeToString(info[12:44]))

	var valueLen uint32
	n, err = binary.Decode(info[44:48], binary.LittleEndian, &valueLen)
	assert.Equal(t, 4, n)
	require.NoError(t, err)
	assert.Equal(t, uint32(9), valueLen)

	assert.Equal(t, value, info[48:57])

	crc := crc32.NewIEEE()
	crc.Write(info[4 : len(info)-4])
	assert.Equal(t, crc.Sum(nil), info[len(info)-4:])

	os.Remove(w.FileName())
}
