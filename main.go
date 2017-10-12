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

	/*
		contact2 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111192"),
			"127.0.0.1:1422")
		contact3 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111193"),
			"127.0.0.1:1423")
		contact4 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111194"),
			"127.0.0.1:1424")
		contact5 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111195"),
			"127.0.0.1:1425")
		contact6 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111196"),
			"127.0.0.1:1426")
		contact7 := kademlia.NewContact(kademlia.NewKademliaID("1FFFFFFF11111111111111111111111111111197"),
			"127.0.0.1:1427")

		me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
			"127.0.0.1:1089")

		kad1 := kademlia.NewKademlia(me)

		kad1.AddRoutingtable(contact2)

		kad2 := kademlia.NewKademlia(contact2)

		kad2.AddRoutingtable(contact4)

		kad3 := kademlia.NewKademlia(contact3)
		kad3.AddRoutingtable(contact4)
		kad3.AddRoutingtable(contact2)

		kad4 := kademlia.NewKademlia(contact4)
		kad4.AddRoutingtable(contact3)
		kad4.AddRoutingtable(contact5)

		kad5 := kademlia.NewKademlia(contact5)
		kad5.AddRoutingtable(contact4)
		kad5.AddRoutingtable(contact6)

		kad6 := kademlia.NewKademlia(contact6)
		kad6.AddRoutingtable(contact5)
		kad6.AddRoutingtable(contact7)

		kad7 := kademlia.NewKademlia(contact7)
		kad7.AddRoutingtable(contact6)

		net1 := kad1.GetNetwork()
		net2 := kad2.GetNetwork()
		net3 := kad3.GetNetwork()
		net4 := kad4.GetNetwork()
		net5 := kad5.GetNetwork()
		net6 := kad6.GetNetwork()
		net7 := kad7.GetNetwork()

		go net1.Listen(me)
		go net2.Listen(contact2)
		go net3.Listen(contact3)
		go net4.Listen(contact4)
		go net5.Listen(contact5)
		go net6.Listen(contact6)
		go net7.Listen(contact7)
	*/
	//net1.AddMessage(contact2.ID)
	//go kad1.LookupContact(contact7.ID)
	/*
		str := "aids in the face"
		hash := kademlia.KademliaID(sha1.Sum([]byte(str)))
		fmt.Println("Det här är ursprungshash i main: " + hash.String())
		kad7.GetNetwork().GetStorage().StoreFile(&hash, str, me.ID.String())
		//fmt.Println("RetrieveFile via main: " + kad1.GetNetwork().GetStorage().RetrieveFile(&hash))
		fmt.Println(kad1.LookupData("2718e9414b3cd72cbe02601fb53842576c5b6435"))
	*/
	//fmt.Print(kad1.Store("aids"))
	//go net1.SendFindContactMessage(&contact2)

	/*for {
		if kad1.GetFound() != false {
			fmt.Println("END")
			break
		}
	}*/
	/*
		for i := 0; i < 3; i++ {
			go net1.SendPingMessage(&contact2)
			go net1.SendPingMessage(&contact7)
		}
	*/
	/*
		net2.GetRoutingTable().PrintRoutingTable()
		net2.UpdateRoutingtable(contact3)
		net2.GetRoutingTable().PrintRoutingTable()

		fmt.Println("GOLANG är fan sämst")*/
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

		/*
			for j := n - i; j > 0; j = j / 10 {
				if j%10 == 0 {
					//numbers--
				} else {
					numbers++
				}
			}
		*/
		fmt.Println(numbers)
		extra := ""
		for k := 0; k < 3-numbers; k++ {
			extra += "0"
		}
		fmt.Println("127.0.0.1:9" + extra + iString)
		contacts[i-1] = kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFFFF11111111111111111111111111"+extra+iString), "127.0.0.1:7"+extra+iString)

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
		/*
			if m == 0 {
				fmt.Println(contacts[m+1])
				kademlias[m].GetRoutingtable().AddContact(contacts[m+1])
			} else if m == n-2 {
				kademlias[m].GetRoutingtable().AddContact(contacts[m-1])
			} else if m > 5 && m < n-7 {
				//kademlias[m].GetRoutingtable().AddContact(contacts[m+4])
				kademlias[m].GetRoutingtable().AddContact(contacts[m+3])
				kademlias[m].GetRoutingtable().AddContact(contacts[m+2])
				kademlias[m].GetRoutingtable().AddContact(contacts[m+1])
				kademlias[m].GetRoutingtable().AddContact(contacts[m-1])
				kademlias[m].GetRoutingtable().AddContact(contacts[m-2])
				kademlias[m].GetRoutingtable().AddContact(contacts[m-3])
				kademlias[m].GetRoutingtable().AddContact(contacts[m-4])
			} else {
				kademlias[m].GetRoutingtable().AddContact(contacts[m+1])
			}
		*/
	}
	time.Sleep(2000 * time.Millisecond)
	/*
		str := "aids in the face"
		hash := kademlia.KademliaID(sha1.Sum([]byte(str)))
		fmt.Println("Det här är ursprungshash i main: " + hash.String())
		go kademlias[n/10].Store(str)
		time.Sleep(1000 * time.Millisecond)
	kademlias[n/2].LookupData(hash.String())*/

	//kademlias[0].LookupContact(contacts[n-2].ID)

}

func start100() {
	for i := 0; i < 50; i++ {
		iString := strconv.Itoa(i)
		var contact kademlia.Contact
		if i < 10 {
			contact = kademlia.NewContact(kademlia.NewKademliaID(
				"FFFFFFFFFFF1111111111111111111111111190"+iString), "127.0.0.1:990"+iString)
		} else {
			contact = kademlia.NewContact(kademlia.NewKademliaID(
				"FFFFFFFFFFF111111111111111111111111119"+iString), "127.0.0.1:99"+iString)
		}
		kad := kademlia.NewKademlia(contact)
		net := kad.GetNetwork()
		go net.Listen(contact)
	}
}
