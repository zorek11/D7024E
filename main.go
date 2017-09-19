package main

import (
  kademlia "D7024E-Kademlia/d7024e"
  "sync"
  //"fmt"

  //"log"
  //"github.com/golang/protobuf/proto"
  //"D7024E-Kademlia/protobuf"
)

func main () {
  var mutex = &sync.Mutex{}
  mutex.Lock()
  defer mutex.Unlock()
  contact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8001")
  //contact2 := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8002")

  net := kademlia.Network{}
  go net.SendPingMessage(&contact)
  //go net.SendPingMessage(&contact)
  //go net.SendPingMessage(&contact2)
  kademlia.Listen("127.0.0.1", 8001)

}
