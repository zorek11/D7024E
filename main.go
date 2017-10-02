package main

import (
	kademlia "D7024E-Kademlia/d7024e"
	"fmt"
	"sync"
	//"fmt"
)

//export GOPATH=$HOME/go
func main() {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	contact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF11111111111111111111111111111111"),
		"127.0.0.1:7777")
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF11111111111111111111111111111112"),
		"127.0.0.1:7778")
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF11111111111111111111111111111113"),
		"127.0.0.1:7779")

	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
		"127.0.0.1:9999")

	//rt := kademlia.NewRoutingTable(me)
	//rt.AddContact(contact)
	//rt.AddContact(contact2)
	me.CalcDistance(contact.ID)
	contact.CalcDistance(me.ID)
	me.CalcDistance(contact2.ID)

	kadc := kademlia.NewKademlia(contact)
	kadc.AddRoutingtable(me)
	kadc.AddRoutingtable(contact2)
	kadc.AddRoutingtable(contact3)
	contacts1 := kadc.GetRoutingtable().FindClosestContacts(contact.ID, 20)
	for i := range contacts1 {
		fmt.Println(contacts1[i].String())
	}

	netc := kademlia.NewNetwork(contact, kadc)
	netc.SendPingMessage(&contact)

	kad := kademlia.NewKademlia(me)
	net := kademlia.NewNetwork(me, kad)
	go net.Listen(me)
	go net.Listen(contact)
	go net.Listen(contact2)
	go kadc.LookupContact(&contact)
	/*for i := 0; i < 5; i++ {
		go net.SendPingMessage(&contact)
		go net.SendPingMessage(&contact3)
	}
	/*
		//rt := kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
		for i := 1; i < 10; i++ {
			rt.AddContact(kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost:800"+strconv.Itoa(i)))
		}
		contacts := rt.FindClosestContacts(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)
		for i := range contacts {
			fmt.Println(contacts[i].String())
		}
	*/
	for {
	}
}
