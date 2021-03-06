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
	me           Contact
	target       *KademliaID
	response     [][]Contact
	temp         *Contact
	rt           *RoutingTable
	mtx          *sync.Mutex
	dataFound    string
	pingList     []Ping
	pingResponse []*Contact
	storage      Storage
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
	fmt.Println("Added data: ", network.dataFound)
}

func (network *Network) GetData() string {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	return network.dataFound
}

func (network *Network) AddResponse(c []Contact) {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.response = append(network.response, [][]Contact{c}...)
	//fmt.Println("\nResponse: ", network.response)
}

func (network *Network) RemoveFirstResponse() {

	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.response = network.response[1:]
}

func (network *Network) GetResponse() [][]Contact {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	return network.response
}

func (network *Network) GetStorage() *Storage {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	return &network.storage
}

func (network *Network) GetRoutingTable() *RoutingTable {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	return network.rt
}

func (network *Network) ResetData() {
	network.mtx.Lock()
	defer network.mtx.Unlock()
	network.dataFound = ""
}

func (network *Network) PopData() string {
	temp := network.dataFound
	network.dataFound = ""
	return temp
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
	buf := make([]byte, 4096)
	for {
		time.Sleep(5 * time.Millisecond)
		n, _, err := Conn.ReadFromUDP(buf)
		if string(buf[:n]) == "KILL" {
			break
		}

		go messagehandler.handleMessage(channel, &me, network)
		channel <- buf[:n]

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
		time.Sleep(5 * time.Millisecond) //sleep to ofload the CPU and avoid buffer-error.
	}
}

/**
* Sends a ping message and waits for timeout.
 */
func (network *Network) SendPingMessage(contact *Contact) bool {

	//build and send message
	message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
	send(contact.Address, message)

	//wait for timeout (2sec)
	time.Sleep(time.Second * 2)
	network.mtx.Lock()
	if network.existsInPing(contact, network.pingResponse) {
		//fmt.Println("\nContact alive:", contact.Address)
		network.mtx.Unlock()
		return true
	} else {
		//fmt.Println("\nContact dead:", contact.Address)
		network.mtx.Unlock()
		return false
	}

}

/*
* Checks if a contacts is in the list or not, returns bool. If the contact exists, it deletes it.
 */
func (network *Network) existsInPing(c *Contact, contacts []*Contact) bool {
	for i := 0; i < len(contacts); i++ {
		if c.ID.Equals(contacts[i].ID) {
			contacts = append(contacts[:i], contacts[i+1:]...)
			return true
		}
	}
	return false
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	message := buildMessage([]string{"LookupContact", network.me.ID.String(), network.me.Address, network.target.String()})
	send(contact.Address, message)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) {
	message := buildMessage([]string{"LookupData", network.me.ID.String(), network.me.Address, hash})
	send(contact.Address, message)
}

func (network *Network) SendStoreMessage(contact *Contact, key *KademliaID, value string) {
	message := buildMessage([]string{"StoreData", network.me.ID.String(), network.me.Address, key.String(), value})
	send(contact.Address, message)
}

func (network *Network) SendPinMessage(contact *Contact, key *KademliaID) {
	message := buildMessage([]string{"Pin", network.me.ID.String(), network.me.Address, key.String()})
	send(contact.Address, message)
}

func (network *Network) SendUnpinMessage(contact *Contact, key *KademliaID) {
	message := buildMessage([]string{"Unpin", network.me.ID.String(), network.me.Address, key.String()})
	send(contact.Address, message)
}

/**
* Updates a routing table accoring to the Kademlia specs.
* Uses network.ping
 */
func (network *Network) UpdateRoutingtable(contact Contact) {
	//network.mtx.Lock()
	bucket := network.rt.buckets[network.rt.getBucketIndex(contact.ID)]
	var element *list.Element
	bucket.mtx.Lock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}
	bucket.mtx.Unlock()

	if element == nil { //If the contact is not in my bucket
		bucket.mtx.Lock()
		if bucket.list.Len() < bucketSize { //if my bucket has empty slots
			bucket.list.PushFront(contact)
			bucket.mtx.Unlock()
		} else {
			lastContact := bucket.list.Back().Value.(Contact)
			bucket.mtx.Unlock()
			//network.mtx.Unlock()
			ping := network.SendPingMessage(&lastContact)
			//network.mtx.Lock()
			if !ping { //if I have no resonse add delete contact and add new
				bucket.RemoveContact(lastContact)
				bucket.AddContact(contact)
			}
		}
	} else { //if I have the element move it to front
		bucket.mtx.Lock()
		bucket.list.MoveToFront(element)
		bucket.mtx.Unlock()
	}
}

/**
* Sends a protobuf message to the address via UDP
 */
func send(Address string, message *protobuf.KademliaMessage) {
	if len(Address) >= 14 {
		//fmt.Println("send to anddress: ", Address)
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
