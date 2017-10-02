package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	me       Contact
	target   *Contact
	response []Contact
	kademlia *Kademlia
}

func NewNetwork(me Contact, kad *Kademlia) *Network {
	network := &Network{}
	network.me = me
	network.kademlia = kad
	return network
}

func (network *Network) AddMessage(c *Contact) {
	network.target = c
}

func (network *Network) AddResponse(c []Contact) {
	network.response = c
}

func Listen(me Contact, network *Network) {
	messagehandler := NewMessageHandler(network)

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
		go messagehandler.handleMessage(channel, me, *network)
		_, _, err := Conn.ReadFromUDP(buf)
		//fmt.Print("Connection recived: ", UDPaddr)
		channel <- buf

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
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
		Label:      proto.String("LookupContact"),
		Senderid:   proto.String(network.me.ID.String()),
		SenderAddr: proto.String(network.me.Address),
		Lookupcontact: &protobuf.KademliaMessage_LookupContact{
			ID:       proto.String(contact.ID.String()),
			Address:  proto.String(contact.Address),
			Distance: proto.String(contact.Distance.String()),
		},
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
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
