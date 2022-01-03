package main

import (
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
		fmt.Println("Please provide a port number!")
		return
	}
	PORT := ":" + arguments[1]

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Printf("Got error trying to listen on port %s; err %v\n", PORT, err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Closing UDP connection")
			return
		}

		data := []byte(strconv.Itoa(rand.Int()))
		fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)

		if err != nil {
			fmt.Printf("Got error trying to send data to %s; err %v", addr.IP.String(), err)
			return
		}
	}
}
