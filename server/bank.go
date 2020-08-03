
package main

import (
	//"math"
	//"errors"
	"fmt"
)


type BankEntry struct {
	Bal	int
	SSN	uint32
}

type BankType map[uint32]*BankEntry //dereference necessary to allow m[uid].Bal += x

func PrintMap(m BankType) { //debugging purpose
	for k, v := range m {
		fmt.Println(k," : ",*v)
	}
}
func InitBank() (BankType) {
	var m BankType = make(BankType)
	return m
}

func CheckSSN(m BankType, ssn uint32) (bool, error) { //checks if a SSN is registered
	for _, v := range m {
		if v.SSN == ssn {
			return true, nil
		}
	}
	return false, nil
}
func CheckUid(m BankType, uid uint32) (bool, error) { //true = uid exists
	_, prs := m[uid]
	return prs, nil
}

func Open(m BankType, ssn uint32) (uint32, error) {//BankType is passed as reference
	var uid uint32 = 0
	for {
		_, prs := m[uid]
		if prs == false {
			break
		}
		uid++
	}
	m[uid] = &BankEntry{Bal: 0,SSN: ssn}
	return uid, nil
}

func Mod(m BankType, uid uint32, amount int) (error) {
	//Assumes uid exists

	m[uid].Bal = m[uid].Bal + amount
	return nil
	
}
func Delete(m BankType, uid uint32) (uint16) {
	delete(m, uid)
	return 0
}





