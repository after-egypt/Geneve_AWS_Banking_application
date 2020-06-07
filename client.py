import socket 
import sys
import math

clientMsg = str.encode(input("msg:"))

Msg = bytearray()
Msg += bytes.fromhex("0001")
Msg += bytes.fromhex("03")
OptLen = math.ceil(len(clientMsg)/4)
#print(OptLen,len(clientMsg))
if OptLen > 2**5:
    print("Err: Message too long")
    sys.exit()
Msg += OptLen.to_bytes(1, byteorder = 'big')
Msg += clientMsg
Msg += (0).to_bytes(OptLen*4-len(clientMsg), byteorder = 'big')

print(Msg)
serverAddrPort = ("127.0.0.1",10001)
bufferSize = 1024

UDPClientSocket = socket.socket(family=socket.AF_INET,type=socket.SOCK_DGRAM)
UDPClientSocket.sendto(Msg,serverAddrPort)

serverMsg = UDPClientSocket.recvfrom(bufferSize)
print("Server says",serverMsg[0])
