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

	kademlia.nt.AddMessage(target)
	contacts := kademlia.nt.rt.FindClosestContacts(target, count)
	result := contacts

	for j := 0; j < alpha; j++ {
		if j >= len(contacts) {
			//fmt.Println("BREAK", j)
			break
		}
		//result[j] = contacts[j]
		go kademlia.nt.SendFindContactMessage(&contacts[j])
		contacted = append(contacted, []Contact{contacts[j]}...)
	}

	//result = result[0:len(contacts)]
	//	for i := 0; i < len(contacts); i++ {
	//	result[i] = contacts[i]
	//}
	//fmt.Println("I get here")
	kademlia.start = time.Now()
	t := time.Now()
	same := 0
	for {
		t = time.Now()
		//fmt.Println(len(kademlia.GetNetwork().GetResponse()))

		if t.Sub(kademlia.start).Nanoseconds() > 10000000000 {
			//fmt.Println("we got the timeout")
			//fmt.Println("\nhere is the routing table")
			//kademlia.nt.mtx.Lock()
			//kademlia.nt.rt.PrintRoutingTable()
			//kademlia.nt.mtx.Unlock()
			fmt.Println("If 1")
			return result
		}
		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			//fmt.Println("If 2")
			//fmt.Println(kademlia.GetNetwork().GetResponse())
			kademlia.start = time.Now()
			tempAlpha := alpha
			temp := kademlia.nt.GetResponse()[0]
			result = kademlia.checkContacts(result, temp)
			if len(result) < alpha {
				tempAlpha = len(result)
			}
			for i := 0; i < tempAlpha && i < len(result) && tempAlpha < count; i++ {
				if existsIn(result[i], contacted) || result[i].ID.Equals(kademlia.nt.rt.me.ID) {
					if tempAlpha < len(result) {
						tempAlpha++
					}

				} else {
					go kademlia.nt.SendFindContactMessage(&result[i])
					contacted = append(contacted, []Contact{result[i]}...)
				}
			}
			kademlia.nt.RemoveFirstResponse()
			if tempAlpha >= count {
				//fmt.Println("--------------------------------------we got the result for: ", kademlia.nt.rt.me.String())
				//fmt.Println("and result: ", result)
				return result
			} else if tempAlpha >= len(result) {
				same++
				if same > 5 {
					//fmt.Println("--------------------------------------we got the result for: ", kademlia.nt.rt.me.String())
					//fmt.Println("and result: ", result)
					return result
				} else if !(len(kademlia.GetNetwork().GetResponse()) > 0) {
					time.Sleep(200 * time.Millisecond)
					if !(len(kademlia.GetNetwork().GetResponse()) > 0) {
						//fmt.Println("---------------------------len--------we got the result for: ", kademlia.nt.rt.me.String())
						//fmt.Println("and result: ", result)
						return result
					}
				}
			}
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
	fmt.Println("--------------------------------LOOKUP DATA--------------------------------------")
	//fmt.Println("----------------------------------------------------------------------")
	kademlia.nt.ResetData()
	contacted := make([]Contact, 0)
	target := NewKademliaID(hash)
	kademlia.nt.AddMessage(target)

	if len(kademlia.nt.storage.RetrieveFile(target)) > 0 { //check in my storage and purge if overdue
		if kademlia.nt.storage.RetrieveTimeSinceStore(target) < time.Hour*24 {
			fmt.Println("found target locally: ", kademlia.nt.storage.RetrieveFile(target))
			return kademlia.nt.storage.RetrieveFile(target)
		} else {
			kademlia.nt.storage.DeleteFile(target)

		}
	}

	result := kademlia.nt.rt.FindClosestContacts(target, count)
	//fmt.Println(len(contacts))

	for j := 0; j < alpha; j++ {
		if j >= len(result) {
			//fmt.Println("BREAK", j)
			break
		}
		go kademlia.nt.SendFindDataMessage(hash, &result[j])
		contacted = append(contacted, []Contact{result[j]}...)
	}
	kademlia.start = time.Now()
	t := time.Now()
	//same := 0

	for {
		t = time.Now()
		if len(kademlia.nt.GetData()) > 0 {
			fmt.Println("\n we got the data: ", kademlia.nt.GetData())
			return kademlia.nt.GetData()
		}
		if t.Sub(kademlia.start).Nanoseconds() > 1000000000 {
			fmt.Println("\nwe got the timeout")
			//fmt.Println("\nhere is the routing table")
			//kademlia.nt.mtx.Lock()
			kademlia.nt.rt.PrintRoutingTable()
			//kademlia.nt.mtx.Unlock()
			return kademlia.nt.GetData()
		}
		/////////////////
		if len(kademlia.GetNetwork().GetResponse()) > 0 {
			//fmt.Println(kademlia.GetNetwork().GetResponse())
			kademlia.start = time.Now()
			tempAlpha := alpha
			temp := kademlia.nt.GetResponse()[0]
			result = kademlia.checkContacts(result, temp)
			if len(result) < alpha {
				tempAlpha = len(result)
			}
			for i := 0; i < tempAlpha && i < len(result) && tempAlpha < count; i++ {
				if existsIn(result[i], contacted) || result[i].ID.Equals(kademlia.nt.rt.me.ID) {
					if tempAlpha < len(result) {
						tempAlpha++
					}

				} else {
					go kademlia.nt.SendFindDataMessage(hash, &result[i])
					contacted = append(contacted, []Contact{result[i]}...)
				}
				//fmt.Println("-", i, "-")
			}
			kademlia.nt.RemoveFirstResponse()
		}
	}
}

func (kademlia *Kademlia) checkData() bool {
	if len(kademlia.nt.GetData()) > 0 {
		return true
	}
	return false
}

func (kademlia *Kademlia) Store(data string) {
	fmt.Println("IM totally gonna store: ", data)
	hashdata := []byte(data)
	key := KademliaID(sha1.Sum(hashdata))
	contacts := kademlia.LookupContact(&key)
	fmt.Println("AFTER LOOKUP IN STORE")
	for j := range contacts {
		go kademlia.nt.SendStoreMessage(&contacts[j], &key, data)
	}

	timer := time.NewTimer(time.Second * 10)
	<-timer.C
	fmt.Println("Timer expired")
	go kademlia.Store(data)
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
