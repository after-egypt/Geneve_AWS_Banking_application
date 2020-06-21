package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	//"encoding/binary"
	"github.com/google/gopacket/layers"
	"math/bits"
)

func CheckError(err error, loop bool) {
	if err != nil {
		fmt.Println("Error: ", err)
		if loop == false {
			os.Exit(0)
		}
	}
}

type optHead struct{ //TLV header
	Class   uint16
	Type    uint8
	Flags	uint8
	Length  uint8
}
type optType0 struct{ //Client Withdraw, Deposit, Balance
	UID	uint32 //UserID
	SeqNum	uint32 //Sequence ID/Number
	Value   int
}
type optType1 struct{ //Client Open Acount
	RandID	uint32
	SeqID	uint32
}
type optType2 struct{ //Server Account Opened
	RandID	uint32
	SeqNum	uint32
	UID	uint32
}
type optType3 struct{ //Server Response (change balance)
	UID	uint32
	SeqNum	uint32
	Fatal	bool
	Error	uint16
	Bal	int
}

func decodeOptHead(data []byte) (optHead, []byte) {
	var opt optHead

	opt.Class = binary.BigEndian.Uint16(data[0:2])
	opt.Type = data[2]
	opt.Flags = data[3] >> 4
	opt.Length = (data[3]&0x1f)*4 + 4

	msg = make([]byte, opt.Length-4)
	copy(msg, data[4:opt.Length])

	return opt, msg
}

func decodeOpt0(data []byte) (optType0, error) {
	var opt optType0

	opt.UID = binary.BigEndian.Uint32(data[0:4])
	opt.SeqNum = binary.BigEndian.Uint32(data[4:8])
	opt.Value = binary.BigEndian


func encodeOptHead(opt *optHead) ([]byte, error) {
	msg := make([]byte, 4)

	binary.BigEndian.PutUint16(msg[0:2], opt.Class)
	binary.BigEndian.PutUint8(msg[2], opt.Type)
	if opt.Length >> 5 != 0 {
		return msg, errors.New("encodeOptHead: Length too long")
	}
	binary.BigEndian.PutUint8(msg[3], opt.Flags<<5 + opt.Length)

	return msg, nil
}

		
	
func genevePkg(Type uint8, msg uint) ([]byte, error) {
	OptClass := 0x1234
	ByteLen := uint(math.Ceil(float64(bits.Len(msg)) / 8))
	Length := uint(math.Ceil(float64(ByteLen) / 4))
	out := make([]byte, Length*4+4)

	if Length>>5 != 0 {
		return out, errors.New("GenevePkg: Length too long")
	}

	out[0] = byte(OptClass >> 8)
	out[1] = byte(OptClass & 0x00FF)
	out[2] = byte(Type)
	out[3] = byte(Length)
	//fmt.Println(fmt.Sprintf("%x",msg))
	for i := 0; bits.Len(msg) != 0; i++ {
		out[4+i] = byte(msg & 0x0000FF)
		//fmt.Println(fmt.Sprintf("%x",out[4+i]))
		msg = msg >> 8
	}
	return out, nil
}

type message struct {
	Addr *net.UDPAddr
	Data []byte
}

func listener(queueSize int, queue chan message, wg sync.WaitGroup) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err, false)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err, false)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		CheckError(err,true)

		if len(queue) >= queueSize {
			fmt.Println("Error listener: queue is full")
		} else {
			queue <- message{addr, buf[0:n]}
			wg.Add(1)
			wg.Done()
		}
	}
}

func main() {
	queueSize := 20
	queue := make(chan message, queueSize)
	var wg sync.WaitGroup
	go listener(queueSize, queue, wg)

	ServerAddr2, err := net.ResolveUDPAddr("udp", ":10002")
	CheckError(err, false)

	buf := make([]byte, 1024)
	var ttype uint8
	var bal uint

	for {
		if len(queue) == 0 {
			wg.Wait()
		}
		[]byte buf <- queue

		optHead woptHead, []byte data = decodeOptHead(buf)
		switch woptHead.Type{
		case 1:


		CheckError(err, true)

		raw := hex.EncodeToString(buf)
		fmt.Println("Received:", text, "ttype:", ttype, "From:", addr, "Raw:", raw)
		ClientConn, err := net.DialUDP("udp", ServerAddr2, addr)
		CheckError(err, true)

		ttype = 0x36
		bal = 0xFCE1265EE
		buf2, err := GenevePkg(uint8(ttype), bal)
		CheckError(err, true)
		fmt.Println("Sending ttype:", ttype, "Number:", bal, "Raw:", hex.EncodeToString(buf2))
		//buf2 := []byte(string(buf[0:n]))
		_, err = ClientConn.Write(buf2)
		CheckError(err, true)

		ClientConn.Close()

	}
}
