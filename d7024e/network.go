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
	me        Contact
	target    *KademliaID
	response  [][]Contact
	temp      *Contact
	rt        *RoutingTable
	mtx       *sync.Mutex
	dataFound string
	pingList  []Ping
	storage   Storage
}

type Ping struct {
	Address  string
	Response bool
	Done     chan bool
	Queue    int
}

func NewNetwork(me Contact, rt *RoutingTable, st Storage) Network {
	network := Network{}
	network.me = me
	network.rt = rt
	network.mtx = &sync.Mutex{}
	network.dataFound = ""
	network.storage = st
	return network
}

func (network *Network) AddMessage(c *KademliaID) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.target = c
}

func (network *Network) AddData(s string) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.dataFound = s
}

func (network *Network) GetData() string {
	return network.dataFound
}

func (network *Network) AddResponse(c []Contact) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.response = append(network.response, [][]Contact{c}...)
	fmt.Println("\nResponse: ", network.response)
	fmt.Println(" in : ", network.me)
}

func (network *Network) RemoveFirstResponse() {

	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.response = network.response[1:]
}

func (network *Network) GetResponse() [][]Contact {
	return network.response
}

func (network *Network) GetStorage() *Storage {
	return &network.storage
}

func (network *Network) GetRoutingTable() *RoutingTable {
	return network.rt
}

/**
* Establishes a UDP-listner on an adress and handles incoming traffic in a differnt go-routine
 */
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

/**
* Sends a ping message and waits for timeout.
 */
func (network *Network) SendPingMessage(contact *Contact) {
	//build and send ping message
	//make sure to send ping in order and wait for response.
	pingIndex := IndexInSlice(contact.Address, network.pingList)
	if pingIndex > -1 { //if address in list
		network.pingList[pingIndex].Queue = network.pingList[pingIndex].Queue + 1
		<-network.pingList[pingIndex].Done //wait for previous ping to finsh
	} else {
		network.pingList = append(network.pingList, Ping{contact.Address, false, make(chan bool, 1), 0}) //add to pingList
	}

	//build and send message
	message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
	send(contact.Address, message)

	//wait for timeout (2sec)
	time.Sleep(time.Second * 2)
	pingIndex = IndexInSlice(contact.Address, network.pingList)
	if network.pingList[pingIndex].Response {
		fmt.Println("\nContact alive:", contact.Address)
	} else {
		fmt.Println("\nContact dead:", contact.Address)
	}
	network.pingList[pingIndex].Queue = network.pingList[pingIndex].Queue - 1 //decrease queue
	if network.pingList[pingIndex].Queue >= 0 {                               //if there is someone in the queue release the channel
		network.pingList[pingIndex].Done <- true
		fmt.Println(network.pingList)
	} else {
		network.pingList = append(network.pingList[:pingIndex], network.pingList[pingIndex+1:]...) //delete contact from list
		fmt.Println(network.pingList)
	}

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	message := &protobuf.KademliaMessage{
		Label:      proto.String("LookupContact"),
		Senderid:   proto.String(network.me.ID.String()),
		SenderAddr: proto.String(network.me.Address),
		Lookupcontact: &protobuf.KademliaMessage_LookupContact{
			Id: proto.String(network.target.String()),
		},
	}
	send(contact.Address, message)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) {
	message := &protobuf.KademliaMessage{
		Label:      proto.String("LookupData"),
		Senderid:   proto.String(network.me.ID.String()),
		SenderAddr: proto.String(network.me.Address),
		Key:        proto.String(hash),
	}
	send(contact.Address, message)

}

func (network *Network) SendStoreMessage(contact *Contact, key *KademliaID, value string) {
	message := buildMessage([]string{"StoreData", network.me.ID.String(), network.me.Address, key.String(), value})
	send(contact.Address, message)
}

/**
* Updates a routing table accoring to the Kademlia specs.
* Uses network.ping
 */
func (network *Network) UpdateRoutingtable(contact Contact) {
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
			pingIndex := IndexInSlice(contact.Address, network.pingList)
			if network.pingList[pingIndex].Response { //if I have no resonse add delete contact and add new
				bucket.RemoveContact(lastContact)
				bucket.AddContact(contact)
			}
		}
	} else { //if I have the element move it to front
		bucket.list.MoveToFront(element)
	}
}

/**
* Sends a protobuf message to the address via UDP
 */
func send(Address string, message *protobuf.KademliaMessage) {
	data, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("Marshal Error: ", err)
	}

	Conn, err := net.Dial("udp", Address)
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	defer Conn.Close()
	_, err = Conn.Write(data)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
}

/**
* Finds the index of a elememt in a Pingslice. Returns -1 on nonexisting element.
 */
func IndexInSlice(Address string, list []Ping) int {
	i := 0
	for _, x := range list {
		if x.Address == Address {
			return i
		}
		i++
	}
	return -1
}
