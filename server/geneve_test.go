package main

import(
	"reflect"
	"testing"
)

func TestOpt(t *testing.T) {
	formatted := []OptHead{{
			Class:	0x0560,
			Type:	0x53,
			Flags:	0x0,
			Length:	0x1,
		}, {
			Class:	0x0f03,
			Type:	0xc9,
			Flags:	0x4,
			Length:	0x1a,
		},
	}
	raw := [][]byte{
		{0x05, 0x60, 0x53, 0x01},
		{0x0f, 0x03, 0xc9, 0x9a},
	}

	for i:=0; i<len(raw); i++ {
		dGot := DecodeOptHead(raw[i])

		if !reflect.DeepEqual(dGot, formatted[i]) {
			t.Errorf("Decode mismatch:\nExpected:\n%#v\nGot:\n%#v", formatted[i],dGot)
		}

		eGot := EncodeOptHead(formatted[i])
		if !reflect.DeepEqual(eGot, raw[i]) {
			t.Errorf("Encode mismatch:\nExpected:\n%#v\nGot:\n%#v",raw[i],eGot)
		}
	}
}
func TestClientID(t *testing.T) {
	formatted := []ClientID{
		{0x6fa2bb15, 0xa82d09ad},
		{0xfe3b6d3f, 0xd33512b1},
	}
	raw := [][]byte{
		{0x6f, 0xa2, 0xbb, 0x15, 0xa8, 0x2d, 0x09, 0xad},
		{0xfe, 0x3b, 0x6d, 0x3f, 0xd3, 0x35, 0x12, 0xb1},
	}

	for i:=0; i<len(raw); i++ {
		dGot := DecodeClientID(raw[i])
		if !reflect.DeepEqual(dGot, formatted[i]) {
			t.Errorf("Decode mismatch:\nExpected:\n%#v\nGot:\n%#v", formatted[i],dGot)
		}

		eGot := EncodeClientID(formatted[i])
		if !reflect.DeepEqual(eGot, raw[i]) {
			t.Errorf("Encode mismatch:\nExpected:\n%#v\nGot:\n%#v",raw[i],eGot)
		}
	}
}
func TestEncodeNewAcc(t *testing.T) {
	formatted := [][]uint32{
		{0x7e24edc0, 0xc586594d, 0xa9fd9ec6},
		{0x97177d5a, 0x5a961cb9, 0xe02eee01},
	}
	raw := [][]byte{
		{0x00, 0x00, 0x02, 0x03, 0x7e, 0x24, 0xed, 0xc0, 0xc5, 0x86, 0x59, 0x4d, 0xa9, 0xfd, 0x9e, 0xc6},
		{0x00, 0x00, 0x02, 0x03, 0x97, 0x17, 0x7d, 0x5a, 0x5a, 0x96, 0x1c, 0xb9, 0xe0, 0x2e, 0xee, 0x01},
	}

	for i:=0; i<len(raw); i++ {
		eGot,_ := EncodeNewAcc(formatted[i][0],formatted[i][1],formatted[i][2])
		if !reflect.DeepEqual(eGot, raw[i]) {
			t.Errorf("Encode mismatch:\nExpected:\n%#v\nGot:\n%#v",raw[i],eGot)
		}
	}
}
func TestEncodeModBal(t *testing.T) {
	type entry struct {
		Cid	ClientID
		ErrCode	uint16
		Bal	int
	}
	formatted := []entry{
		{ClientID{0xdaf7f9a3, 0x3344ad2f}, 0x9d9d, -0x049c},
		{ClientID{0x049c7e74, 0xa854cd71}, 0x282b, 0xcd7110},
	}

	raw := [][]byte{
		{0x00, 0x00, 0x03, 0x03, 0xda, 0xf7, 0xf9, 0xa3, 0x33, 0x44, 0xad, 0x2f, 0x9d, 0x9d, 0x84, 0x9c},
		{0x00, 0x00, 0x03, 0x04, 0x04, 0x9c, 0x7e, 0x74, 0xa8, 0x54, 0xcd, 0x71, 0x28, 0x2b, 0x00, 0xcd, 0x71, 0x10, 0x00, 0x00},
	}
	

	for i:=0; i<len(raw); i++ {
		eGot,_ := EncodeModBal(formatted[i].Cid,formatted[i].ErrCode,formatted[i].Bal)
		if !reflect.DeepEqual(eGot, raw[i]) {
			t.Errorf("Encode mismatch:\nExpected:\n%#v\nGot:\n%#v",raw[i],eGot)
		}
	}

}
func TestDecodeOpenAcc(t *testing.T) {
	type entry struct{
		Ssn	uint32
		SeqNum	uint32
		Name	string
	}
	formatted := []entry {
		{0x8c493710, 0xfabc342c, ""},
		{0xb630817c, 0x42f2e40f, "ab12:)" },
	}
	raw := [][]byte{
		{0x8c, 0x49, 0x37, 0x10, 0xfa, 0xbc, 0x34, 0x2c},
		{0xb6, 0x30, 0x81, 0x7c, 0x42, 0xf2, 0xe4, 0x0f, 0x61, 0x62, 0x31, 0x32, 0x3a, 0x29},
	}
	for i:=0; i<len(raw); i++ {
		var dGot entry
		dGot.Ssn, dGot.SeqNum, dGot.Name = DecodeOpenAcc(raw[i])
		if !reflect.DeepEqual(dGot, formatted[i]) {
			t.Errorf("Decode mismatch:\nExpected:\n%#v\nGot:\n%#v",formatted[i],dGot)
		}
	}
}
func TestDecodeModBal(t *testing.T) {
	formatted := []int {
		-0x9d,
		0x51c328e2fabc,
	}
	raw := [][]byte{
		{0x80, 0x00, 0x00, 0x9d},
		{0x00, 0x00, 0x51, 0xc3, 0x28, 0xe2, 0xfa, 0xbc},
	}
	for i:=0; i<len(raw); i++ {
		dGot,_ := DecodeModBal(raw[i])

		if !reflect.DeepEqual(dGot, formatted[i]) {
			t.Errorf("Decode mismatch:\nExpected:\n%#v\nGot:\n%#v",formatted[i],dGot)
		}
	}
}
