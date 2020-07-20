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
	wantHead := optHead{
		Class:  0xABCD,
		Type:   0x03,
		Flags:  0x0,
		Length: 0x1,
	}
	wantData := []byte{
		0x01, 0x23, 0x45, 0x67,
	}

	gotHead, gotData := decodeOptHead(optionTest1)
	if !reflect.DeepEqual(gotHead, wantHead) {
		t.Errorf("decodeOptHead test 1: header mismatch. Want: %#v\nGot: %#v\n", wantHead, gotHead)
	}
	if !reflect.DeepEqual(gotData, wantData) {
		t.Errorf("decodeOptHead test 1: data mismatch. Want: %#v\nGot: %#v\n", wantData, gotData)
	}
}

