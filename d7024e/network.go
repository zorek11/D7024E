package d7024e

import (
	"fmt"
	"net"
	"strconv"
	//"D7024E-Kademlia/protobuf"
	"github.com/golang/protobuf/proto"
)

type Network struct {
	me Contact
}

func Listen(ip string, port int) {
	Addr, err1 := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	Conn, err2 := net.ListenUDP("udp", Addr)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error: ", err1, "\n", err2)
	}

	//read connection
	defer Conn.Close()
	buf := make([]byte, 1024)
	for {
		_, _, err := Conn.ReadFromUDP(buf)
		go handleMessage(buf, ip+":"+strconv.Itoa(port))
		if err != nil {
			fmt.Println("Read Error: ", err)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	message := buildMessage("ping", contact.ID.String(), contact.Address)
	data, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("Marshal Error: ", err)
	}
	Conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	defer Conn.Close()

	_, err = Conn.Write(data)
	if err != nil {
		fmt.Println("Write Error: ", err)
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
