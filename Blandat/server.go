package main

import (
    "fmt"
    "net"
    "strconv"
)
func main() {
  Listen("127.0.0.1", 8001)
}

func Listen(ip string, port int) {
	// TODO
  fmt.Println("ListnerUp")
	ServerAddr,err1 := net.ResolveUDPAddr("udp",ip+":"+strconv.Itoa(port))
	ServerConn, err2 := net.ListenUDP("udp", ServerAddr)
  if (err1  != nil) || (err2 != nil){
      fmt.Println("Error: ", err1, err2)
  }

	//read connection
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
			n,addr,err := ServerConn.ReadFromUDP(buf)
			fmt.Println("Received ",string(buf[0:n]), " from ",addr)
			if err != nil {
					fmt.Println("Error: ",err)
			}
	}
}
