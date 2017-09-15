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
	ServerAddr,err := net.ResolveUDPAddr("udp",ip+":"+portString)
	ServerConnection, err := net.ListenUDP("udp", ServerAddr)

	//read connection
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
			n,addr,err := ServerConn.ReadFromUDP(buf)
			fmt.Println("Received ",string(buf[0:n]), " from ",addr)
			if err != nil {
					fmt.Println("Error: ",err)
			}
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
	ServerAddr,err := net.ResolveUDPAddr("udp",ip+":"+portString)
	LocalAddr, err := net.ResolveUDPAddr("udp", ip+":"+portString)
	Connection, err := net.DialUDP("udp", LocalAddr, ServerAddr)

	defer Conn.Close()
	_,err := Conn.Write("ping")
	if err != nil { //check if ping failed
			fmt.Println(err)
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
