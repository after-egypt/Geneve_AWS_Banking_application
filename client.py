import socket 


clientMsg = input("msg:")
bytesToSend = str.encode(clientMsg)
serverAddrPort = ("127.0.0.1",10001)
bufferSize = 1024

UDPClientSocket = socket.socket(family=socket.AF_INET,type=socket.SOCK_DGRAM)
UDPClientSocket.sendto(bytesToSend,serverAddrPort)

serverMsg = UDPClientSocket.recvfrom(bufferSize)
print("Server says {}:".format(serverMsg[0]))
