package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	me       Contact
	target   *Contact
	response [][]Contact
	temp     *Contact
	rt       *RoutingTable
	storage  Storage
}

func NewNetwork(me Contact, rt *RoutingTable, st Storage) Network {
	network := Network{}
	network.me = me
	network.rt = rt
	network.storage = st
	return network
}

func (network *Network) AddMessage(c *Contact) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	network.target = c
}

func (network *Network) GetTemp() *Contact {
	return network.temp
}

/*func (network *Network) GetKademlia() *Kademlia {
	return network.kademlia
}*/

func (network *Network) AddResponse(c []Contact) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	network.response = append(network.response, [][]Contact{c}...)
	fmt.Println("\nResponse: ", network.response)
}

func (network *Network) RemoveFirstResponse() {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	network.response = network.response[1:]
}

func (network *Network) GetResponse() [][]Contact {
	return network.response
}
func (network *Network) AddTempResponse(c *Contact) {
	network.temp = c
}

func (network *Network) Listen(me Contact) {
	messagehandler := NewMessageHandler(network)
	Addr, err1 := net.ResolveUDPAddr("udp", me.Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error Listen: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	channel := make(chan []byte)
	buf := make([]byte, 1024)

	for {

		go messagehandler.handleMessage(channel, me, network)
		n, _, err := Conn.ReadFromUDP(buf)
		channel <- buf[0:n]
		//fmt.Println("Connection recived: ", string(buf[0:n]), " \nfrom ", addr)

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
	//time.Sleep(time.Second * 2)

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	message := &protobuf.KademliaMessage{
		Label:      proto.String("LookupContact"),
		Senderid:   proto.String(network.me.ID.String()),
		SenderAddr: proto.String(network.me.Address),
		Lookupcontact: &protobuf.KademliaMessage_LookupContact{
			Id:      proto.String(network.target.ID.String()),
			Address: proto.String(network.target.Address),
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

func (network *Network) SendStoreMessage(contact *Contact, key *KademliaID, value string) {

	message := &protobuf.KademliaMessage{
		Label:         proto.String("StoreData"),
		Senderid:      proto.String(network.me.ID.String()),
		SenderAddr:    proto.String(network.me.Address),
		Key: proto.String(key.String()),
		Value: proto.String(value),
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
