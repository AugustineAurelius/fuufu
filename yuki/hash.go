package yuki

import "crypto/sha256"

func hashKey(key, hashedKey []byte) {
	hash := sha256.New()
	hash.Write(key)
	hash.Sum(hashedKey[:0])
}
