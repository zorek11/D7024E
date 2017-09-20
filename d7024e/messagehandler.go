package d7024e

import (
	"D7024E-Kademlia/protobuf"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

type MessageHandler struct {
}

func buildMessage(lable string, id string, addr string) *protobuf.KademliaMessage {
	message := &protobuf.KademliaMessage{
		Label:      proto.String(lable),
		Senderid:   proto.String(id),
		SenderAddr: proto.String(addr),
	}
	return message

}

func handleMessage(channel chan []byte, me Contact) {
	data := <-channel
	message := &protobuf.KademliaMessage{}
	err := proto.Unmarshal(data, message)
	if err != nil {
		//log.Fatal("unmarshaling error: ", err)
	}
	switch *message.Label {
	case "ping":
		fmt.Println("\n", message)
		handlePing(message, me)

	case "pingResponse":
		fmt.Print("\n", message)

	case "LookupContact":
		fmt.Print("\n", message)

	case "LookupData":
		fmt.Print("\n", message)

	case "StoreData":
		fmt.Print("\n", message)

	default:
		fmt.Println("PANIC")

	}
}

func handlePing(message *protobuf.KademliaMessage, me Contact) {
	response := buildMessage("pingResponse", me.ID.String(), me.Address)
	data, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("Marshal Error: ", err)
	}
	Conn, err := net.Dial("udp", message.GetSenderAddr())
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	defer Conn.Close()

	_, err = Conn.Write(data)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
}

func handlePingResponse(message *protobuf.KademliaMessage, me Contact) {

}
