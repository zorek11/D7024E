package d7024e

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestsimulateN(t *testing.T) {
	n := 100
	//max 100
	fmt.Println("enter simulation")
	contacts := make([]Contact, n)
	for i := 1; i < n; i++ {
		iString := strconv.Itoa(i)
		numbers := (int(math.Log10(float64(i)))) + 1

		extra := ""
		for k := 0; k < 3-numbers; k++ {
			extra += "0"
		}
		contacts[i-1] = NewContact(NewKademliaID("FFFFFFFFFFF11111111111111111111111111"+extra+iString), "127.0.0.1:6"+extra+iString)

	}
	kademlias := make([]*Kademlia, n)
	for l := 0; l < n; l++ {
		kademlias[l] = NewKademlia(contacts[l])
		go kademlias[l].GetNetwork().Listen(contacts[l])
	}
	for m := 1; m < n-1; m++ {
		kademlias[m].GetRoutingtable().AddContact(contacts[0])
		go kademlias[m].LookupContact(contacts[m].ID)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5000 * time.Millisecond)

	//NewAPI("127.0.0.1:9999", kademlias[0])
}

func TestRepublish(t *testing.T) {
	simulateN(100)
	killNode("127.0.0.1:6001")

}

func killNode(address string) {
	Conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	buf := []byte("KILL")
	defer Conn.Close()
	_, err = Conn.Write(buf)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
}

func simulateN(n int) {
	//max 100
	fmt.Println("enter simulation")
	contacts := make([]Contact, n)
	for i := 1; i < n; i++ {
		iString := strconv.Itoa(i)
		numbers := (int(math.Log10(float64(i)))) + 1

		extra := ""
		for k := 0; k < 3-numbers; k++ {
			extra += "0"
		}
		contacts[i-1] = NewContact(NewKademliaID("FFFFFFFFFFF11111111111111111111111111"+extra+iString), "127.0.0.1:6"+extra+iString)

	}
	kademlias := make([]*Kademlia, n)
	for l := 0; l < n; l++ {
		kademlias[l] = NewKademlia(contacts[l])
		go kademlias[l].GetNetwork().Listen(contacts[l])
	}
	for m := 1; m < n-1; m++ {
		kademlias[m].GetRoutingtable().AddContact(contacts[0])
		go kademlias[m].LookupContact(contacts[m].ID)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(3000 * time.Millisecond)

	for {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("still alive")
	}

}
