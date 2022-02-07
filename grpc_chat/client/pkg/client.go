package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	chatpb "github.com/RafaelPiloto10/grpc_chat/chatpb"
)

func JoinChannel(ctx context.Context, account *chatpb.Account, client chatpb.ServerClient) {
	stream, err := client.JoinChannel(ctx, account)

	if err != nil {
		log.Printf("got error trying to join channel %v; %v", account.GetChannel(), err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		for {
			in, err := stream.Recv()

			if err == io.EOF {
				wg.Done()
				return
			}

			if err != nil {
				log.Fatalf("[Client] Failed to receive message from channel %v; %v",account.GetChannel(),  err)
			}

			if account.GetId() != in.GetSenderId() {
				fmt.Printf("[%v] (%v) -> %v \n", time.Now().String(), in.GetSenderId(), in.GetBody())
			}
		}
	}()

	wg.Wait()
	fmt.Printf("[Client] Stopped listening to stream on channel %v", account.GetChannel())
}

func LeaveChannel(ctx context.Context, account *chatpb.Account, client chatpb.ServerClient) {
	_, err := client.LeaveChannel(ctx, account)

	if err != nil {
		fmt.Printf("got error trying to leave channel; %v", err)
	}
}

func SendMessage(ctx context.Context, message *chatpb.Message, client chatpb.ServerClient) {
	response, err := client.SendMessage(ctx, message)	
	if err != nil || response.Status == chatpb.Response_ERROR {
		fmt.Printf("got error trying to send message; %v", err)
	}
	fmt.Printf("[Client] Sent message!\n")
}
