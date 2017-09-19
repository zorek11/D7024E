package d7024e

import (
	"flag"
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
	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
	contacts := kademlia.rt.FindClosestContacts(target.ID, count)
	if target.ID != contacts[0].ID {
		for i := 0; i < alpha; i++ {
			//s := strings.Split(contacts[i].Address, ":")
			go Listen(*target)
		}
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
	flag.Parse()
	// string to int
	i, err := strconv.Atoi(port)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	return i
}
