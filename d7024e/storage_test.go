package d7024e

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	storage := NewStorage()
	key := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	publisher := "127.0.0.1"

	storage.StoreFile(key, "testString", publisher)
	retrivestored := storage.RetrieveFile(key)

	storage.DeleteFile(key)
	retrivedeleted := storage.RetrieveFile(key)
	if len(retrivedeleted) == 0 {
		retrivedeleted = "Deleted."
	}
	storage.StoreFile(key, "testString", publisher)
	storage.Pin(key)
	publisher = storage.RetrievePublisher(key)
	time := storage.RetrieveTimeSinceStore(key)
	pin := storage.RetrievePin(key)

	storage.DeleteFile(key)
	pinnedDelete := storage.RetrieveFile(key)
	storage.Unpin(key)
	storage.DeleteFile(key)
	unpinnneddelete := storage.RetrieveFile(key)

	fmt.Printf("\nSTORE\nStored: %v Delete: %v Publisher: %v Time: %v Pinned: %v DeleteofPinned: %v DeleteofUnpinned: %v \n\n",
		retrivestored, retrivedeleted, publisher, time, pin, pinnedDelete, unpinnneddelete)
}
