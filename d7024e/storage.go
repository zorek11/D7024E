package d7024e

type Storage struct {
  	publisherht  map[KademliaID][]byte
    valueht  map[KademliaID][]byte
    //pin   boolean //TODO: ADD LATER???
  }

  func NewStorage() Storage {
    publisherht = make(map[KademliaID][]byte)
    valueht = make(map[KademliaID][]byte)
  	return storage
  }

  func (storage *Storage) StoreFile(key KademliaID, value string, publisher string){
      storage.publisherht[key] = publisher
      storage.valueht[key] = value
  }
  func (storage *Storage) RetrieveFile(key KademliaID{
      storage.valueht[key]
  }

	func (storage *Storage) RetrievePublisher(key KademliaID, value string, publisher Contact){
			storage.publisherht[key]
	}
