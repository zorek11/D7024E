package d7024e

import (
	"fmt"
	"os"
	"strconv"
)

const count = 20
const alpha = 3

type Kademlia struct {
	rt    *RoutingTable
	items []string
}

func (kademlia *Kademlia) AddRoutingtable(c Contact) {
	kademlia.rt.AddContact(c)
}

func (kademlia *Kademlia) GetRoutingtable() *RoutingTable {
	return kademlia.rt
}

func NewKademlia(self Contact) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
	kademlia.rt = NewRoutingTable(self)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	contacts := kademlia.rt.FindClosestContacts(target.ID, count)
	thisalpha := 0
	if len(contacts) < alpha {
		thisalpha = len(contacts)
	} else {
		thisalpha = alpha
	}

	if contacts[0].ID == target.ID {
		//return &contacts[0]
		fmt.Println("Target found: " + target.String())
		fmt.Println("With address: " + contacts[0].String())
		return
	}

	networks := make([]*Network, thisalpha)

	for i := 0; i < thisalpha; i++ {
		networks[i] = NewNetwork(kademlia.rt.me, kademlia)
		networks[i].AddMessage(target)
		go networks[i].SendFindContactMessage(&contacts[i])
	}

	for k := 0; k <= thisalpha; k++ {
		if k == thisalpha {

			k = 0
		}
		if networks[k].response != nil {
			if networks[k].response[0].ID == target.ID {
				fmt.Println("Target found: " + target.String())
				fmt.Println("With address: " + networks[k].response[0].String())
				return
				//return networks[k].response[0]
			} else {
				tempCon := networks[k].response[0]
				networks[k] = NewNetwork(kademlia.rt.me, kademlia)
				networks[k].AddMessage(target)
				go networks[k].SendFindContactMessage(&tempCon)
			}
		}

	}
	//return &contacts[0]

}

func (kademlia *Kademlia) LookupData(hash string) {
	target := NewContact(NewKademliaID(hash), "")
	contacts := kademlia.rt.FindClosestContacts(target.ID, count)
	thisalpha := 0
	if len(contacts) < alpha {
		thisalpha = len(contacts)
	} else {
		thisalpha = alpha
	}

	if contacts[0].ID == target.ID {
		//return &contacts[0]
		fmt.Println("Target found: " + target.String())
		fmt.Println("With address: " + contacts[0].String())
		return
	}

	networks := make([]*Network, thisalpha)

	for i := 0; i < thisalpha; i++ {
		networks[i] = NewNetwork(kademlia.rt.me, kademlia)
		networks[i].AddMessage(&target)
		go networks[i].SendFindContactMessage(&contacts[i])
	}

	for k := 0; k <= thisalpha; k++ {
		if k == thisalpha {
			k = 0
		}
		if networks[k].response != nil {
			if networks[k].response[0].ID == target.ID {
				fmt.Println("Target found: " + target.String())
				fmt.Println("With address: " + networks[k].response[0].String())
				return
				//return networks[k].response[0]
			} else {
				tempCon := networks[k].response[0]
				networks[k] = NewNetwork(kademlia.rt.me, kademlia)
				networks[k].AddMessage(&target)
				go networks[k].SendFindContactMessage(&tempCon)
			}
		}
	}
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
