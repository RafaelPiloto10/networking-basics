package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Printf("Please provide a port number!\n")
		return
	}

	PORT := ":" + arguments[1]
	conn, err := net.Listen("tcp4", PORT)

	if err != nil {
		fmt.Printf("Got error trying to listen on tcp4 %v; %v\n", PORT, err)
		return
	}

	defer conn.Close()

	for {
		client, err := conn.Accept()
		if err != nil {
			fmt.Printf("Got error trying to accept a connection; %v\n", err)
			return
		}

		go handleConnection(client)
	}

}

func handleConnection(client net.Conn) {
	fmt.Printf("Serving %s\n", client.RemoteAddr().String())
	rand.Seed(time.Now().UnixNano())

	for {
		netData, err := bufio.NewReader(client).ReadString('\n')
		if err != nil {
			fmt.Printf("Got error trying to read incoming data from %s; %v\n", client.RemoteAddr().String(), err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339)
		fmt.Printf("[%v] %s: %s\n", myTime, client.RemoteAddr().String(), temp)
		if temp == "STOP" {
			break
		}

		response := strconv.Itoa(rand.Int()) + "\n"
		client.Write([]byte(string(response)))
	}

	fmt.Println("Client session has closed")
	client.Close()
}
