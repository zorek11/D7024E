package d7024e

import (
	"crypto/sha1"
	"fmt"
	"time"
)

const count = 20
const alpha = 3

type Kademlia struct {
	nt    Network
	items []string
	found bool
	start time.Time
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

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	kademlia.nt.AddMessage(target.ID)
	contacts := kademlia.nt.rt.FindClosestContacts(target.ID, count)
	fmt.Println(len(contacts))
	result := make([]Contact, 20)
	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			fmt.Println("BREAK", j)
			break
		}
		result[j] = contacts[j]
		go kademlia.nt.SendFindContactMessage(&contacts[j])
	}
	result = result[0:len(contacts)]
	for i := 0; i < len(contacts); i++ {
		result[i] = contacts[i]
	}
	kademlia.start = time.Now()
	t := time.Now()
	for {
		t = time.Now()
		//fmt.Println(len(kademlia.GetNetwork().GetResponse()))

		if t.Sub(kademlia.start) > 5000000000 {
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			fmt.Println("we got the results")
			return result
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {

			temp := kademlia.GetNetwork().GetResponse()[0]
			tempAlpha := alpha
			for i := 0; i < tempAlpha; i++ {
				if i >= len(temp) {
					break
				}
				if existsIn(temp[i], result) {
					tempAlpha++
				} else {
					go kademlia.nt.SendFindContactMessage(&temp[i])
				}
			}
			result = kademlia.checkContacts(result, temp)
			fmt.Println("\n\nthis is the result so far: ", result)
			kademlia.nt.RemoveFirstResponse()
		}
	}
}

func existsIn(c Contact, contacts []Contact) bool {
	for i := 0; i < len(contacts); i++ {
		if c.ID.Equals(contacts[i].ID) {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) checkContacts(this []Contact, addition []Contact) []Contact {
	for j := 0; j < len(addition); j++ {
		addition[j].CalcDistance(kademlia.nt.target)
	}
	var temp ContactCandidates
	temp.Append(this)
	temp.Append(addition)
	temp.Sort()
	k := 0
	for k < count && k < len(temp.contacts)-1 {
		if temp.contacts[k].ID.Equals(temp.contacts[k+1].ID) {
			temp.contacts = append(temp.contacts[:k], temp.contacts[k+1:]...)

		} else {
			kademlia.start = time.Now()
			k++
		}
	}
	if len(temp.contacts) < count {
		return temp.contacts
	}
	return temp.contacts[0:count]
}

//TODO: Implement some kind of deletion if timestamp overdue. (PURGE)
func (kademlia *Kademlia) LookupData(hash string) {
	target := KademliaID(sha1.Sum([]byte(hash)))

	if len(kademlia.nt.storage.RetrieveFile(&target)) > 0 {
		fmt.Println("File retrieved in LookupData: " + kademlia.nt.storage.RetrieveFile(&target))
		fmt.Println("File found locally in LookupData: " + target.String())
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

//TODO: call Store again after a specific time to store again(REPUBLISH)
func (kademlia *Kademlia) Store(data string) KademliaID {
	//TODO: LookupContact find 20 closest somehow. This kademlia doesn't know all contacts in network.
	hashdata := []byte(data)
	key := KademliaID(sha1.Sum(hashdata))
	//contacts := kademlia.LookupContact(&target) //How it should work
	contacts := kademlia.nt.rt.FindClosestContacts(&key, count)

	for j := range contacts {
		go kademlia.nt.SendStoreMessage(&contacts[j], &key, data)
	}
	return key
}
