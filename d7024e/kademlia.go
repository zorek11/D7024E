package d7024e

import (
	"fmt"
	"os"
	"strconv"
)

const count = 20
const alpha = 3

type Kademlia struct {
	nt    Network
	items []string
}

func (kademlia *Kademlia) AddRoutingtable(c Contact) {
	kademlia.nt.rt.AddContact(c)
}

func (kademlia *Kademlia) GetRoutingtable() *RoutingTable {
	return kademlia.nt.rt
}

func (kademlia *Kademlia) GetNetwork() *Network {
	return &kademlia.nt
}

func NewKademlia(self Contact) (kademlia *Kademlia) {
	kademlia = new(Kademlia)
	kademlia.nt = NewNetwork(self, NewRoutingTable(self))
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {

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
	kademlia.nt.AddMessage(target)

	go kademlia.nt.SendFindContactMessage(&contacts[0])
	//}
	for {
		if kademlia.GetNetwork().GetResponse() != nil {
			//fmt.Println(tempnetwork.GetResponse()[0].String())
			break
		}
	}
	/*
		for {
			for k := 0; k < thisalpha; k++ {

				if networks[k].GetResponse() != nil {
					fmt.Println(networks[k].GetResponse()[0].String())
					break
				}

				//fmt.Println(networks[k].GetResponse())
				//fmt.Println(networks[k].response == nil)
				if networks[k].GetResponse() != nil {
					if networks[k].response[0].ID == target.ID {
						fmt.Println("Target found: " + target.String())
						fmt.Println("With address: " + networks[k].response[0].String())
						return
						//return networks[k].response[0]
					} else {
						tempAlpha := alpha % (len(networks[k].response) + 1)
						tempNetworks := getNewNetworks(kademlia, networks[k].GetResponse(), tempAlpha, target)
						sendFindContactForAll(tempNetworks, target)
						networks = remove(networks, k)
						networks = combineNetworks(networks, tempNetworks)
						thisalpha = len(networks)

					}
				}

			}
		}
	*/
	//return &contacts[0]

}

func (kademlia *Kademlia) LookupData(hash string) {
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

/*
func getNewNetworks(kademlia *Kademlia, contacts []Contact, alpha int, target *Contact) []*Network {
	networks := make([]*Network, alpha)

	for i := 0; i < alpha; i++ {
		networks[i] = NewNetwork(kademlia.rt.me, kademlia)
		networks[i].AddMessage(target)
	}
	return networks
}
*/
func sendFindContactForAll(networks []*Network, target *Contact) {
	for i := 0; i < len(networks); i++ {
		go networks[i].SendFindContactMessage(target)
	}
}

func remove(networks []*Network, s int) []*Network {
	return append(networks[:s], networks[s+1:]...)
}

func combineNetworks(first []*Network, second []*Network) []*Network {
	networks := make([]*Network, (len(second) + len(first)))
	for i := 0; i < len(first); i++ {
		networks[i] = first[i]
	}
	for j := len(first) - 1; j < len(second); j++ {
		networks[j] = first[j-len(first)+1]
	}
	return networks

}
