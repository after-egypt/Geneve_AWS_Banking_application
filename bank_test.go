package main

import (
	"testing"
	"reflect"
)

var m1 = map[uint32]int{
	0 : 200,
	2 : 432,
	3 : 0,
	4 : -123,
}
var m2 = map[uint32]int{
	1 : 2,
	2 : 9019234,
	3 : 3832,
	7 : 254,
}
var m3 = map[uint32]int{
	0 : 971,
	2 : 1678,
	3 : 11115,
	19123 : 48052985,
}
/*func TestInit(t *testing.T) {
	m := Init()
} */
func TestOpen1(t *testing.T) {
	var uids = []uint32{}
	var errs = []error{}
	for i:=0; i<3;i++ {
		uid, err :=Open(m1)
		uids = append(uids, uid)
		errs = append(errs, err)
	}
	var m1Want = map[uint32]int{
		0 : 200,
		1 : 0,
		2 : 432,
		3 : 0,
		4 : -123,
		5 : 0,
		6 : 0,
	}
	var uidsWant = []uint32{
		1,5,6,
	}
	var errsWant = []error{
		nil,nil,nil,
	}
	if !reflect.DeepEqual(m1,m1Want) {
		t.Error("Database mismatch\nwant:",m1Want,"\ngot:",m1)
	}
	if !reflect.DeepEqual(uids,uidsWant) {
		t.Error("UID mismatch\nwant:",uidsWant,"\ngot:",uids)
	}
	if !reflect.DeepEqual(errs,errsWant) {
		t.Error("Error mismatch\nwant:",errsWant,"\ngot:",errs)
	}
}

func TestOpen2(t *testing.T) {
	var uids = []uint32{}
	var errs = []error{}
	for i:=0; i<5;i++ {
		uid, err :=Open(m2)
		uids = append(uids, uid)
		errs = append(errs, err)
	}
	var m2Want = map[uint32]int{
		0 : 0,
		1 : 2,
		2 : 9019234,
		3 : 3832,
		4 : 0,
		5 : 0,
		6 : 0,
		7 : 254,
		8 : 0,
	}
	var uidsWant = []uint32{
		0,4,5,6,8,
	}
	var errsWant = []error{
		nil,nil,nil,nil,nil,
	}
	if !reflect.DeepEqual(m2,m2Want) {
		t.Error("Database mismatch\nwant:",m2Want,"\ngot:",m2)
	}
	if !reflect.DeepEqual(uids,uidsWant) {
		t.Error("UID mismatch\nwant:",uidsWant,"\ngot:",uids)
	}
	if !reflect.DeepEqual(errs,errsWant) {
		t.Error("Error mismatch\nwant:",errsWant,"\ngot:",errs)
	}
}
func TestMod(t *testing.T) {
	var bals = []int{}	
	bal, _ := Mod(m3, 0, -2000)
	bals = append(bals, bal)
	bal, _ = Mod(m3, 3, 111)
	bal, _ = Mod(m3, 3, 2903)
	bals = append(bals, bal)
	bal, _ = Mod(m3, 19123, 0)
	bals = append(bals, bal)
	
	var m3Want = map[uint32]int {
		0 : -1029,
		2 : 1678,
		3 : 14129,
		19123 : 48052985,
	}
	var balsWant = []int{
		-1029,14129,48052985,
	}
	if !reflect.DeepEqual(m3,m3Want) {
		t.Error("Database mismatch\nwant:",m3Want,"\ngot:",m3)
	}
	if !reflect.DeepEqual(bals,balsWant) {
		t.Error("Bal mismatch\nwant:",balsWant,"\ngot:",bals)
	}
}

