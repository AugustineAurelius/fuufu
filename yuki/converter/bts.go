package converter

import "unsafe"

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	if len(b) < 1 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}
