package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type MessageHandler struct {
	network *Network
	mtx     *sync.Mutex
}

func NewMessageHandler(net *Network) *MessageHandler {
	mes := &MessageHandler{}
	mes.network = net
	mes.mtx = &sync.Mutex{}

	return mes
}

/**
* Messagehandler for a listner. Handles all messages in a switch and takes according actions.
 */
func (this *MessageHandler) handleMessage(channel chan []byte, me *Contact, network *Network) {

	data := <-channel
	message := &protobuf.KademliaMessage{}
	err := proto.Unmarshal(data, message)
	if err != nil {
		fmt.Println(err)
	}
	sender := NewContact(NewKademliaID(message.GetSenderid()), message.GetSenderAddr())
	//this.mtx.Lock()
	fmt.Print("\n\nListner:", me, "\nSender: ", sender, "\nMessage: ", message)
	//this.mtx.Unlock()
	network.UpdateRoutingtable(sender) //update routingtable on all RPCs
	switch *message.Label {
	case "ping":
		response := buildMessage([]string{"pong", me.ID.String(), me.Address})
		send(message.GetSenderAddr(), response)
	case "pong":
		//this.mtx.Lock()
		network.pingResponse = true
		//this.mtx.Unlock()
		//pingIndex := IndexInSlice(message.GetSenderAddr(), network.pingList)
		//network.pingList[pingIndex].Response = true

	case "LookupContact":
		fmt.Println("\nlookup this id: ", *message.Lookupcontact.Id)
		id := NewKademliaID(*message.Lookupcontact.Id)
		//this.mtx.Lock()
		temp := network.rt.FindClosestContacts(id, 20) //no recursion
		//this.mtx.Unlock()
		//fmt.Println("in lookupcontact case - temp: ", temp)
		r := ""
		for i := 0; i < len(temp); i++ {
			r = r + temp[i].String() + "\n"
		}
		//this.mtx.Lock()
		response := buildMessage([]string{"LookupContactResponse", me.ID.String(), me.Address, r})
		//this.mtx.Unlock()
		send(message.GetSenderAddr(), response)

	case "LookupContactResponse":
		s := string(message.Data)
		contactList := unparseContacts(s)
		if len(contactList) > 0 {
			network.AddResponse(contactList)
		}

	case "LookupData":
		//fmt.Println("this is message key", message.Key)
		key := NewKademliaID(*(message.Key))
		storage := network.storage.RetrieveFile(key)
		if len(storage) > 0 { //if data found
			if network.storage.RetrieveTimeSinceStore(key) < time.Hour*24 {
				//this.mtx.Lock()
				response := buildMessage([]string{"LookupDataResponse", me.ID.String(), me.Address, storage})
				//this.mtx.Unlock()
				send(message.GetSenderAddr(), response)
			} else {
				if network.storage.RetrievePin(key) == false {
					network.storage.DeleteFile(key)
				}
			}
		} else { //return K-closest
			//this.mtx.Lock()
			temp := network.rt.FindClosestContacts(key, 20) //no recursion
			//this.mtx.Unlock()
			//fmt.Println("\nthis is temp", temp)
			r := ""
			for i := 0; i < len(temp); i++ {
				r = r + temp[i].String() + "\n"
			}
			//this.mtx.Lock()
			response := buildMessage([]string{"LookupContactResponse", me.ID.String(), me.Address, r})
			//this.mtx.Unlock()
			send(message.GetSenderAddr(), response)
		}
	case "LookupDataResponse":
		s := string(message.Data)
		network.AddData(s)

	case "StoreData":
		key := NewKademliaID(*(message.Key))
		value := *(message.Value)
		senderid := *(message.Senderid)
		network.storage.StoreFile(key, value, senderid)
		//network.storage.RetrieveFile(key)
		//fmt.Println("WOOT WOOOT + " network.storage.RetrieveFile(key))
	case "Pin":
		key := NewKademliaID(*(message.Key))
		network.storage.Pin(key)
		//fmt.Println("AUHHUFWHKFHAEHFAJEFHJFEA")
	case "Unpin":
		key := NewKademliaID(*(message.Key))
		network.storage.Unpin(key)
	default:
		fmt.Println("PANIC in switch")

	}

}

/**
* Takes a string of contacts and parses it to a slice of contacts
 */
func unparseContacts(input string) []Contact {
	var contactList []Contact
	split := strings.Split(input, "\n")
	for i := 0; i < len(split)-1; i++ {
		out := strings.Split(split[i], "\"")
		contactList = append(contactList, NewContact(NewKademliaID(out[1]), out[3]))

	}
	return contactList
}

/*
 */
func parseContacts(input []string) *protobuf.KademliaMessage {
	message := &protobuf.KademliaMessage{
		Label:      proto.String("LookupContactResponse"),
		Senderid:   proto.String(input[0]),
		SenderAddr: proto.String(input[1]),
	}
	return message
}

func buildContact(message *protobuf.KademliaMessage_LookupContact) Contact {
	return NewContact(NewKademliaID(*message.Id), *message.Address)
}

/**
* Builds a protobuf message from a input array
 */
func buildMessage(input []string) *protobuf.KademliaMessage {
	if input[0] == "ping" || input[0] == "pong" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
		}
		return message
	}
	if input[0] == "LookupContact" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Data:       []byte(input[3]),
		}
		return message
	}
	if input[0] == "LookupData" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Data:       []byte(input[3]),
			Key:        proto.String(input[4]),
		}
		return message
	}
	if input[0] == "LookupContactResponse" || input[0] == "LookupDataResponse" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Data:       []byte(input[3]),
		}
		return message
	}
	if input[0] == "StoreData" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Key:        proto.String(input[3]),
			Value:      proto.String(input[4]),
		}
		return message
	}
	if input[0] == "Pin" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Key:        proto.String(input[3]),
		}
		return message
	}
	if input[0] == "Unpin" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Key:        proto.String(input[3]),
		}
		return message
	} else {
		message := &protobuf.KademliaMessage{
			Label:      proto.String("Error"),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
		}
		return message
	}
}
