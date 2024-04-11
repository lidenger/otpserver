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

func Generate32Str() string {
	return GenerateStr()
}

func Generate24Str() string {
	str := Generate32Str()
	return str[0:24]
}

func Generate16Str() string {
	str := Generate32Str()
	return str[0:16]
}

// GetArrFirstItem 获取数组第一个元素
func GetArrFirstItem[T any](arr []*T) *T {
	if len(arr) > 0 {
		return arr[0]
	} else {
		return nil
	}
}

func Eqs(target string, params ...string) bool {
	for _, param := range params {
		if param == target {
			return true
		}
	}
	return false
}
