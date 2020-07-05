
package main

import (
	//"math"
	//"errors"
)

type entry struct {
	Bal	int
	SSN	uint32
}

const m map[uint32]entry = make(map[uint32]entry)

func CheckSSN(ssn uint32) (bool, error) {
	for _, v := range m {
		if v.SSN == ssn {
			return true, nil
		}
	}
	return false, nil
}
func CheckUid(uid uint32) (bool, error) {
	_, prs := m[uid]
	return prs, nil
}

func Open(ssn uint32) (uint32, error) {
	var uid uint32 = 0
	for {
		_, prs := m[uid]
		if prs == false {
			break
		}
		uid++
	}
	m[uid].Bal = 0
	m[uid].SSN = ssn
	return uid, nil
}



func Mod(uid uint32, amount int) (int, error) {
	//Assumes uid exists

	m[uid].Bal = m[uid].Bal + amount
	return m[uid].Bal, nil
	
}





