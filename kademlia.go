package d7024e

const count = 20
const alpha = 3

type Kademlia struct {
	rt *RoutingTable
}

func NewKademlia(self Contact, rt *RoutingTable) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
	kademlia.rt = NewRoutingTable(self)
	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
	contacts := kademlia.rt.FindClosestContacts(target.ID, count)
	if target.ID != contacts[0].ID {
		for i := 0; i < alpha; i++ {
			//go LookupContact()
		}
	} else {
		return contacts[0]
	}
	return contacts[0]
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
