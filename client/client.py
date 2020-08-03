import socket 
import sys
import math
import binascii

optClass = bytes.fromhex("0000")
serverAddrPort = ("192.168.1.13",10001)
bufferSize = 1024

def utf8len(s):
    return int(len(s.encode('utf-8')))

def EncodeOpt(tlvType, flags, length):
    msg = bytearray()
    msg += optClass
    msg += tlvType.to_bytes(1, byteorder='big')
    msg += ((flags<<5) + length).to_bytes(1, byteorder='big')
    print("EncodeOpt Result:",binascii.hexlify(msg),"length is supposed to be:",((flags<<5) + length))
    return msg

def GetSeqNum():
        return bytes.fromhex("ABCDEF12")



def craftOpenAcc():
    print("craftOpenAcc starts")
    ssn = int(input("Social Security Number: "))
    name = str(input("Name: "))
    length = int(2+(utf8len(name)/4))
    
    msg = bytearray()
    msg += EncodeOpt(1,0,length)
    msg += ssn.to_bytes(4, byteorder='big')
    msg += GetSeqNum()
    msg += name.encode('utf-8')
    return msg

def craftModAcc(negative): #true = negative
    print("craftModBal starts")
    uid = int(input("User ID: "))
    if (negative):
        value = int(input("Amount to Withdraw:"))
    else:
        value = int(input("Amount to Deposit:"))

    if (value < 0):
        negative = not negative
        value = int(abs(value))

    length = math.ceil((value.bit_length() + 1)/32) + 2
    print("Using TLV Length:",length)

    msg = bytearray()
    msg += EncodeOpt(0, 0, length)
    msg += uid.to_bytes(4, byteorder='big')
    msg += GetSeqNum()
    msg += value.to_bytes((length-2)*4, byteorder='big')
    if (negative):
        msg[12] = msg[12] | 0x80 # 1000 0000
    return msg
def craftDelAcc():
    uid = int(input("User ID: "))

    msg = bytearray()
    msg += EncodeOpt(5, 0, 2)
    msg += uid.to_bytes(4, byteorder='big')
    msg += GetSeqNum()
    return msg


while True:
    option = int(input("Which type of transaction?\n0: Open Account\n1: Deposit\n2: Withdraw\n3: Delete Account\n(type a number)"))
    print(type(option))

    def inputSwitch(argument):
        def craftModAccTrue():
            return craftModAcc(True)
        def craftModAccFalse():
            return craftModAcc(False)
        def returnEmpty():
            return b''
        print("inputSwitch starts")
        switcher= {
            0: craftOpenAcc,
            1: craftModAccFalse,
            2: craftModAccTrue,
            3: craftDelAcc
            }
        return switcher.get(argument, returnEmpty)

    msg = bytearray()
    msg += inputSwitch(option)()

    print("end of loop result",binascii.hexlify(msg),"\n",type(msg))
    if (msg != b''):
        break
    print("Invalid Option!")

UDPClientSocket = socket.socket(family=socket.AF_INET,type=socket.SOCK_DGRAM)
UDPClientSocket.sendto(msg,serverAddrPort)
print("Request Sent")
while True:
    servermsg = UDPClientSocket.recvfrom(bufferSize)
    servermsg = bytearray(servermsg[0])
    print("Server says",binascii.hexlify(servermsg))
    if (servermsg[4:12] == msg[4:12]):
        break

def serverOpenAcc(msg):
    uid = int.from_bytes(servermsg[12:], "big")
    print("You user ID is:",uid)
    return

def serverChangeBal(msg):
    sign = 1
    if ((servermsg[14] & 0x80) != 0):
        servermsg[14] = servermsg[14] & 0x7F
        sign = -1
    TLVLen = int(msg[3] & 0x1F)
    numLen = (2**(TLVLen -2)) #num of bytes for balance
    bal = int.from_bytes(servermsg[14:14+numLen], "big")
    bal = sign*bal
    print("Your balance is:",bal)
    return

def serverFatalError(msg):
    errNum = int.from_bytes(servermsg[12:14], "big")
    print("Fatal Error",errNum)
    return

def serverDelAcc(msg):
    errCode = int.from_bytes(msg[12:14], "big")
    if errCode != 0:
        print("Warning:",0)
    print("Account deletion successful")

def responseSwitch(argument):
    switcher= {
            2:  serverOpenAcc,
            3:  serverChangeBal,
            4:  serverFatalError,
            6:  serverDelAcc,
            }
    return switcher.get(argument)

TLVtype = servermsg[2]
responseSwitch(TLVtype)(msg)

'''def GenevePkg(ttype,clientmsg):
    msg = bytearray()
    Length = math.ceil(len(clientmsg)/4)
    #print(Length,len(clientmsg))
    if Length>>5 != 0:
        print("Err: Message too long")
        sys.exit()
    msg += bytes.fromhex("ABCD")
    msg += ttype.to_bytes(1, byteorder='big')
    msg += Length.to_bytes(1, byteorder='big')
    msg += (0).to_bytes(Length*4-len(clientmsg), byteorder='big')
    msg += clientmsg
    return msg

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
msg = GenevePkg(ttype,clientmsg)
print(msg)
serverAddrPort = ("127.0.0.1",10001)
bufferSize = 1024

UDPClientSocket = socket.socket(family=socket.AF_INET,type=socket.SOCK_DGRAM)
UDPClientSocket.sendto(msg,serverAddrPort)

servermsg = UDPClientSocket.recvfrom(bufferSize)
ttype, length, bal = GeneveUnPkg(servermsg[0])
#print("Server says",servermsg[0][0],"Type:",type(servermsg[0]))
print("Server says Type:",ttype,"Number:",bal,"Raw:",servermsg[0])
'''
