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

	contact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFFFF11111111111111111111111111111"),
		"127.0.0.1:1991")
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111112"),
		"127.0.0.1:1722")
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111113"),
		"127.0.0.1:1723")

	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
		"127.0.0.1:9999")

	//rt := kademlia.NewRoutingTable(me)
	//rt.AddContact(contact)
	//rt.AddContact(contact2)
	me.CalcDistance(contact.ID)
	contact.CalcDistance(me.ID)
	me.CalcDistance(contact2.ID)

	kad1 := kademlia.NewKademlia(me)
	net1 := kademlia.NewNetwork(me, kad1)
	//kad1.AddRoutingtable(contact)
	kad1.AddRoutingtable(contact3)

	kad2 := kademlia.NewKademlia(contact2)
	net2 := kademlia.NewNetwork(contact2, kad2)
	kad2.AddRoutingtable(me)
	kad2.AddRoutingtable(contact3)

	kad3 := kademlia.NewKademlia(contact3)
	net3 := kademlia.NewNetwork(contact3, kad3)
	kad3.AddRoutingtable(me)
	kad3.AddRoutingtable(contact2)
	/*
		contacts1 := kad1.GetRoutingtable().FindClosestContacts(contact.ID, 20)
		for i := range contacts1 {
			fmt.Println(contacts1[i].String())
		}
	*/
	//netc.SendPingMessage(&contact)

	go net1.Listen(me)
	go net2.Listen(contact2)
	go net3.Listen(contact3)
	// go net1.SendPingMessage(&contact2)
	go net1.GetKademlia().LookupContact(&contact2)
	/*for i := 0; i < 5; i++ {
		go net1.SendPingMessage(&contact2)
		//go net.SendPingMessage(&contact3)
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
