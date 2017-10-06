package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"container/list"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	me          Contact
	target      *Contact
	response    [][]Contact
	temp        *Contact
	rt          *RoutingTable
	mtx         sync.Mutex
	pingResp    bool
	pingChannel chan bool
}

func NewNetwork(me Contact, rt *RoutingTable) Network {
	network := Network{}
	network.me = me
	network.rt = rt
	network.pingChannel = make(chan bool)
	return network
}

func (network *Network) AddMessage(c *Contact) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.target = c
}

func (network *Network) GetTemp() *Contact {
	return network.temp
}

/*func (network *Network) GetKademlia() *Kademlia {
	return network.kademlia
}*/

func (network *Network) AddResponse(c []Contact) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.response = append(network.response, [][]Contact{c}...)
	fmt.Println("\nResponse: ", network.response)
}

func (network *Network) RemoveFirstResponse() {

	network.mtx.Lock()
	defer network.mtx.Unlock()
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
		go messagehandler.handleMessage(channel, &me, network)
		n, _, err := Conn.ReadFromUDP(buf)
		channel <- buf[0:n]
		//fmt.Println("Connection recived: ", string(buf[0:n]), " \nfrom ", addr)

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	//build and send ping message

	network.pingResp = false
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
	//wait for timeout (2sec)
	time.Sleep(time.Second * 2)
	if network.pingResp {
		fmt.Println("Contact alive:", contact.Address)
		//network.pingChannel <- true
	} else {
		fmt.Println("Contact dead:", contact.Address)
		//network.pingChannel <- false
	}
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

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) updateRoutingtable(contact *Contact) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	bucket := network.rt.buckets[network.rt.getBucketIndex(contact.ID)]
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil { //If the contact is not in my bucket
		if bucket.list.Len() < bucketSize { //if my bucket has empty slots
			bucket.list.PushFront(contact)
		} else {
			lastContact := bucket.list.Back().Value.(Contact)
			go network.SendPingMessage(&lastContact)
			if !network.pingResp { //if I have no resonse add delete contact and add new
				bucket.RemoveContact(lastContact)
				bucket.AddContact(*contact)
			}
		}
	} else { //if I have the element move it to front
		bucket.list.MoveToFront(element)
	}
}
