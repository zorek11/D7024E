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

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	fmt.Println("----------------------------------------------------------------------")
	contacted := make([]Contact, 0)

	kademlia.nt.AddMessage(target)
	contacts := kademlia.nt.rt.FindClosestContacts(target, count)
	fmt.Println(len(contacts))
	result := make([]Contact, 20)
	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			fmt.Println("BREAK", j)
			break
		}
		result[j] = contacts[j]
		go kademlia.nt.SendFindContactMessage(&contacts[j])
		contacted = append(contacted, []Contact{contacts[j]}...)
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

		if t.Sub(kademlia.start) > 1000000000 {
			fmt.Println("we got the timeout")
			fmt.Println("\nhere is the routing table")
			kademlia.nt.rt.PrintRoutingTable()
			return result
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {

			temp := kademlia.GetNetwork().GetResponse()[0]
			tempAlpha := alpha
			result = kademlia.checkContacts(result, temp)
			for i := 0; i < tempAlpha; i++ {
				if i >= len(result) {
					break
				}
				if existsIn(result[i], contacted) || result[i].ID.Equals(kademlia.nt.rt.me.ID) {
					tempAlpha++
				} else {
					go kademlia.nt.SendFindContactMessage(&result[i])
					contacted = append(contacted, []Contact{result[i]}...)
				}
			}
			if tempAlpha == 20 {
				fmt.Println("we got the result")
				fmt.Println("\nhere is the routing table")
				kademlia.nt.rt.PrintRoutingTable()
				return result
			}

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
	kademlia.checkDuplicates(this, temp.contacts)
	k := 0
	for k < count && k < len(temp.contacts)-1 {
		if temp.contacts[k].ID.Equals(temp.contacts[k+1].ID) {
			temp.contacts = append(temp.contacts[:k], temp.contacts[k+1:]...)
			k++
		} else {
			//kademlia.start = time.Now()
			k++
		}
	}
	kademlia.checkDuplicates(this, temp.contacts)
	if len(temp.contacts) < count {
		return temp.contacts
	}
	return temp.contacts[0:count]
}

func (kademlia *Kademlia) checkDuplicates(contacts []Contact, temp []Contact) {
	if len(contacts) < count {
		return
	}
	for i := 0; i < len(contacts); i++ {
		if !(contacts[i].ID.Equals(temp[i].ID)) {
			kademlia.start = time.Now()
			return
		}
	}

}

//TODO: Implement some kind of deletion if timestamp overdue. (PURGE)
func (kademlia *Kademlia) LookupData(hash string) string {
	target := NewKademliaID(hash)
	kademlia.nt.AddMessage(target)
	if len(kademlia.nt.storage.RetrieveFile(target)) > 0 {
		fmt.Println("found target locally")
		return kademlia.nt.storage.RetrieveFile(target)
	}
	contacted := make([]Contact, 0)
	contacts := kademlia.nt.rt.FindClosestContacts(target, count)
	fmt.Println(len(contacts))
	result := make([]Contact, 20)
	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			fmt.Println("BREAK", j)
			break
		}
		result[j] = contacts[j]
		go kademlia.nt.SendFindDataMessage(hash, &contacts[j])
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
		if len(kademlia.nt.GetData()) > 0 {
			fmt.Println("we got the data: ", kademlia.nt.GetData())
			return kademlia.nt.GetData()
		}
		if t.Sub(kademlia.start) > 1000000000 {
			fmt.Println("we got the timeout")
			fmt.Println("\nhere is the routing table")
			kademlia.nt.rt.PrintRoutingTable()
			return kademlia.nt.GetData()
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {

			temp := kademlia.GetNetwork().GetResponse()[0]
			tempAlpha := alpha
			result = kademlia.checkContacts(result, temp)
			for i := 0; i < tempAlpha; i++ {
				if i >= len(result) {
					break
				}
				if existsIn(result[i], contacted) || result[i].ID.Equals(kademlia.nt.rt.me.ID) {
					tempAlpha++
				} else {
					go kademlia.nt.SendFindDataMessage(hash, &result[i])
					contacted = append(contacted, []Contact{result[i]}...)
				}
			}
			if tempAlpha == 20 {
				fmt.Println("we looked through all")
				fmt.Println("\nhere is the routing table")
				kademlia.nt.rt.PrintRoutingTable()
				return kademlia.nt.GetData()
			}

			fmt.Println("\n\nthis is the result so far: ", result)
			kademlia.nt.RemoveFirstResponse()
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
