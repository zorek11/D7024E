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

func handleMessage(data []byte, me Contact) {
	message := &protobuf.KademliaMessage{}
	err := proto.Unmarshal(data, message)
	if err != nil {
		//log.Fatal("unmarshaling error: ", err)
	}
	switch *message.Label {
	case "ping":
		fmt.Println("\n", message)
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
	case "pingResponse":
		fmt.Print("\n", message)
		//TODO:
	default:

	}
}
