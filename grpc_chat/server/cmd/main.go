package main

import (
	"flag"
	"fmt"

	server "github.com/RafaelPiloto10/grpc_chat/server/pkg"
)

var port = flag.Uint("port", 42069, "The port to host the server on")

func main() {	
	flag.Parse()

	fmt.Println("Booting server...")
	s := server.NewServer()
	fmt.Println("Serving server...")
	s.Serve(uint32(*port))
}
