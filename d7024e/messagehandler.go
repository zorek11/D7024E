package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"

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
				ID:      proto.String(input[3]),
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
	response := (*protobuf.KademliaMessage)(nil)
	send := false
	fmt.Println("I am: ", me)
	switch *message.Label {
	case "ping":
		fmt.Println("\n", message)
		response = buildMessage([]string{"pong", me.ID.String(), me.Address})
		send = true

	case "pong":
		fmt.Print("\n", message)
		response = buildMessage([]string{"pong", me.ID.String(), me.Address})

	case "LookupContact":
		fmt.Print("\n", message, "\n\n")
		contact := buildContact(message.Lookupcontact)
		temp := network.kademlia.rt.FindClosestContacts(contact.ID, 20)
		response = buildMessage([]string{"LookupContactResponse", me.ID.String(), me.Address, temp[0].ID.String(), temp[0].Address})
		send = true
	case "LookupContactResponse":

		fmt.Print("\n", message)
		response := buildContact(message.Lookupcontact)
		network.AddTempResponse(&response)
	case "LookupData":
		fmt.Print("\n", message)

	case "StoreData":
		fmt.Print("\n", message)

	default:
		fmt.Println("PANIC in switch")

	}
	if send { //marshal and send message
		data, err = proto.Marshal(response)
		if err != nil {
			fmt.Println("Marshal Error: ", err)
		}

		Conn, err := net.Dial("udp", network.kademlia.GetRoutingtable().me.Address)
		if err != nil {
			fmt.Println("UDP-Error: ", err)
		}
		defer Conn.Close()
		_, err = Conn.Write(data)
		if err != nil {
			fmt.Println("Write Error: ", err)
		}
		send = false
	}

}

func buildContact(message *protobuf.KademliaMessage_LookupContact) Contact {
	return NewContact(NewKademliaID(*message.ID), *message.Address)
}
