package d7024e

import (
	"fmt"
	"strings"
	"testing"
)

func TestRoutingTable(t *testing.T) {

	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8000")
	contact2 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "127.0.0.1:8001")
	contact3 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "127.0.0.1:8002")
	contact4 := NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "127.0.0.1:8002")
	contact5 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "127.0.0.1:8002")
	contact6 := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "127.0.0.1:8002")

	kademlia := NewKademlia(contact1)
	kademlia2 := NewKademlia(contact2)

	kademlia.rt.AddContact(contact2)
	kademlia.rt.AddContact(contact3)
	kademlia.rt.AddContact(contact4)
	kademlia.rt.AddContact(contact5)
	kademlia.rt.AddContact(contact6)

	kademlia2.rt.AddContact(contact1)
	kademlia2.rt.AddContact(contact3)
	kademlia2.rt.AddContact(contact4)
	kademlia2.rt.AddContact(contact5)
	kademlia2.rt.AddContact(contact6)

	s1 := strings.Split(kademlia.rt.me.Address, ":")
	fmt.Println(s1)
	net := Network{}
	go net.SendPingMessage(&contact2)
	Listen("127.0.0.1", 8001)
	//

	//s2 := strings.Split(kademlia.rt.me.Address, ":")
	//Listen(s2[0], IntConverter(s2[1]))

	//n1 := new(Network)
	//PingMessage(&contact1)

	//contact := kademlia.LookupContact(&contact1)
	//fmt.Println(contact.ID)
}
