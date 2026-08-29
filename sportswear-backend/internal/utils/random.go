package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString 生成随机字符串（用于 API Key、Secret 等）
func GenerateRandomString(length int) string {
	bytes := make([]byte, (length+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}
