package main

import (
	kademlia "D7024E-Kademlia/d7024e"
	"crypto/sha1"
	"fmt"
	"net"
	"strings"
)

/*
 */
type API struct {
	kademlia *kademlia.Kademlia
}

/*
 */
func NewAPI(address string, kademlia *kademlia.Kademlia) {
	api := &API{}
	api.kademlia = kademlia
	go api.Listener(address)
}

/*
 */
func (api *API) Listener(Address string) {
	Addr, err1 := net.ResolveUDPAddr("udp", Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error Listen: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	buf := make([]byte, 4096)
	for {
		n, addr, err := Conn.ReadFromUDP(buf)
		go handleTraffic(buf[:n], api, addr.String())
		//fmt.Println("Received ", string(buf[:n]), " from ", addr)

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
	}
}

/*
 */
func handleTraffic(traffic []byte, api *API, sender string) {
	out := strings.Split(string(traffic), ",")
	switch out[0] {
	case "store":
		r := kademlia.KademliaID(sha1.Sum([]byte(out[1])))
		api.kademlia.Store(out[1])
		UDPsend(r.String(), sender)
	case "cat":
		r := api.kademlia.LookupData(out[1])
		if r == "" {
			r = "Data not found, sorry!"
		}
		UDPsend(r, sender)
	case "pin":
		api.kademlia.Pin(out[1])
	case "unpin":
		api.kademlia.Unpin(out[1])
	case "default":
	}
}

/*
 */
func UDPsend(data string, address string) {
	Conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	buf := []byte(data)
	defer Conn.Close()
	_, err = Conn.Write(buf)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
}
