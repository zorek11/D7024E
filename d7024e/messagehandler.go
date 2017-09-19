package d7024e

import (
  "D7024E-Kademlia/protobuf"
  "github.com/golang/protobuf/proto"
  "fmt"
)

type MessageHandler struct {


}

func buildMessage(lable string) *protobuf.KademliaMessage {
  message := &protobuf.KademliaMessage {
    Label: proto.String(lable),
  }
  return message

}

func handleMessage(data []byte) {
  message := &protobuf.Test{}
  err := proto.Unmarshal(data, message)
  if err != nil {
    //log.Fatal("unmarshaling error: ", err)
  }
  switch *message.Label {
  case "ping":
    fmt.Println(message)
  default:
    fmt.Println("panic")
  }
}
