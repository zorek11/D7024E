package d7024e

import (
	"fmt"
)

const count = 20
const alpha = 3

type Kademlia struct {
<<<<<<< HEAD
	rt *RoutingTable
	ht  map[KademliaID][]byte
=======
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
>>>>>>> 2e01575072c29155d4862dbf7c2399005958e8b0
}

func NewKademlia(self Contact) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
<<<<<<< HEAD
	kademlia.rt = NewRoutingTable(self)
	kademlia.ht = make(map[KademliaID][]byte)
	//kademlia.network = new(Network)
	//s := strings.Split(self.Address, ":")
=======
	kademlia.nt = NewNetwork(self, NewRoutingTable(self))
	kademlia.found = false
>>>>>>> 2e01575072c29155d4862dbf7c2399005958e8b0
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	kademlia.nt.AddMessage(target.ID)
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
	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			break
		}
		go kademlia.nt.SendFindContactMessage(&contacts[j])
	}
	//}
	for {
		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			//fmt.Println("Response in kademlia: ", kademlia.GetNetwork().GetResponse())
			//if kademlia.GetNetwork().GetResponse()[0] != nil {
			temp := kademlia.GetNetwork().GetResponse()[0]
			if temp[0].ID.String() == target.ID.String() {
				fmt.Println("This is the correct ID String: " + temp[0].ID.String())
				kademlia.found = true
				return
			} else {

				for i := 0; i < alpha; i++ {
					if i >= len(temp) {
						break
					}
					//fmt.Println("This is the new: ", temp[i])
					go kademlia.nt.SendFindContactMessage(&temp[i])

				}

				kademlia.nt.RemoveFirstResponse()

			}
			//}

		}
	}

}

func (kademlia *Kademlia) LookupData(hash string) {

	target := NewKademliaID(hash)

	kademlia.nt.AddMessage(target)
	contacts := kademlia.nt.rt.FindClosestContacts(target, count)
	fmt.Println(len(contacts))
	if contacts[0].ID == target {
		//TODO change to check local hash
		fmt.Println("Target found: " + target.String())
		fmt.Println("With address: " + contacts[0].String())
		return
	}

	//tempnetwork := NewNetwork(nt.rt, kademlia)

	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			break
		}
		go kademlia.nt.SendFindContactMessage(&contacts[j])
	}
	//}
	for {
		if kademlia.nt.GetData() != "" {
			fmt.Println(kademlia.nt.GetData())
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			//fmt.Println("Response in kademlia: ", kademlia.GetNetwork().GetResponse())
			//if kademlia.GetNetwork().GetResponse()[0] != nil {
			temp := kademlia.GetNetwork().GetResponse()[0]
			if temp[0].ID.String() == target.String() {
				fmt.Println("This is the correct ID String: " + temp[0].ID.String())
				kademlia.found = true
				return
			} else {

				for i := 0; i < alpha; i++ {
					if i >= len(temp) {
						break
					}
					//fmt.Println("This is the new: ", temp[i])
					go kademlia.nt.SendFindContactMessage(&temp[i])

				}

				kademlia.nt.RemoveFirstResponse()

			}
			//}

		}
	}

<<<<<<< HEAD
func (kademlia *Kademlia) LookupData(hash string) {
			data := kademlia.ht[key]
=======
>>>>>>> 2e01575072c29155d4862dbf7c2399005958e8b0
}

func (kademlia *Kademlia) Store(key string, value string) {
			kademlia.ht[key] = value


}
