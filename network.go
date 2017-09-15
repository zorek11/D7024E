package d7024e

import (
    "fmt"
    "net"
    "strconv"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
  fmt.Println("ListnerUp")
	ServerAddr,err1 := net.ResolveUDPAddr("udp",ip+":"+strconv.Itoa(port))
	ServerConn, err2 := net.ListenUDP("udp", ServerAddr)
  if (err1  != nil) || (err2 != nil){
      fmt.Println("Error: ", err1,"\n", err2)
  }

	//read connection
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
			n,addr,err := ServerConn.ReadFromUDP(buf)
			fmt.Println("Received ",string(buf[0:n]), " from ",addr)
			if err != nil {
					fmt.Println("Error: ",err)
			}
      break;
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
	ServerAddr,err := net.ResolveUDPAddr("udp", contact.Address)
	LocalAddr, err2 := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	Conn, err3 := net.DialUDP("udp", LocalAddr, ServerAddr)
  if (err  != nil) || (err2  != nil) || (err3  != nil) {
      fmt.Println("Error: ", err, "\n", err2, "\n", err3)
  }

  buf := []byte("PING")
  _,err = Conn.Write(buf)
  if err  != nil {
      fmt.Println("Error: ", err)
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
