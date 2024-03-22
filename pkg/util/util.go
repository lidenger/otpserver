package util

import (
	"github.com/google/uuid"
	"strings"
)

// GenerateStr 生成32位随机字符串
func GenerateStr() string {
	str, _ := uuid.NewUUID()
	return strings.ReplaceAll(str.String(), "-", "")
}

// GetArrFirstItem 获取数组第一个元素
func GetArrFirstItem[T any](arr []*T) *T {
	if len(arr) > 0 {
		return arr[0]
	} else {
		return nil
	}
}
