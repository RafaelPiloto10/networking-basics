package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	client "github.com/RafaelPiloto10/grpc_chat/client/pkg"
	"github.com/RafaelPiloto10/grpc_chat/chatpb"
	"google.golang.org/grpc"
)

var id = flag.Uint("id", 0, "The username to join the chat with")
var addr = flag.String("addr", "0.0.0.0:42069", "The server port to connect to")

func main() {
	flag.Parse()
	
	fmt.Println("Booting client app....")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial(*addr, opts...)
	
	if err != nil {
		log.Fatalf("could not connect to server; err %v", err)
	}

	defer conn.Close()

	ctx := context.Background()
	c := chatpb.NewServerClient(conn)
	acc := &chatpb.Account{
		Id: uint32(*id),
		Channel: "",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "!") {
			tokens := strings.Split(msg, " ")
			cmd := tokens[0]
			if cmd == "!join" {
				acc.Channel = tokens[1]
				go client.JoinChannel(ctx, acc, c)
			} else if cmd == "!quit" {
				go client.LeaveChannel(ctx, acc, c)
			} else {
				fmt.Println("error: Unrecognized command! Use (!join <channel> or !quit <channel>)")
			}
		} else {
			go client.SendMessage(ctx, &chatpb.Message{
				ReceiverId: 1000,
				SenderId: acc.GetId(),
				Body: msg,
				ChannelName: acc.GetChannel(),
			}, c)
		}
		fmt.Printf("> ")
	}

}
