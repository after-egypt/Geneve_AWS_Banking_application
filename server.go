package main
import("fmt"
       "net"
       "os"
       "encoding/hex"
       "errors"
       "math"
       //"encoding/binary"
       "math/bits"
)

func CheckError(err error,loop bool) {
	if err != nil {
		fmt.Println("Error: ",err)
		if loop==false{
			os.Exit(0)
		}
	}
}
func GenevePkg(errCode uint8, msg uint)([]byte, error){
	OptClass := 0x1234
	ByteLen := uint(math.Ceil(float64(bits.Len(msg))/8))
	Length := uint(math.Ceil(float64(ByteLen)/4))
	out := make([]byte, Length*4 + 4)

	if(Length>>5 != 0){
		return out, errors.New("GenevePkg: Length too long")
	}

	out[0] = byte(OptClass>>8)
	out[1] = byte(OptClass & 0x00FF)
	out[2] = byte(errCode)
	out[3] = byte(Length)
	//fmt.Println(fmt.Sprintf("%x",msg))
	for i:= 0; bits.Len(msg) != 0; i++ {
		out[4+i] = byte(msg & 0x0000FF)
		//fmt.Println(fmt.Sprintf("%x",out[4+i]))
		msg = msg>>8
	}
	return out, nil
}
func GeneveUnPkg(msg []byte)(uint8, string, error){
	OptClass := 0xABCD
	var ttype uint8
	var text string
	if(int(msg[0])<<8 + int(msg[1]) != OptClass){
		return 0xFF,"",errors.New("GeneveUnPkg: OptClass incorrect")
	} else{
		ttype = uint8(msg[2])

		if(msg[3]>>5 != 0){
			return 0xFF,"",errors.New("GeneveUnPkg: Length too long")
		} else{

			i := 0
			for msg[4+i] == 0x00 {
				i++}
			text = string(msg[4+i:])
		}
	}
	return ttype, text, nil
}



func main() {

	ServerAddr,err := net.ResolveUDPAddr("udp",":10001")
	CheckError(err,false)

	ServerAddr2, err := net.ResolveUDPAddr("udp",":10002")
	CheckError(err,false)

	ServerConn,err := net.ListenUDP("udp",ServerAddr)
	CheckError(err,false)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	var errCode uint8
	var bal uint
	for {
		n,addr,err := ServerConn.ReadFromUDP(buf)
		CheckError(err,true)

		ttype, text, err := GeneveUnPkg(buf[0:n])
		CheckError(err,true)

		raw := hex.EncodeToString(buf[0:n])
		fmt.Println("Received:",text,"Type:",ttype,"From:",addr,"Raw:",raw)


		ClientConn, err := net.DialUDP("udp",ServerAddr2,addr)
		CheckError(err,true)

		errCode = 0x36
		bal = 0xFCE1265EE
		buf2, err := GenevePkg(errCode,bal)
		CheckError(err,true)
		fmt.Println("Sending Err Code:",errCode,"Number:",bal,"Raw:",hex.EncodeToString(buf2))
		//buf2 := []byte(string(buf[0:n]))
		_,err = ClientConn.Write(buf2)
		CheckError(err,true)

		ClientConn.Close()

	}
}
