package d7024e

import (
	"fmt"
	"os"
	"strconv"
)

const count = 20
const alpha = 3

type Kademlia struct {
	rt *RoutingTable
}

func NewKademlia(self Contact) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
	kademlia.rt = NewRoutingTable(self)
	//kademlia.network = new(Network)
	//s := strings.Split(self.Address, ":")
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) *Contact {
	contacts := kademlia.rt.FindClosestContacts(target.ID, count)

	for j, contact := range contacts {
		if contact.ID == kademlia.rt.me.ID {
			return &contacts[j]
		}
	}

	networks := make([]*Network, alpha)

	for i := 0; i < alpha; i++ {
		networks[i] = NewNetwork(kademlia.rt.me, kademlia)
		networks[i].AddMessage(target)
		go networks[i].SendFindContactMessage(&contacts[i])
	}
	//go Listen(kademlia.rt.me)
	fails := 0
	stuck := 0
	for k := 0; k <= alpha; k++ {
		if k == alpha {
			k = 0
		}
		if networks[k].response.ID == target.ID {
			return networks[k].response
		}
		if networks[k].response == nil && fails+alpha <= len(contacts) {
			fails++
			networks[k] = NewNetwork(kademlia.rt.me, kademlia)
			networks[k].AddMessage(target)
			go networks[k].SendFindContactMessage(&contacts[alpha+fails])
		} else if networks[k].response == nil && fails+alpha > len(contacts) {
			stuck++
			if stuck >= 3 {
				break
			}
		} else {
			stuck = 0
		}
	}
	return &contacts[0]
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func IntConverter(port string) int {
	i, err := strconv.Atoi(port)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	return i
}
