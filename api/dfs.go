package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	address := "127.0.0.1:9999"
	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("Missing subcommand. The avalible commands are: \n-store, \n-cat, \n-pin, \n-unpin")
		os.Exit(1)
	}
	if len(os.Args) < 3 {
		fmt.Println("Missing argument. Every subcommand requires one argument.")
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		fmt.Println("To many arguments. Every subcommand requires one argument")
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "store":
		fmt.Printf("This is: %s with argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)
	case "cat":
		fmt.Printf("This is: %s with argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)

	case "pin":
		fmt.Printf("This is: %s with argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)

	case "unpin":
		fmt.Printf("This is: %s with argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)

	default:
		fmt.Printf("Unavalibe subcommand: %s \n", os.Args[1])
		os.Exit(1)
	}
}

/*
* Sends data to the address by udp
 */
func send(command string, args string, address string) {
	ServerAddr, _ := net.ResolveUDPAddr("udp", address)
	LocalAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9988")
	Conn, _ := net.DialUDP("udp", LocalAddr, ServerAddr)
	/*if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	*/buf := []byte(command + "," + args)
	_, err := Conn.Write(buf)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
	Conn.Close()
	if command == "store" || command == "cat" {
		Listener(LocalAddr.String())
	}
}

/*
 */
func Listener(Address string) {
	fmt.Println("GO LISTEN")
	Addr, err1 := net.ResolveUDPAddr("udp", Address)
	Conn, err2 := net.ListenUDP("udp", Addr)

	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error Listen: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	buf := make([]byte, 4096)
	for {
		n, addr, err := Conn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
		break
	}
}

//go build dfs.go
//sudo cp dfs /usr/local/bin
