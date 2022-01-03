package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port")
		return
	}

	addr := arguments[1]

	client, err := net.Dial("tcp", addr)

	if err != nil {
		fmt.Printf("Got error trying to dial host %s; %v\n", addr, err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(client, text+"\n")

		message, _ := bufio.NewReader(client).ReadString('\n')
		fmt.Print("->:" + message)

		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP Client exiting...")
			return
		}
	}
}
