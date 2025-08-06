package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GeneratePassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
