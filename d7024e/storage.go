package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

type Storage struct {
  	publisherht  map[KademliaID][]byte
    valueht  map[KademliaID][]byte
    //pin   boolean //TODO: ADD LATER???
  }


  func NewStorage(publisher Contact, value string,) Storage {
    publisherht = make(map[KademliaID][]byte)
    valueht = make(map[KademliaID][]byte)
  	return Storage{publisher, ht, value}
  }

  func (storage *Storage) StoreFile(key KademliaID, value string, publisher Contact){
      storage.publisherht[key] = publisher
      storage.valueht[key] = value
  }
  func (storage *Storage) RetrieveFile(key KademliaID, value string, publisher Contact){
      storage.publisherht[key]
      storage.valueht[key]
  }
