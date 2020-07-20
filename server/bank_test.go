package main

import (
	"testing"
	"reflect"
)
/*
type entry struct {
	Bal	int
	SSN	uint32
} */
//type BankType map[uint32]*entry 
func TestOpenAndMod(t *testing.T) {
	var set1 = BankType{
		0         : {Bal: 200,       SSN: 0},
		3         : {Bal: -132,      SSN: 3892341},
		4         : {Bal: 11115,     SSN: 4194745990},
		6         : {Bal: 1534300814,SSN: 689607142},
		624556489 : {Bal: 1994362715,SSN: 2748336316},
	}
	var opens = []uint32 {
		3176915754,
		1135689063,
		1960039307,
		2595918420,
	}
	type modEntry struct {
		Uid	uint32
		Val	int
	}
	var mods = []modEntry {
		{0         ,-10},
		{4         ,-42368842},
		{624556489 ,0},
		{6         ,33674053},
	}
	wantFinal := BankType{
		0         : {190       ,0},
		1         : {0         ,3176915754},
		2         : {0         ,1135689063},
		3         : {-132      ,3892341},
		4         : {-42357727 ,4194745990},
		5         : {0         ,1960039307},
		6         : {1567974867,689607142},
		7         : {0         ,2595918420},
		624556489 : {1994362715,2748336316},
	}

	for i := 0; i<len(opens); i++ {//must be done in order
		Open(set1, uint32(opens[i]))
	}
	for _, entry := range mods {
		_,_ = Mod(set1, entry.Uid, entry.Val)
	}

	if !reflect.DeepEqual(wantFinal, set1) {
		t.Error("Mod/Open mismatch\nwant:")
		for k, v := range wantFinal {
			t.Error(k," : ",*v)
		}
		t.Error("\ngot:")
		for k, v := range set1 {
			t.Error(k," : ",*v)
		}
	}
}
func TestSSNCheck(t *testing.T){
	var set1 = BankType{
		0         : {Bal: 200,       SSN: 0},
		3         : {Bal: -132,      SSN: 3892341},
		4         : {Bal: 11115,     SSN: 4194745990},
		6         : {Bal: 1534300814,SSN: 689607142},
		624556489 : {Bal: 1994362715,SSN: 2748336316},
	}

	type check struct {
		Num	uint32
		Ans	bool
	}
	var SSNChecks = []check {
		{0         ,true},
		{1801916073,false},
		{1075435019,false},
		{2748336316,true},
		{2667815423,false},
	}
	for _, layer := range SSNChecks {
		res, _ := CheckSSN(set1, layer.Num)
		if res != layer.Ans {
			t.Error("SSN mismatch\nwant:\n",layer,"\ngot:\n",res)
		}
	}
}
func TestUIDCheck(t *testing.T) {
	var set1 = BankType{
		0         : {Bal: 200,       SSN: 0},
		3         : {Bal: -132,      SSN: 3892341},
		4         : {Bal: 11115,     SSN: 4194745990},
		6         : {Bal: 1534300814,SSN: 689607142},
		624556489 : {Bal: 1994362715,SSN: 2748336316},
	}

	type check struct {
		Num	uint32
		Ans	bool
	}
	var UIDChecks = []check {
		{815093642, false},
		{3,         true},
		{0,         true},
		{5,         false},
		{624556489, true},
		{3251927295,false},
	}
	for _, layer := range UIDChecks {
		res, _ := CheckUid(set1, layer.Num)
		if res != layer.Ans {
			t.Error("UID mismatch\nwant:\n",layer,"\ngot:\n",res)
		}
	}
}
