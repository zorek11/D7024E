package d7024e

import (
	"container/list"
	"sync"
)

type bucket struct {
	list *list.List
	mtx  *sync.Mutex
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.mtx = &sync.Mutex{}
	return bucket
}

//this.mtx.Lock()
//defer this.mtx.Unlock()

func (bucket *bucket) AddContact(contact Contact) {
	bucket.mtx.Lock()
	defer bucket.mtx.Unlock()
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

func (bucket *bucket) RemoveContact(contact Contact) {
	bucket.mtx.Lock()
	defer bucket.mtx.Unlock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			bucket.list.Remove(e)
			break
		}
	}
}

func (bucket *bucket) ContactinBucket(contact Contact) bool {
	bucket.mtx.Lock()
	defer bucket.mtx.Unlock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			return true
		}
	}
	return false
}

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	bucket.mtx.Lock()
	defer bucket.mtx.Unlock()
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
