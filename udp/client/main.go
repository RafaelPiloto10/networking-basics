package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("PLease provide a host:port string")
		return
	}

	ADDR := args[1]

	s, err := net.ResolveUDPAddr("udp4", ADDR)
	client, err := net.DialUDP("udp4", nil, s)

	if err != nil {
		fmt.Printf("Got error trying to dial UDP addr %s; %v\n", ADDR, err)
		return
	}

	fmt.Printf("Connected to UDP Server %s\n", client.RemoteAddr().String())

	defer client.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err := client.Write(data)
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Closing UDP client session!")
			return
		}

		if err != nil {
			fmt.Printf("Got error trying to write data to UDP server at %s; %v\n", client.RemoteAddr().String(), err)
			return
		}

		buffer := make([]byte, 1024)

		n, _, err := client.ReadFromUDP(buffer)

		if err != nil {
			fmt.Printf("Got error trying to read from server; err %v\n", err)
			return
		}

		fmt.Printf("-> %s\n", buffer[0:n])
	}
}
