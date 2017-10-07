package d7024e

import (
	"fmt"
)

type Storage struct {
	publisherht map[*KademliaID]string
	valueht     map[*KademliaID]string
	//pin   boolean //TODO: ADD LATER???
}

func NewStorage() Storage {
	var storage Storage
	storage.publisherht = make(map[*KademliaID]string)
	storage.valueht = make(map[*KademliaID]string)
	return storage
}

func (storage *Storage) StoreFile(key *KademliaID, value string, publisher string) {
	//fmt.Println(value)
	storage.publisherht[key] = publisher
	storage.valueht[key] = value

}
func (storage *Storage) RetrieveFile(key *KademliaID) string {
	fmt.Println("here is the retrieved file string: ", storage.valueht[key])
	return storage.valueht[key]
}

func (storage *Storage) RetrievePublisher(key *KademliaID, value string, publisher string) string {
	return storage.publisherht[key]
}
