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

func TestGenerate16Str(t *testing.T) {
	str := Generate16Str()
	t.Log(str)
	if len(str) != 16 {
		t.Fatal()
	}
}

func TestGenerate24Str(t *testing.T) {
	str := Generate24Str()
	t.Log(str)
	if len(str) != 24 {
		t.Fatal()
	}
}

func TestGenerate32Str(t *testing.T) {
	str := Generate32Str()
	t.Log(str)
	if len(str) != 32 {
		t.Fatal()
	}
}
