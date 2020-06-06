package main
import("fmt"
       "net"
       "os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ",err)
		os.Exit(0)
	}
}

func main() {
	ServerAddr,err := net.ResolveUDPAddr("udp",":10001")
	CheckError(err)

	ServerAddr2, err := net.ResolveUDPAddr("udp",":10002")
	CheckError(err)

	ServerConn,err := net.ListenUDP("udp",ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n,addr,err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received: ",string(buf[0:n])," From: ",addr)
		if err != nil {
			fmt.Println(err)
		}

		ClientConn, err := net.DialUDP("udp",ServerAddr2,addr)
		CheckError(err)

		buf2 := []byte(string(buf[0:n]))
		_,err = ClientConn.Write(buf2)
		if err != nil{
			fmt.Println(err)
		}
		ClientConn.Close()

	}
}
