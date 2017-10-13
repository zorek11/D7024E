package main

import (
	kademlia "D7024E-Kademlia/d7024e"
	"fmt"
	"math"
	"strconv"
	"time"
)

//export GOPATH=$HOME/go
func main() {

	simulateN(100)

	for {
	}
}

func simulateN(n int) {
	//max 100
	fmt.Println("enter simulation")
	contacts := make([]kademlia.Contact, n)
	for i := 1; i < n; i++ {
		iString := strconv.Itoa(i)
		numbers := (int(math.Log10(float64(i)))) + 1

		extra := ""
		for k := 0; k < 3-numbers; k++ {
			extra += "0"
		}
		contacts[i-1] = kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFFFF11111111111111111111111111"+extra+iString), "127.0.0.1:6"+extra+iString)
		contacts[i-1].CalcDistance(kademlia.NewKademliaID("FFFFFFFFFFF11111111111111111111111111005"))
	}
	kademlias := make([]*kademlia.Kademlia, n)
	for l := 0; l < n; l++ {
		kademlias[l] = kademlia.NewKademlia(contacts[l])
		go kademlias[l].GetNetwork().Listen(contacts[l])
	}
	for m := 1; m < n-1; m++ {
		kademlias[m].GetRoutingtable().AddContact(contacts[0])
		go kademlias[m].LookupContact(contacts[m].ID)
		time.Sleep(100 * time.Millisecond)
	}
	var cc kademlia.ContactCandidates
	cc.Append(contacts[:len(contacts)-1])
	cc.Sort()

	first := kademlias[10].LookupContact(kademlia.NewKademliaID("FFFFFFFFFFF11111111111111111111111111005"))

	//second := kademlias[5].LookupContact(contacts[10].ID)

	fmt.Println("first \n", first, "\n")
	fmt.Println("second \n", cc.GetContacts(20))

	time.Sleep(3000 * time.Millisecond)

	NewAPI("127.0.0.1:9999", kademlias[0])
	for {
		time.Sleep(1000 * time.Millisecond)
		//fmt.Println("still alive")
	}

}
