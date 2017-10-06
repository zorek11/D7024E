package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
)

type MessageHandler struct {
	network *Network
}

func NewMessageHandler(net *Network) *MessageHandler {
	mes := &MessageHandler{}
	mes.network = net
	return mes
}

func buildMessage(input []string) *protobuf.KademliaMessage {
	if input[0] == "ping" || input[0] == "pong" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
		}
		return message
	}
	if input[0] == "LookupContactResponse" {
		message := &protobuf.KademliaMessage{
			Label:      proto.String(input[0]),
			Senderid:   proto.String(input[1]),
			SenderAddr: proto.String(input[2]),
			Lookupcontact: &protobuf.KademliaMessage_LookupContact{
				Id:      proto.String(input[3]),
				Address: proto.String(input[4]),
			},
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

func parseContacts(input []string) *protobuf.KademliaMessage {
	message := &protobuf.KademliaMessage{
		Label:      proto.String("LookupContactResponse"),
		Senderid:   proto.String(input[0]),
		SenderAddr: proto.String(input[1]),
	}
	return message
}

func (this *MessageHandler) handleMessage(channel chan []byte, me Contact, network *Network) {
	data := <-channel
	message := &protobuf.KademliaMessage{}
	err := proto.Unmarshal(data, message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\n ========================= \nI am: ", me)
	switch *message.Label {
	case "ping":
		fmt.Println("\n", message)
		response := buildMessage([]string{"pong", me.ID.String(), me.Address})
		send(message.GetSenderAddr(), response)
	case "pong":
		fmt.Print("\n", message)
		network.pingResp = true

	case "LookupContact":
		fmt.Print("\n", message, "\n\n")
		id := NewKademliaID(*message.Lookupcontact.Id)
		temp := network.rt.FindClosestContacts(id, 20) //no recursion

		//==================================
		r := ""
		for i := 0; i < len(temp); i++ {
			r = r + temp[i].String() + "\n"
		}
		response := &protobuf.KademliaMessage{
			Label:      proto.String("LookupContactResponse"),
			Senderid:   proto.String(me.ID.String()),
			SenderAddr: proto.String(me.Address),
			Data:       []byte(r),
		}
		send(message.GetSenderAddr(), response)

	case "LookupContactResponse":
		//response := buildContact(message.Lookupcontact)
		s := string(message.Data)
		contactList := unparse(s)
		if len(contactList) > 0 {
			network.AddResponse(contactList)
		}

	case "LookupData":
		fmt.Print("\n", message)

	case "StoreData":
		key := NewKademliaID(*(message.Key))
		value := *(message.Value)
		senderid := *(message.Senderid)
		network.storage.StoreFile(key, value, senderid)
		network.storage.RetrieveFile(key)
		fmt.Println(network.storage.RetrieveFile(key))

	default:
		fmt.Println("PANIC in switch")

	}

}
func send(Address string, message *protobuf.KademliaMessage) {
	data, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("Marshal Error: ", err)
	}

	Conn, err := net.Dial("udp", Address)
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	defer Conn.Close()
	_, err = Conn.Write(data)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
}

func unparse(input string) []Contact {
	var contactList []Contact
	split := strings.Split(input, "\n")
	for i := 0; i < len(split)-1; i++ {
		out := strings.Split(split[i], "\"")
		contactList = append(contactList, NewContact(NewKademliaID(out[1]), out[3]))

	}
	return contactList
}

func buildContact(message *protobuf.KademliaMessage_LookupContact) Contact {
	return NewContact(NewKademliaID(*message.Id), *message.Address)
}
