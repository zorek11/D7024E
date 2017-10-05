package d7024e

import (
	"fmt"
	"sync"
)

const count = 20
const alpha = 3

type Kademlia struct {
	nt    Network
	items []string
	found bool
}

func (kademlia *Kademlia) AddRoutingtable(c Contact) {
	kademlia.nt.rt.AddContact(c)
}

func (kademlia *Kademlia) GetRoutingtable() *RoutingTable {
	return kademlia.nt.rt
}
func (kademlia *Kademlia) GetFound() bool {
	return kademlia.found
}

func (kademlia *Kademlia) GetNetwork() *Network {
	return &kademlia.nt
}

func NewKademlia(self Contact) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
	kademlia.nt = NewNetwork(self, NewRoutingTable(self))
	kademlia.found = false
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	var mutex = &sync.Mutex{}

	contacts := kademlia.nt.rt.FindClosestContacts(target.ID, count)
	//thisalpha := alpha % (len(contacts) + 1)
	fmt.Println(len(contacts))
	if contacts[0].ID == target.ID {
		//return &contacts[0]
		fmt.Println("Target found: " + target.String())
		fmt.Println("With address: " + contacts[0].String())
		return
	}

	//tempnetwork := NewNetwork(nt.rt, kademlia)
	kademlia.nt.AddMessage(target)

	go kademlia.nt.SendFindContactMessage(&contacts[0])
	//}
	for {

		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			//fmt.Println("Response in kademlia: ", kademlia.GetNetwork().GetResponse())
			//if kademlia.GetNetwork().GetResponse()[0] != nil {
			temp := kademlia.GetNetwork().GetResponse()[0]
			fmt.Println("this is first contact: ", temp)
			if temp[0].ID.String() == target.ID.String() {
				fmt.Println("This is the correct ID String: " + temp[0].ID.String())
				kademlia.found = true
				return
			} else {

				mutex.Lock()
				for i := 0; i < alpha; i++ {
					if i >= len(temp) {
						break
					}
					//fmt.Println("This is the new: ", temp[i])
					go kademlia.nt.SendFindContactMessage(&temp[i])

				}

				kademlia.nt.RemoveFirstResponse()
				mutex.Unlock()

			}
			//}

		}
	}

}

func (kademlia *Kademlia) LookupData(hash string) {

	target := NewContact(NewKademliaID(hash), "")
	contacts := kademlia.nt.rt.FindClosestContacts(target.ID, count)
	//thisalpha := alpha % (len(contacts) + 1)
	fmt.Println(len(contacts))
	if contacts[0].ID == target.ID {
		//return &contacts[0]
		fmt.Println("Target found: " + target.String())
		fmt.Println("With address: " + contacts[0].String())
		return
	}

	//tempnetwork := NewNetwork(nt.rt, kademlia)
	kademlia.nt.AddMessage(&target)

	go kademlia.nt.SendFindContactMessage(&contacts[0])
	//}
	x := 0
	for {
		if x < len(kademlia.GetNetwork().GetResponse()) {
			//fmt.Println("Response in kademlia: ", kademlia.GetNetwork().GetResponse())
			//if kademlia.GetNetwork().GetResponse()[0] != nil {
			temp := kademlia.GetNetwork().GetResponse()[x]
			x++
			if temp[0].ID.String() == target.ID.String() {
				fmt.Println("This is the correct ID String: " + temp[0].ID.String())
				kademlia.found = true
				return
			} else {
				var mutex = &sync.Mutex{}
				mutex.Lock()
				for i := 0; i < alpha; i++ {
					if i >= len(temp) {
						break
					}
					//fmt.Println("This is the new: ", temp[i])
					go kademlia.nt.SendFindContactMessage(&temp[i])

				}

				//kademlia.nt.RemoveFirstResponse()
				mutex.Unlock()

			}
			//}

		}
	}

}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
