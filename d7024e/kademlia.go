package d7024e

import (
	"crypto/sha1"
	"fmt"
	"sync"
	"time"
)

const count = 20
const alpha = 3

type Kademlia struct {
	nt    Network
	items []string
	found bool
	start time.Time
	mtx   *sync.Mutex
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
	kademlia.mtx = &sync.Mutex{}
	kademlia.nt = NewNetwork(self, NewRoutingTable(self), NewStorage())
	kademlia.found = false
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	var contacted []Contact
	var alphaContacts []Contact
	var result ContactCandidates
	kademlia.nt.AddMessage(target)
	contacts := kademlia.nt.rt.FindClosestContacts(target, count)
	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			break
		}
		contacts[j].CalcDistance(target)
		go kademlia.nt.SendFindContactMessage(&contacts[j])
		result.Append([]Contact{contacts[j]})

		contacted = append(contacted, []Contact{contacts[j]}...)
		alphaContacts = append(alphaContacts, []Contact{contacts[j]}...)
	}
	//fmt.Println("I get here")
	start := time.Now()
	t := time.Now()
	alphaTime := time.Now()
	checkedAlpha := 0
	checkedSecond := 0
	checkedThird := 0
	checkedFailure := 0
	//alphaCount := alpha
	//fmt.Println(result.contacts)
	for {
		//fmt.Println("second: ", checkedSecond, " third: ", checkedThird)
		t = time.Now()
		current := t.Sub(start).Nanoseconds()
		if t.Sub(alphaTime).Nanoseconds() > 3000000000 && checkedAlpha < len(alphaContacts) {
			if len(contacts) > len(alphaContacts) {
				go kademlia.nt.SendFindContactMessage(&contacts[len(alphaContacts)])
				alphaContacts = append(alphaContacts, []Contact{contacts[len(alphaContacts)]}...)
			}

		}
		if current > 5000000000 || checkedFailure >= len(alphaContacts) {
			fmt.Println("\nreturning", kademlia.nt.rt.me.ID.String(), "\n")
			kademlia.nt.rt.PrintRoutingTable()
			fmt.Println("")
			if len(alphaContacts) >= len(contacts) {
				if len(result.contacts) < count {
					return result.contacts[:len(result.contacts)]
				}
				return result.contacts[:count]
			}
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			start = time.Now()
			if existsIn(kademlia.GetNetwork().GetResponse()[0][0], alphaContacts) {
				alphaTime = time.Now()
				temp := sortContacts(kademlia.GetNetwork().GetResponse()[0][0], kademlia.GetNetwork().GetResponse()[0][1:])
				temp = kademlia.calcDist(temp)
				result.contacts = kademlia.checkContacts(result.contacts, temp)
				checkedAlpha++
				for i := 1; i < len(temp); i++ {
					if !(existsIn(temp[i], contacted)) &&
						!(temp[i].ID.Equals(kademlia.nt.rt.me.ID)) {
						go kademlia.nt.SendFindContactMessage(&temp[i])
						contacted = append(contacted, []Contact{temp[i]}...)
					}
				}
			} else if existsIn(kademlia.GetNetwork().GetResponse()[0][0], contacted) {
				//fmt.Println("in second loop")
				temp := sortContacts(kademlia.GetNetwork().GetResponse()[0][0], kademlia.GetNetwork().GetResponse()[0][1:])
				temp = kademlia.calcDist(temp)
				result.contacts = kademlia.checkContacts(result.contacts, temp)
				temporaryChecked := checkedSecond
				for i := 1; i < len(temp); i++ {
					if !(existsIn(temp[i], contacted)) &&
						!(temp[i].ID.Equals(kademlia.nt.rt.me.ID)) {
						checkedSecond++
						go kademlia.nt.SendFindContactMessage(&temp[i])
						contacted = append(contacted, []Contact{temp[i]}...)
					}
					if temporaryChecked >= checkedSecond {
						checkedFailure++
					}

				}
			} else {
				checkedThird++
				temp := sortContacts(kademlia.GetNetwork().GetResponse()[0][0], kademlia.GetNetwork().GetResponse()[0][1:])
				result.contacts = kademlia.checkContacts(result.contacts, temp)
			}
			kademlia.nt.RemoveFirstResponse()
		}
	}
}

func (kademlia *Kademlia) calcDist(contacts []Contact) []Contact {
	for i := 0; i < len(contacts); i++ {
		contacts[i].CalcDistance(kademlia.nt.target)
	}
	return contacts
}

