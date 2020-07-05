package vbinary

import (
	//"math"
//"encoding/binary"
) 

type ByteOrder interface {
	Uint([]byte) uint
	Int([]byte) int
	Int8([]byte) int8
	Int16([]byte) int16
	Int32([]byte) int32
	Int64([]byte) int64
	PutInt8([]byte) int8
	PutInt16([]byte) int16
	PutInt32([]byte) int32
	PutInt64([]byte) int64
}

var LittleEndian littleEndian

var BigEndian bigEndian

type littleEndian struct{}

type bigEndian struct{}

func (bigEndian) Int32(b []byte) int32 {
	_ = b[3]
	sign := b[0] >> 7
	if sign == 0 {
		return int32(uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | (uint32(b[0])&0x7f)<<24)
	} else {
		return -1 * int32(uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | (uint32(b[0])&0x7f)<<24)
	}
}

func (bigEndian) Int64(b []byte) int64 {
	_ = b[7]
	sign := b[0] >> 7
	if sign == 0 {
		return int64(uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
			uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | (uint64(b[0])&0x7f)<<56)
	} else {
		return -1 * int64(uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 | 
			uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | (uint64(b[0])&0x7f)<<56)
	}
}

func (bigEndian) PutInt16(b []byte, v int16) {
	_ = b[1]
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	if v & 0x80 == 1 { //v is negative
		b[0] = byte(b[0] + 1<<15)
	}
}

func (bigEndian) PutInt32(b []byte, v int32) {
	_ = b[3]
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
	if v & 0x8000 == 1 { //v is negative
		b[0] = byte(b[0] + 1<<31)
	}
}

func (bigEndian) PutInt64(b []byte, v int64) {
	_ = b[7]
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	if v & 0x80000000 == 1 { //v is negative
		b[0] = byte(b[0] + 1<<63)
	}
}



