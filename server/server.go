package main

import (
	"encoding/hex"
	//"errors"
	"fmt"
	//"sync"
	//"math"
	"net"
	"os"
	//"github.com/google/gopacket/layers"
	//"math/bits"
)

func CheckError(err error, loop bool) {
	if err != nil {
		fmt.Println("Error: ", err)
		if loop == false {
			os.Exit(0)
		}
	}
}
const OpenAccType uint8 = 1
const queueSize int = 20

var ServerSendAddr *net.UDPAddr

type Message struct {
	Addr *net.UDPAddr
	Data []byte
}
func clientModBal(m BankType, cid ClientID, recv Message, length uint8) error {
	fmt.Println("clientModBal activated")
	value, err := DecodeModBal(recv.Data[12:])
	CheckError(err, true)

	err = Mod(m, cid.UID, value)
	CheckError(err, true)
	fmt.Println("New Bank Entry is:",*m[cid.UID])

	msg,_ := EncodeModBal(cid, 0, m[cid.UID].Bal)
	Send(ServerSendAddr, recv.Addr, msg)

	return nil
}

func clientOpenAccount(m BankType, recv Message, length uint8) error {
	fmt.Println("clientOpenAccount Activated")
	ssn, seqNum, _ := DecodeOpenAcc(recv.Data[4:])
	//needs log checking
	prs, err := CheckSSN(m, ssn)
	CheckError(err, true)

	if !prs {
		uid, err := Open(m, ssn)
		CheckError(err, true)

		msg, _ := EncodeNewAcc(ssn, seqNum, uid)
		fmt.Println("message encoded")
		Send(ServerSendAddr, recv.Addr, msg)
	}
	return nil
}

func idVerify(m BankType, recv Message) (ClientID, bool) {
	fmt.Println("idVerify Activated")
	cid := DecodeClientID(recv.Data[4:12])
	//needs log checking
	fmt.Println("CID:",cid)
	prs, _ := CheckUid(m, cid.UID)
	fmt.Println("Account Exists?",prs)
	return cid, prs //true means account exists
}

func listener(queueSize int, queue chan Message) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err, false)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err, false)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	fmt.Println("listener ready")
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
	queue := make(chan Message, queueSize)
	go listener(queueSize, queue)
	
	var err error
	ServerSendAddr, err = net.ResolveUDPAddr("udp", ":10002")
	CheckError(err, false)

	var wOptHead OptHead
	m := InitBank()
	fmt.Println("main loop ready")
	for {
		recv := <-queue // := should automatically wait if queue is empty
		fmt.Println("\nrecieved:\n"+hex.EncodeToString(recv.Data),"\nfrom:",recv.Addr)
		wOptHead = DecodeOptHead(recv.Data) 

		if wOptHead.Type == OpenAccType {//const
			clientOpenAccount(m, recv, wOptHead.Length)
		} else {
			cid, prs := idVerify(m, recv)	
			if prs {
				switch wOptHead.Type {
				case 0: //Client Change Bal
					err := clientModBal(m, cid, recv, wOptHead.Length)
					CheckError(err, true)
				}
			} else {
				fmt.Println("ID not found:",cid,prs)
				PrintMap(m)
				//send server error
			}
		}

		fmt.Println("done")
	}
}
