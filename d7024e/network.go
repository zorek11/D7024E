package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"
	"crypto/sha1"
	"io"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	me       Contact
	target   *Contact
	response *Contact
	kademlia *Kademlia
}

func NewNetwork(me Contact, kad *Kademlia) *Network {
	net := &Network{}
	net.me = me
	net.kademlia = kad
	return net
}

func (network *Network) AddMessage(c *Contact) {
	network.target = c
}

func (network *Network) AddResponse(c *Contact) {
	network.response = c
}

func Listen(me Contact) {
	Addr, err1 := net.ResolveUDPAddr("udp", me.Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	channel := make(chan []byte)
	buf := make([]byte, 1024)
	for {
		go handleMessage(channel, me)
		_, _, err := Conn.ReadFromUDP(buf)
		//fmt.Print("Connection recived: ", UDPaddr)
		channel <- buf
		if err != nil {
			fmt.Println("Read Error: ", err)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	message := buildMessage("ping", network.me.ID.String(), network.me.Address)
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
	message := &protobuf.KademliaMessage{
		Label:         proto.String("LookupContact"),
		Senderid:      proto.String(network.me.ID.String()),
		SenderAddr:    proto.String(network.me.Address),
		LookupContact: proto.String(network.target.String()),
	}
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

func (network *Network) SendFindDataMessage(hash string) {
		//TODO

}

func (network *Network) SendStoreMessage(data string) {
	hash := sha1.new()
	io.WriteString(h, data)
	key := h.Sum(nil)
	//TODO: implement below if byte array and not string. For final solution.
	//	data := []byte("This page intentionally left blank.")
	//fmt.Printf("% x", sha1.Sum(data))

	message := &protobuf.KademliaMessage{
		Label:         proto.String("StoreData"),
		Senderid:      proto.String(network.me.ID.String()),
		SenderAddr:    proto.String(network.me.Address),
		Key: proto.String(key),
		Value: prot.String(data)
	}
	contacts := kademlia.rt.FindClosestContacts(network.me.ID, count)

	for j, contact := range contacts {
		data, err := proto.Marshal(message)
		if err != nil {
			fmt.Println("Marshal Error: ", err)
		}
		Conn, err := net.Dial("udp", contacts[j].Address)
		if err != nil {
			fmt.Println("UDP-Error: ", err)
		}
		defer Conn.Close()

		_, err = Conn.Write(data)
		if err != nil {
			fmt.Println("Write Error: ", err)
		}
	}
}
