package util

import (
	"testing"
)

func TestGenerateStr(t *testing.T) {
	s := GenerateStr()
	if len(s) != 32 {
		t.Fatal()
	}
}

type test struct {
}

func TestGetArrFirstItem(t *testing.T) {
	var arr1 []*test = nil
	item := GetArrFirstItem(arr1)
	if item != nil {
		t.Fatal()
	}
	var arr2 = make([]*test, 0)
	item = GetArrFirstItem(arr2)
	if item != nil {
		t.Fatal()
	}
	var arr3 = make([]*test, 0)
	arr3 = append(arr3, &test{})
	item = GetArrFirstItem(arr3)
	if item == nil {
		t.Fatal()
	}
}
