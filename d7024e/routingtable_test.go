package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	kademlia := NewKademlia(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	kademlia.rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	kademlia.rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	kademlia.rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	kademlia.rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	kademlia.rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	contact1 := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")
	kademlia.rt.AddContact(contact1)

	contacts := kademlia.rt.FindClosestContacts(NewKademliaID("1111111500000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}
