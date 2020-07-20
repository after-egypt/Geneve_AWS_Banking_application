import socket 
import sys
import math

clientMsg = str.encode(input("msg:"))

def GenevePkg(ttype,clientMsg):
    Msg = bytearray()
    Length = math.ceil(len(clientMsg)/4)
    #print(Length,len(clientMsg))
    if Length>>5 != 0:
        print("Err: Message too long")
        sys.exit()
    Msg += bytes.fromhex("ABCD")
    Msg += ttype.to_bytes(1, byteorder='big')
    Msg += Length.to_bytes(1, byteorder='big')
    Msg += (0).to_bytes(Length*4-len(clientMsg), byteorder='big')
    Msg += clientMsg
    return Msg

def GeneveUnPkg(tlv):
    print(type(tlv[1:1]))
    #python automatically interprets bytearray[idx] as type int
    optClass = tlv[0]<<8 | tlv[1]
    print(hex(optClass), hex(tlv[0]<<8),hex(tlv[1]))
    if optClass != 0x1234:
        print("Err: Incorrect Option Class")
        sys.exit()
    ttype = tlv[2]
    length = tlv[3]
    bal = int.from_bytes(tlv[4:],"little")
    print(hex(bal))
    return ttype, length, bal
Msg = GenevePkg(ttype,clientMsg)
print(Msg)
serverAddrPort = ("127.0.0.1",10001)
bufferSize = 1024

UDPClientSocket = socket.socket(family=socket.AF_INET,type=socket.SOCK_DGRAM)
UDPClientSocket.sendto(Msg,serverAddrPort)

serverMsg = UDPClientSocket.recvfrom(bufferSize)
ttype, length, bal = GeneveUnPkg(serverMsg[0])
#print("Server says",serverMsg[0][0],"Type:",type(serverMsg[0]))
print("Server says Type:",ttype,"Number:",bal,"Raw:",serverMsg[0])