func sortContacts(c Contact, contacts []Contact) []Contact {
	var temp []Contact
	for i := 0; i < len(contacts); i++ {
		if contacts[i].ID.Less(c.ID) {
			temp = append(temp, []Contact{contacts[i]}...)
		}
	}
	return temp
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
	for o := 0; o < len(addition); o++ {
		if !(existsIn(addition[o], this)) {
			temp.Append([]Contact{addition[o]})
		}
	}
	temp.Sort()
	kademlia.checkDuplicates(this, temp.contacts)

	if len(temp.contacts) < count {
		return temp.contacts
	}
	result := temp.contacts[0:count]
	return result
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

func (kademlia *Kademlia) LookupData(hash string) string {
	target := NewKademliaID(hash)
	var contacted []Contact
	var alphaContacts []Contact
	var result ContactCandidates
	kademlia.nt.AddMessage(target)
	contacts := kademlia.nt.rt.FindClosestContacts(target, count)
	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			break
		}
		contacts[j].CalcDistance(target)
		go kademlia.nt.SendFindDataMessage(hash, &contacts[j])
		result.Append([]Contact{contacts[j]})

		contacted = append(contacted, []Contact{contacts[j]}...)
		alphaContacts = append(alphaContacts, []Contact{contacts[j]}...)
	}
	//fmt.Println("I get here")
	start := time.Now()
	t := time.Now()
	alphaTime := time.Now()
	checkedAlpha := 0
	checkedSecond := 0
	checkedThird := 0
	checkedFailure := 0
	//alphaCount := alpha
	//fmt.Println(result.contacts)
	for {
		//fmt.Println("second: ", checkedSecond, " third: ", checkedThird)
		t = time.Now()
		current := t.Sub(start).Nanoseconds()
		if len(kademlia.nt.GetData()) > 0 {
			return kademlia.nt.GetData()
		}
		if t.Sub(alphaTime).Nanoseconds() > 1000000000 && checkedAlpha < len(alphaContacts) {
			if len(contacts) > len(alphaContacts) {
				go kademlia.nt.SendFindDataMessage(hash, &contacts[len(alphaContacts)])
				alphaContacts = append(alphaContacts, []Contact{contacts[len(alphaContacts)]}...)
			}

		}
		if current > 3000000000 || checkedFailure >= len(alphaContacts) {
			fmt.Println("\nreturning", kademlia.nt.rt.me.ID.String(), "\n", result.contacts, "\n-")
			if len(alphaContacts) >= len(contacts) {
				if len(result.contacts) < count {
					return kademlia.nt.GetData()

				}
				return kademlia.nt.GetData()

			}
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			start = time.Now()
			if existsIn(kademlia.GetNetwork().GetResponse()[0][0], alphaContacts) {
				alphaTime = time.Now()
				temp := sortContacts(kademlia.GetNetwork().GetResponse()[0][0], kademlia.GetNetwork().GetResponse()[0][1:])
				temp = kademlia.calcDist(temp)
				result.contacts = kademlia.checkContacts(result.contacts, temp)
				checkedAlpha++
				for i := 1; i < len(temp); i++ {
					if !(existsIn(temp[i], contacted)) &&
						!(temp[i].ID.Equals(kademlia.nt.rt.me.ID)) {
						go kademlia.nt.SendFindDataMessage(hash, &temp[i])
						contacted = append(contacted, []Contact{temp[i]}...)
					}
				}
			} else if existsIn(kademlia.GetNetwork().GetResponse()[0][0], contacted) {
				//fmt.Println("in second loop")
				temp := sortContacts(kademlia.GetNetwork().GetResponse()[0][0], kademlia.GetNetwork().GetResponse()[0][1:])
				temp = kademlia.calcDist(temp)
				result.contacts = kademlia.checkContacts(result.contacts, temp)
				temporaryChecked := checkedSecond
				for i := 1; i < len(temp); i++ {
					if !(existsIn(temp[i], contacted)) &&
						!(temp[i].ID.Equals(kademlia.nt.rt.me.ID)) {
						checkedSecond++
						go kademlia.nt.SendFindDataMessage(hash, &temp[i])
						contacted = append(contacted, []Contact{temp[i]}...)
					}
					if temporaryChecked >= checkedSecond {
						checkedFailure++
					}

				}
			} else {
				checkedThird++
				temp := sortContacts(kademlia.GetNetwork().GetResponse()[0][0], kademlia.GetNetwork().GetResponse()[0][1:])
				result.contacts = kademlia.checkContacts(result.contacts, temp)
			}
			kademlia.nt.RemoveFirstResponse()
		}
	}

}

//TODO: call Store again after a specific time to store again(REPUBLISH)
func (kademlia *Kademlia) Store(data string) {
	fmt.Println("IM totally gonna store: ", data)
	//TODO: LookupContact find 20 closest somehow. This kademlia doesn't know all contacts in network.
	hashdata := []byte(data)
	key := KademliaID(sha1.Sum(hashdata))
	//contacts := kademlia.LookupContact(&target) //How it should work
	//kademlia.nt.mtx.Lock()
	//contacts := kademlia.nt.rt.FindClosestContacts(&key, count)
	contacts := kademlia.LookupContact(&key)
	fmt.Println("AFTER LOOKUP IN STORE")
	//kademlia.nt.mtx.Unlock()
	for j := range contacts {
		fmt.Println(contacts[j])
		go kademlia.nt.SendStoreMessage(&contacts[j], &key, data)
	}
	fmt.Println("DONE")
}

func (kademlia *Kademlia) Pin(target string) {
	key := NewKademliaID(target)
	contacts := kademlia.LookupContact(key)

	for j := range contacts {
		go kademlia.nt.SendPinMessage(&contacts[j], key)
	}
}

func (kademlia *Kademlia) Unpin(target string) {
	key := NewKademliaID(target)
	contacts := kademlia.LookupContact(key)

	for j := range contacts {
		go kademlia.nt.SendUnpinMessage(&contacts[j], key)
	}
}
