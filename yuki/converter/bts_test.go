package converter_test

import (
	"testing"

	"github.com/AugustineAurelius/fuufu/yuki/converter"
	"github.com/stretchr/testify/require"
)

func Test_StringToBytes(t *testing.T) {
	str := "Hello world!"
	strByte := []byte(str)

	res := converter.StringToBytes(str)

	require.Equal(t, strByte, res)
}

func Test_BytesToString(t *testing.T) {
	str := "Hello world!"
	strByte := []byte(str)

	res := converter.BytesToString(strByte)
	require.Equal(t, str, res)
}
