package main

import (
	kademlia "D7024E-Kademlia/d7024e"
	"fmt"
	//"fmt"
)

//export GOPATH=$HOME/go
func main() {

	var input string

	contact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF11111111111111111111111111111111"),
		"127.0.0.1:7777")
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF11111111111111111111111111118888"),
		"127.0.0.1:7779")
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
		"127.0.0.1:9999")
	rt := kademlia.NewRoutingTable(me)
	rt.AddContact(contact)
	net := kademlia.NewNetwork(me)
	go kademlia.Listen(me)
	go kademlia.Listen(contact)
	go kademlia.Listen(contact2)
	for i := 1; i < 5; i++ {
		go net.SendPingMessage(&contact)
	}
	//go net.SendPingMessage(&contact)
	//go net.SendPingMessage(&contact2)
	/*
		for i := 1; i < 10; i++ {
			rt.AddContact(kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost:800"+strconv.Itoa(i)))
		}
		contacts := rt.FindClosestContacts(me.ID, 20)
		for i := range contacts {
			fmt.Println(contacts[i].String())
		}

	*/
	fmt.Scanln(&input)
	//for {

	//}

}
