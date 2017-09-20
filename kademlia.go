package d7024e

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const count = 20
const alpha = 3

type Kademlia struct {
	rt      *RoutingTable
	network *Network
}

func NewKademlia(self Contact) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
	kademlia.rt = NewRoutingTable(self)
	kademlia.network = new(Network)
	//s := strings.Split(self.Address, ":")
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
	contacts := kademlia.rt.FindClosestContacts(target.ID, count)

	for i := 0; i < len(contacts); i++ {
		if target.ID == contacts[i].ID {
			return contacts[i]
		}
	}

	for i := 0; i < alpha; i++ {
		s := strings.Split(contacts[i].Address, ":")
		fmt.Println(s)
		//go Listen(s[0], IntConverter(s[1]))
	}
	return contacts[0]
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
