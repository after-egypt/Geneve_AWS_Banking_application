package main

import(
	"encoding/binary"
	"errors"
	"notgo/vbinary"
	"net"
	"math/bits"
	"math"
)
		
type ClientID struct {
	UID    uint32
	SeqNum uint32
}
type OptHead struct { //TLV header
	Class  uint16
	Type   uint8
	Flags  uint8
	Length uint8
}

func DecodeOptHead(data []byte) (OptHead, []byte) {
	var opt OptHead

	opt.Class = binary.BigEndian.Uint16(data[0:2])
	opt.Type = data[2]
	opt.Flags = data[3] >> 4
	opt.Length = (data[3] & 0x1f)

	var msg = make([]byte, opt.Length*4)
	copy(msg, data[4:])

	return opt, msg
}

func EncodeOptHead(opt OptHead) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[:2], opt.Class)
	b[2] = byte(opt.Type)
	b[3] = byte((opt.Length & 0x1f) + opt.Flags << 5)

	return b
}

func EncodeClientID(cid ClientID) []byte {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint32(msg[:4], cid.UID)
	binary.BigEndian.PutUint32(msg[4:8], cid.SeqNum)
}

func Send2 (ClientConn *net.UDPConn, ssn uint32, seqNum uint32, uid uint32) error {
	msg := make([]byte, 16)
	msg = EncodeOptHead(OptHead{0x0000, 0x02, 0x0, 0x3})
	
	binary.BigEndian.PutUint32(msg[4:8], ssn)
	binary.BigEndian.PutUint32(msg[8:12], seqNum)
	binary.BigEndian.PutUint32(msg[12:16], uid)

	_, err := ClientConn.Write(msg)
	return err
}

func Send3 (ClientConn *net.UDPConn, cid ClientID, errCode int16, bal int) error {
	balLen := math.Ceil(float64( bits.Len(uint(math.Abs(bal))) )/16)*2 //Number of Bytes, rounded to nearest 16
	length := uint8(math.Ceil((10 + balLen)/4))
	msg := make([]byte, length*4 + 4)
	msg := EncodeOptHead(OptHead{0x0000, 0x03, 0x0, length})

	copy(msg[4:12], EncodeClientID(cid))

	vbinary.BigEndian.PutInt16(msg[12:14], errCode)
	switch balLen {
	case 2:
		vbinary.BigEndian.PutInt16(msg[14:16], int16(bal))
	case 4:
		vbinary.BigEndian.PutInt32(msg[14:18], int32(bal))
	case 6 || 8:
		vbinary.BigEndian.PutInt64(msg[14:22], int64(bal))
	default:
		return errors.New("Balance byteLength calculation messed up")
	}

	_, err := ClientConn.Write(msg)
	return err
}

