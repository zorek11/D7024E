package main

import (
	kademlia "D7024E-Kademlia/d7024e"
	//"fmt"
)

//export GOPATH=$HOME/go
func main() {
	//var mutex = &sync.Mutex{}
	//mutex.Lock()
	//defer mutex.Unlock()

	//contact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFFFF11111111111111111111111111111"),
	//	"127.0.0.1:1191")
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111192"),
		"127.0.0.1:1222")
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111193"),
		"127.0.0.1:1223")
	contact4 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111194"),
		"127.0.0.1:1224")
	contact5 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111195"),
		"127.0.0.1:1225")
	contact6 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111196"),
		"127.0.0.1:1226")
	contact7 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111197"),
		"127.0.0.1:1227")

	/*
		contact8 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111198"),
			"127.0.0.1:1228")
		contact9 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111199"),
			"127.0.0.1:1229")
		contact10 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111110"),
			"127.0.0.1:1210")
		contact11 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111111"),
			"127.0.0.1:1211")
		contact12 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111112"),
			"127.0.0.1:1212")
	*/

	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
		"127.0.0.1:1999")

	//rt := kademlia.NewRoutingTable(me)
	//rt.AddContact(contact)
	//rt.AddContact(contact2)
	//me.CalcDistance(contact.ID)
	//contact.CalcDistance(me.ID)
	//me.CalcDistance(contact2.ID)

	kad1 := kademlia.NewKademlia(me)
	//net1 := kademlia.NewNetwork(me, kad1)
	//kad1.AddRoutingtable(contact)
	kad1.AddRoutingtable(contact2)

	kad2 := kademlia.NewKademlia(contact2)
	//net2 := kademlia.NewNetwork(contact2, kad2)
	kad2.AddRoutingtable(me)
	kad2.AddRoutingtable(contact3)
	kad2.AddRoutingtable(contact4)

	kad3 := kademlia.NewKademlia(contact3)
	//net3 := kademlia.NewNetwork(contact3, kad3)
	kad3.AddRoutingtable(contact4)
	kad3.AddRoutingtable(contact2)
	kad3.AddRoutingtable(me)

	kad4 := kademlia.NewKademlia(contact4)
	//net3 := kademlia.NewNetwork(contact3, kad3)
	kad4.AddRoutingtable(contact3)
	kad4.AddRoutingtable(contact5)
	//kad4.AddRoutingtable(contact12)

	kad5 := kademlia.NewKademlia(contact5)
	//net3 := kademlia.NewNetwork(contact3, kad3)
	kad5.AddRoutingtable(contact4)
	kad5.AddRoutingtable(contact6)
	//kad5.AddRoutingtable(contact11)

	kad6 := kademlia.NewKademlia(contact6)
	//net3 := kademlia.NewNetwork(contact3, kad3)
	kad6.AddRoutingtable(contact5)
	kad6.AddRoutingtable(contact7)
	//kad6.AddRoutingtable(contact10)

	kad7 := kademlia.NewKademlia(contact7)
	//net3 := kademlia.NewNetwork(contact3, kad3)
	kad7.AddRoutingtable(contact6)
	//kad7.AddRoutingtable(contact8)
	//kad7.AddRoutingtable(contact9)

	/*
		contacts1 := kad1.GetRoutingtable().FindClosestContacts(contact.ID, 20)
		for i := range contacts1 {
			fmt.Println(contacts1[i].String())
		}
	*/
	//netc.SendPingMessage(&contact)

	net1 := kad1.GetNetwork()

	net2 := kad2.GetNetwork()

	net3 := kad3.GetNetwork()

	net4 := kad4.GetNetwork()

	net5 := kad5.GetNetwork()

	net6 := kad6.GetNetwork()

	//net7 := kad7.GetNetwork()

	go net1.Listen(me)
	go net2.Listen(contact2)
	go net3.Listen(contact3)
	go net4.Listen(contact4)
	go net5.Listen(contact5)
	go net6.Listen(contact6)
	//go net7.Listen(contact7)
	// go net1.SendPingMessage(&contact2)
	//net1.AddMessage(&contact2)
	//go kad1.LookupContact(&contact7)

	//go net1.SendFindContactMessage(&contact2)
	for i := 0; i < 2; i++ {
		go net1.SendPingMessage(&contact2)
		go net1.SendPingMessage(&contact7)
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
	/*for {
		if kad1.GetFound() != false {
			fmt.Println("END")
			break
		}
	}*/
	for {
	}
}
