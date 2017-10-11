package d7024e

import (
	"sync"
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
	mtx         *sync.Mutex
}

func NewStorage() Storage {
	var storage Storage
	storage.publisherht = make(map[KademliaID]string)
	storage.valueht = make(map[KademliaID]string)
	storage.timeht = make(map[KademliaID]time.Time)
	storage.pinht = make(map[KademliaID]bool)
	storage.mtx = &sync.Mutex{}
	return storage
}

func (storage *Storage) StoreFile(key *KademliaID, value string, publisher string) {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
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
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	delete(storage.valueht, *key)
	delete(storage.publisherht, *key)
	delete(storage.timeht, *key)
	//delete(storage.pinht, *key)
}

func (storage *Storage) RetrieveFile(key *KademliaID) string {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	return storage.valueht[*key]
}

func (storage *Storage) RetrievePublisher(key *KademliaID) string {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	return storage.publisherht[*key]
}

func (storage *Storage) RetrieveTimeSinceStore(key *KademliaID) time.Duration {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	start := time.Now()
	return start.Sub(storage.timeht[*key])
}

func (storage *Storage) Pin(key *KademliaID) {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	storage.pinht[*key] = true
}

func (storage *Storage) Unpin(key *KademliaID) {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	storage.pinht[*key] = false
}

func (storage *Storage) RetrievePin(key *KademliaID) bool {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	return storage.pinht[*key]
}
