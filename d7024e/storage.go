package d7024e

//"fmt"

type Storage struct {
	publisherht map[KademliaID]string
	valueht     map[KademliaID]string
	timeht      map[KademliaID]string
	//pin   boolean //TODO: ADD LATER???
}

func NewStorage() Storage {
	var storage Storage
	storage.publisherht = make(map[KademliaID]string)
	storage.valueht = make(map[KademliaID]string)
	storage.timeht = make(map[KademliaID]string)
	return storage
}

func (storage *Storage) StoreFile(key *KademliaID, value string, publisher string) {
	//fmt.Println(value)
	storage.publisherht[*key] = publisher
	storage.valueht[*key] = value
	//storage.timeht[*key] = t
}
func (storage *Storage) DeleteFile(key *KademliaID) {
	delete(storage.valueht, *key)
	delete(storage.publisherht, *key)
}

func (storage *Storage) RetrieveFile(key *KademliaID) string {
	return storage.valueht[*key]
}

func (storage *Storage) RetrievePublisher(key *KademliaID, value string, publisher string) string {
	return storage.publisherht[*key]
}
