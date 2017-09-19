package main

import (
	kademlia "D7024E-Kademlia/d7024e"
	"sync"
	//"fmt"
)

//export GOPATH=$HOME/go
func main() {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	contact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8001")
	net := kademlia.Network{this}
	go net.SendPingMessage(&contact)

	kademlia.Listen("127.0.0.1", 8002)
	/*
		rt := kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
		for i := 1; i < 10; i++ {
			rt.AddContact(kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost:800"+strconv.Itoa(i)))
		}
		contacts := rt.FindClosestContacts(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)
		for i := range contacts {
			fmt.Println(contacts[i].String())
		}
	*/
}