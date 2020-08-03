## Geneve Banking Application
This is a Bank program using the Geneve protocol to communicate between clients and a server. The client and server can be on different computers.

Geneve is a protocol for sending information like http for websites, but more low level. It encapsulates raw binary with info about which specification should be used to interperet the bits. These specifications are specific to one company. Mine are defined in the file design.txt . 
The purpose of Geneve is in datacenters. TCP (The standard info sharing protocol on the internet) has lots of redundancy because the internet is large and chaotic. This 
#### Features:
- Openingi/Closing Accounts
- Depositing/Withdrawing
- Works over network
- Multiple users
- Uses Geneve for transmitting data

#### Not features:
- Security


#### How to run:	
This program is developed on linux. Theoretically it will work on all operating systems.
The Server's IP and port is hardcoded in client.py near line 7 `serverAddrPort = ("192.168.1.13",10001)`. **This must be changed to the ip of your server computer**. The program can only understant ipv4 addresses.
The Server's Port for sending and recieving are hardcoded in server.go. These
##### Linux:
On the server side:
1. `$ ifconfig`
This is to find the server's ip address. If you're running the client on the same computer, find lo (loopback) and remember the numbers after "inet" (usually 127.0.0.1). If you're running the client on a different computer on the same network, find 
1. `$ cd /the/directory/you/downloaded/this/to/`
1. `$ cd ./server/`
1. `$ go run server.go bank.go geneve.go`
`	This will run the server in the terminal. Info will show up when clients send requests.

On the Client's side:
1. `$ cd ./client/`
1. `$ python3.7 client.py`
1. Follow the instructions printed on screen. Remember your user id because you need that to access your account.

	

####Dependencies:
- Python 3.7
- Go 1.14.2
