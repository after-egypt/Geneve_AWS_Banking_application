import socket 
import sys
import math
import binascii
msg = bytes.fromhex(input())

serverAddrPort = ("127.0.0.1",10001)
bufferSize = 1024

UDPClientSocket = socket.socket(family=socket.AF_INET,type=socket.SOCK_DGRAM)
UDPClientSocket.sendto(msg,serverAddrPort)
print("sent")
serverMsg = UDPClientSocket.recvfrom(bufferSize)
print(binascii.hexlify(serverMsg[0]))
