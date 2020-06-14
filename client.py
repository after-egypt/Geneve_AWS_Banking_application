import socket 
import sys
import math

clientMsg = str.encode(input("msg:"))
ttype = int(input("Choose Type: {0-3}"))

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
    print(type(tlv[0]))
    optClass = int.from_bytes(tlv[0], byteorder='big')<<8 + int.from_bytes(tlv[1], byteorder='big')
    if optClass != 0x1234:
        print("Err: Incorrect Option Class")
        sys.exit()
    ttype = int.from_bytes(tlv[2],byteorder='big')
    length = int.from_bytes(tlv[3],byteorder='big')
    bal = int.from_bytes(tlv[4:],"little")
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
print("Server says Class:",optClass,"Type:",ttype,"Number:",bal,"Raw:",serverMsg.hex())

