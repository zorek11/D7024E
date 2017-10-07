package d7024e

import (
	"crypto/sha1"
	"fmt"
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
	kademlia.nt = NewNetwork(self, NewRoutingTable(self), NewStorage())
	kademlia.found = false
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
	target := KademliaID(sha1.Sum([]byte(hash)))

	if len(kademlia.nt.storage.RetrieveFile(&target)) > 0 {
		fmt.Println("File: " + kademlia.nt.storage.RetrieveFile(&target))
		fmt.Println("File found locally: " + target.String())
		return
	}
	for {

	}
	kademlia.nt.AddMessage(&target)
	contacts := kademlia.nt.rt.FindClosestContacts(&target, count)
	fmt.Println(len(contacts))

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

}

func (kademlia *Kademlia) Store(data string) {
	hashdata := []byte(data)
	key := KademliaID(sha1.Sum(hashdata))

	contacts := kademlia.nt.rt.FindClosestContacts(&key, count)

	for j := range contacts {
		go kademlia.nt.SendStoreMessage(&contacts[j], &key, data)
	}
}
