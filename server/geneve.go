package main

import(
	"encoding/binary"
	"errors"
	"net"
	"math/bits"
	"math"
	"fmt"
	"encoding/hex"
)

const ServerClass uint16 = 0x0000

func checkErr(err error) {
	if err != nil {
		fmt.Println("Error:",err)
	}
}
		
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
func Send(serverAddr *net.UDPAddr, clientAddr *net.UDPAddr, msg []byte) {
	clientConn, err := net.DialUDP("udp", serverAddr, clientAddr)
	checkErr(err)
	_, err = clientConn.Write(msg)
	checkErr(err)
	fmt.Println("sent:",hex.EncodeToString(msg))
	clientConn.Close()
}


func DecodeOptHead(data []byte) (OptHead) { 
	var opt OptHead

	opt.Class = binary.BigEndian.Uint16(data[0:2])
	opt.Type = data[2]
	opt.Flags = data[3] >> 5
	opt.Length = (data[3] & 0x1f)
	
	return opt
}
func EncodeOptHead(opt OptHead) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[:2], opt.Class)
	b[2] = byte(opt.Type)
	b[3] = byte((opt.Length & 0x1f) + opt.Flags << 5)

	return b
}
func DecodeClientID(data []byte) (ClientID) { //data must not include option header
	_ = data[7] //bounds check
	var cid ClientID
	cid.UID = binary.BigEndian.Uint32(data[:4])
	cid.SeqNum = binary.BigEndian.Uint32(data[4:8])
	return cid
}
func EncodeClientID(cid ClientID) []byte {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint32(msg[:4], cid.UID)
	binary.BigEndian.PutUint32(msg[4:8], cid.SeqNum)
	return msg
}
/* 
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Option Class(vendor specific) |   2           |R|R|R| Length  |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                       Social Security Num                     |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                 Sequence Number                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |            User ID                                            |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 */
func EncodeNewAcc (ssn uint32, seqNum uint32, uid uint32) ([]byte, error) { //Encode is not the opposite of Decode!
	msg := make([]byte, 16)
	copy(msg[:4], EncodeOptHead(OptHead{ServerClass, 0x02, 0x0, 0x3}))
	
	binary.BigEndian.PutUint32(msg[4:8], ssn)
	binary.BigEndian.PutUint32(msg[8:12], seqNum)
	binary.BigEndian.PutUint32(msg[12:16], uid)

	return msg, nil
}
/* 
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| Option Class(vendor specific) |   3           |R|R|R| Length  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       User ID                                 |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                 Sequence Number                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|    error code                 |S|     balance                 |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                             +
|                                                               | 
~                                                               ~
|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                               |         reserved              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
func EncodeModBal (cid ClientID, errCode uint16, bal int) ([]byte, error) {//Type specific Encodes currently handle everything including options and clientID, this may change
	balLen := math.Ceil(float64( bits.Len(uint(math.Abs(float64(bal))))+1 )/16)*2 //Number of Bytes, rounded to nearest 16
	length := uint8(math.Ceil((10 + balLen)/4))
	msg := make([]byte, length*4 + 4)
	copy(msg[:4], EncodeOptHead(OptHead{ServerClass, 0x03, 0x0, length}))

	copy(msg[4:12], EncodeClientID(cid))

	binary.BigEndian.PutUint16(msg[12:14], errCode)
	switch balLen {
	case 2:
		binary.BigEndian.PutUint16(msg[14:16], uint16(math.Abs(float64(bal))))
	case 4:
		binary.BigEndian.PutUint32(msg[14:18], uint32(math.Abs(float64(bal))))
	case 6:
		binary.BigEndian.PutUint64(msg[14:22], uint64(math.Abs(float64(bal))))
	case 8:
		binary.BigEndian.PutUint64(msg[14:22], uint64(math.Abs(float64(bal))))
	default:
		return nil, errors.New("Balance byteLength calculation messed up")
	}
	if bal < 0 {
		msg[14] = msg[14] | 0x80 //1000 0000
	}
	
	return msg, nil
}
/*
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Option Class         |      1        |R|R|R| Length  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      Social Security Num                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|             Sequence Number                                   |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+       
|                      Name                                     |
~                                                               ~
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+*/
func DecodeOpenAcc (data []byte) (uint32, uint32, string) {//type specific decodes know nothing except the bytes they need. no options data.
	ssn := binary.BigEndian.Uint32(data[:4])
	seqNum := binary.BigEndian.Uint32(data[4:8])
	name := string(data[8:])
	return ssn, seqNum, name
}
/*
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Option Class         |      0        |R|R|R| Length  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       User ID                                 |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|             Sequence Number                                   |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+       
|S|                      Value                                  |
~-+                                                             ~
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+*/
func DecodeModBal (data []byte) (int, error) {//Function doesn't know cid
	var value int = 0
	length := uint(len(data) /4)
	if (length > 0){
		sign := uint(data[0] & 0x80) //1000 0000
		data[0] = data[0] & 0x7F //0111 1111
		switch length {
		case 1:
			value = int(binary.BigEndian.Uint32(data[0:4]))
		case 2:
			value = int(binary.BigEndian.Uint64(data[0:8]))
		default:
			return 0, errors.New("Length too long for uint --DECODE MOD BAL")
		}
		if (sign != 0) {
			value = -1* value
		}
	}
	return value, nil
}
func EncodeDelAcc (cid ClientID, errCode uint16) ([]byte) {
	msg := make([]byte, 16)
	copy(msg[:4], EncodeOptHead(OptHead{ServerClass, 0x06, 0x0, 0x03}))
	copy(msg[4:12], EncodeClientID(cid))
	binary.BigEndian.PutUint16(msg[12:14],errCode)
	return msg
}
func EncodeFatalErr(cid ClientID, errCode uint16) ([]byte){
	msg := make([]byte, 16)
	copy(msg[:4], EncodeOptHead(OptHead{ServerClass, 0x04, 0x0, 0x03}))
	copy(msg[4:12], EncodeClientID(cid))
	binary.BigEndian.PutUint16(msg[12:14], errCode)
	return msg
}
