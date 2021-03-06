package d7024e

import (
    "fmt"
    "net"
    "time"
    "strconv"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
	portString := strconv.Itoa(port) //convert int to string
	ServerAddr, err := net.ResolveUDPAddr("udp",ip+":"+portString)
	if err  != nil {
			fmt.Println("Error: " , err)
	}
	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	//read connection
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	//for {
			n,addr,err := ServerConn.ReadFromUDP(buf)
			fmt.Println("Received ",string(buf[0:n]), " from ",addr)
			if err != nil {
					fmt.Println("Error: ",err)
			//}
	}

}

func SendPingMessage(ip string, port int) {
	// TODO
	portString := strconv.Itoa(port) //convert int to string
	ServerAddr,err := net.ResolveUDPAddr("udp",ip+":"+portString)
	if err  != nil {
			fmt.Println("Error: " , err)
	}
	thisIp := "127.0.0.1"
	thisPort := "8081"
	LocalAddr, err := net.ResolveUDPAddr("udp", thisIp+":"+thisPort)
	if err  != nil {
			fmt.Println("Error: " , err)
	}
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err  != nil {
			fmt.Println("Error: " , err)
	}

	defer Conn.Close()
	buf := []byte("PING")
	_,err = Conn.Write(buf)
	if err  != nil {
			fmt.Println("Error: " , err)
	}

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
