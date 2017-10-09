package d7024e

/*
import (
	"fmt"
	"net"
)

func CliListen(me Contact) {
	Addr, err1 := net.ResolveUDPAddr("udp", me.Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error Listen: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	channel := make(chan []byte)
	buf := make([]byte, 1024)
	for {

		switch channel {
		case "dick snack":
			fmt.Println("\n", channel)
			n, _, err := Conn.ReadFromUDP(buf)
			channel <- buf[0:n]
			//fmt.Println("Connection recived: ", string(buf[0:n]), " \nfrom ", addr)

			if err != nil {
				fmt.Println("Read Error: ", err)
			}
		}
	}
}
*/
