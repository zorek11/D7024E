package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"
	"time"

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
	} else {
		fmt.Println("\n not fatal error")

	}
	response := (*protobuf.KademliaMessage)(nil)
	send := false
	pingpong := false
	switch *message.Label {
	case "ping":
		fmt.Println("\n", message)
		response = buildMessage([]string{"pong", me.ID.String(), me.Address})
		send = true

	case "pong":
		fmt.Print("\n", message)
		response = buildMessage([]string{"pong", me.ID.String(), me.Address})
		pingpong = true

	case "LookupContact":
		fmt.Print("\n", message, "\n\n")
		contact := buildContact(message.Lookupcontact)
		//network.kademlia.LookupContact(&contact)
		temp := network.kademlia.rt.FindClosestContacts(contact.ID, 20)
		response = parseContacts([]string{temp[0].ID.String(), temp[0].Address})
		send = true
		fmt.Println(response.SenderAddr)
	case "LookupContactResponse":
		if message.Lookupcontact != nil {
			fmt.Print("\n", message)
			response := buildContact(message.Lookupcontact)
			network.AddTempResponse(&response)
		} else {
			return
		}
	case "LookupData":
		fmt.Print("\n", message)

	case "StoreData":
		fmt.Print("\n", message)

	default:
		fmt.Println("PANIC")

	}
	if send { //marshal and send message
		data, err = proto.Marshal(response)
		if err != nil {
			fmt.Println("Marshal Error: ", err)
		}

		//LocalAddr, err := net.ResolveUDPAddr("udp", network.kademlia.GetRoutingtable().me.Address)
		ServerAddr, err := net.ResolveUDPAddr("udp", message.GetSenderAddr())
		Conn, err := net.ListenUDP("udp", ServerAddr)
		if err != nil {
			fmt.Println("UDP-Error: ", err)
		}
		defer Conn.Close()
		fmt.Println("\n writing to UDP")
		fmt.Println(data)
		_, err = Conn.WriteToUDP(data, ServerAddr)
		if err != nil {
			fmt.Println("Write Error: ", err)
		}
		send = false
	}
	if message.GetLabel() == "ping" {
		fmt.Println(time.Now())
		time.After(time.Second * 3)
		fmt.Println("Pung")
		if !pingpong {
			fmt.Println("Timeout: ", message.GetSenderAddr())
		}
	}

}

func buildContact(message *protobuf.KademliaMessage_LookupContact) Contact {
	return NewContact(NewKademliaID(*message.ID), *message.Address)
}
