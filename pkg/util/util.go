package util

import (
	"github.com/google/uuid"
	"strings"
)

// GenerateStr 生成32随机字符串
func GenerateStr() string {
	str, _ := uuid.NewUUID()
	return strings.ReplaceAll(str.String(), "-", "")
}
