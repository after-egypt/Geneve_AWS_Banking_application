package main

import (
	"testing"
	"reflect"
)

var optionTest1 = []byte{
	0xAB, 0xCD, 0x03, 0x01, 0x01, 0x23, 0x45, 0x67,
}

var optionTest2 = []byte{
	0x4d, 0x72, 0x60, 0x36, 0xdd, 0xdd, 0x85, 0x05, 0xda, 0x99, 0x4a, 0xf2,
}



func TestOptDecode(t *testing.T) {
	wantHead := OptHead{
		Class:  0xABCD,
		Type:   0x03,
		Flags:  0x0,
		Length: 0x1,
	}
	wantData := []byte{
		0x01, 0x23, 0x45, 0x67,
	}

	gotHead, gotData := DecodeOptHead(optionTest1)
	if !reflect.DeepEqual(gotHead, wantHead) {
		t.Errorf("test 1: header mismatch. Want: %#v\nGot: %#v\n", wantHead, gotHead)
	}
	if !reflect.DeepEqual(gotData, wantData) {
		t.Errorf("test 1: data mismatch. Want: %#v\nGot: %#v\n", wantData, gotData)
	}
}

func TestOptEncode(t *testing.T) {
	give := OptHead{
		Class:	0x1242,
		Type:	0xA0,
		Flags:	0x0,
		Length: 0x13,
	}
	var want = []byte{
		0x12, 0x42, 0xA0, 0x13,
	}
	got := EncodeOptHead(give)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Header Mismatch. Want: %#v\nGot: %#v\n", want, got)
	}
}

