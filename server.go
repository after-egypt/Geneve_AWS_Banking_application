package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	//"sync"
	//"math"
	"encoding/binary"
	"net"
	"os"
	//"github.com/google/gopacket/layers"
	//"math/bits"
	"notgo/vbinary"
)

func CheckError(err error, loop bool) {
	if err != nil {
		fmt.Println("Error: ", err)
		if loop == false {
			os.Exit(0)
		}
	}
}
const OpenAccType int = 1

type Message struct {
	Addr *net.UDPAddr
	Data []byte
}
func Calculate(x int) (result int) {
	result = x + 2
	return result
}
func clientModBal(cid ClientID, recv Message, length uint8, database *[]int) ([]byte, error) {
	//var opt optType0
	var resp []byte
	switch length {
	case 1:
		value := vbinary.BigEndian.Int32(recv.Data)
	case 2:
		value := vbinary.BigEndian.Int64(recv.Data)
	default:
		resp := make([]byte, 4)
		binary.BigEndian.PutUint16(resp, 1<<15+0x0001)
		return resp, errors.New("value too big")
	}

	*database[cid.UID] += value

}

func clientOpenAccount(recv Message, serverAddr2 *net.UDPAddr) error {
	ssn := binary.BigEndian.Uint32(recv.Data[:4])
	seqNum := binary.BigEndian.Uint32(recv.Data[4:8])
	//needs log checking
	prs, err := CheckSSN(ssn)
	CheckError(err, true)

	if !prs {
		uid, err := Open(ssn)
		CheckError(err, true)

		clientConn, err := net.DialUDP("usp", serverAddr2, recv.Addr)
		Send2(serverAddr2, recv.Addr, ssn, seqNum, uid)
		clientConn.Close()
	}
}

func idVerify(recv Message) (ClientID, []byte, error) {
	var cid ClientID
	cid.UID = binary.BigEndian.Uint32(recv.Data[:4])
	cid.SeqNum = binary.BigEndian.Uint32(recv.Data[4:8])
	if CheckUid(cid.UID) {
		return cid, recv.Data[8:], nil
	} else {
		return cid, recv.Data[8:], errors.New("UID found")
	}
}

func listener(queueSize int, queue chan Message, wg sync.WaitGroup) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err, false)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err, false)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		CheckError(err, true)

		if len(queue) >= queueSize {
			fmt.Println("Error listener: queue is full")
		} else {
			queue <- Message{addr, buf[0:n]}
			
		}
	}
}

func main() {
	queueSize := 20
	queue := make(chan Message, queueSize)
	go listener(queueSize, queue, wg)

	ServerAddr2, err := net.ResolveUDPAddr("udp", ":10002")
	CheckError(err, false)

	buf := make([]byte, 1024)
	var ttype uint8
	var bal uint
	var recv Message
	var wOptHead OptHead
	var cid ClientID

	var database []int

	for {
		recv := <-queue

		wOptHead, recv.Data = decodeOptHead(recv.Data) //decodeOptHead removes the header from the data

		if wOptHead.Type == OpenAccType {//const
			clientOpenAccount(recv)
		} else {
			cid, recv.Data, err := idVerify(recv)	
			checkErr(err,true)
			switch wOptHead.Type {
			case 0:
				Mod(cid.UID, 

		}

		if err == nil {
			switch wOptHead.Type {
			case 0:
				clientModBal(cid, recv, wOptHead.Length, &database)
			case 1:
				clientOpenAccount(recv, wOptHead.Length, &database)
			}
		}
		CheckError(err, true)

		fmt.Println("Recieved:", hex.EncodeToString(recv.Data))
		/*ClientConn, err := net.DialUDP("udp", ServerAddr2, addr)
		CheckError(err, true)

		buf2, err := encodeOptHead(uint8(ttype), bal)
		CheckError(err, true)
		fmt.Println("Sending ttype:", ttype, "Number:", bal, "Raw:", hex.EncodeToString(buf2))
		//buf2 := []byte(string(buf[0:n]))
		_, err = ClientConn.Write(buf2)
		CheckError(err, true)

		ClientConn.Close()*/

	}
}
