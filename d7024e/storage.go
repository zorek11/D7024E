package d7024e

import (
	"time"
)

//"fmt"
type Time struct {
	// contains filtered or unexported fields
}

type Storage struct {
	publisherht map[KademliaID]string
	valueht     map[KademliaID]string
	timeht      map[KademliaID]time.Time
	pinht       map[KademliaID]bool //TODO: ADD LATER???
}

func NewStorage() Storage {
	var storage Storage
	storage.publisherht = make(map[KademliaID]string)
	storage.valueht = make(map[KademliaID]string)
	storage.timeht = make(map[KademliaID]time.Time)
	storage.pinht = make(map[KademliaID]bool)
	return storage
}

func (storage *Storage) StoreFile(key *KademliaID, value string, publisher string) {
	start := time.Now()
	if len(storage.valueht) != 0 {
		if storage.publisherht[*key] == publisher {
			storage.publisherht[*key] = publisher
			storage.valueht[*key] = value
			storage.timeht[*key] = start
			storage.pinht[*key] = false
		}
	} else {
		start := time.Now()
		storage.publisherht[*key] = publisher
		storage.valueht[*key] = value
		storage.timeht[*key] = start
		storage.pinht[*key] = false
	}
}

func (storage *Storage) DeleteFile(key *KademliaID) {
	delete(storage.valueht, *key)
	delete(storage.publisherht, *key)
	delete(storage.timeht, *key)
	//delete(storage.pinht, *key)
}

func (storage *Storage) RetrieveFile(key *KademliaID) string {
	return storage.valueht[*key]
}

func (storage *Storage) RetrievePublisher(key *KademliaID) string {
	return storage.publisherht[*key]
}

func (storage *Storage) RetrieveTimeSinceStore(key *KademliaID) time.Duration {
	start := time.Now()
	return start.Sub(storage.timeht[*key])
}

func (storage *Storage) Pin(key *KademliaID) {
	storage.pinht[*key] = true
}

func (storage *Storage) UnPin(key *KademliaID) {
	storage.pinht[*key] = false
}
