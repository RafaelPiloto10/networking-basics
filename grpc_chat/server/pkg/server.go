package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/RafaelPiloto10/grpc_chat/chatpb"
	"google.golang.org/grpc"
)

type Server struct {
	chatpb.UnimplementedServerServer
	sync.RWMutex

	chats    map[string][]chan *chatpb.Message
	accounts map[uint32]*ChatAccount
}

type ChatAccount struct {
	ctx     context.Context
	cancel  context.CancelFunc
	channel string
	stream  chatpb.Server_JoinChannelServer
}

func NewServer() *Server {
	return &Server{
		chats:    make(map[string][]chan *chatpb.Message),
		accounts: make(map[uint32]*ChatAccount),
	}
}

func (s *Server) Serve(port uint32) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))

	if err != nil {
		log.Fatalf("Could not listen on port %v; err %v", port, err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatpb.RegisterServerServer(grpcServer, s)
	grpcServer.Serve(lis)
}

func (s *Server) JoinChannel(account *chatpb.Account, msgStream chatpb.Server_JoinChannelServer) error {
	s.Lock()

	ctx, cancel := context.WithCancel(context.Background())

	s.accounts[account.GetId()] = &ChatAccount{
		ctx:     ctx,
		cancel:  cancel,
		channel: account.GetChannel(),
		stream:  msgStream,
	}

	if _, ok := s.chats[account.GetChannel()]; !ok {
		s.chats[account.GetChannel()] = []chan *chatpb.Message{}
	}

	c := make(chan *chatpb.Message)
	s.chats[account.GetChannel()] = append(s.chats[account.GetChannel()], c)

	fmt.Printf("[Server] Account %v has joined %v\n", account.GetId(), account.GetChannel())

	s.Unlock()

	for {
		select {
		case <-msgStream.Context().Done():
			fmt.Printf("Stream for %v has closed!\n", account.GetChannel())
			return nil
		case <-ctx.Done():
			fmt.Printf("Context for %v was canceled!\n", account.GetChannel())
			return nil
		case msg := <-c:
			fmt.Printf("[Server] Got message: %v \n", msg)
			msgStream.Send(msg)
		}
	}
}

func (s *Server) LeaveChannel(ctx context.Context, account *chatpb.Account) (*chatpb.Empty, error) {
	s.Lock()
	defer s.Unlock()
	fmt.Printf("[Server] User %v has left %v\n", account.GetId(), account.GetChannel())
	s.accounts[account.GetId()].cancel()
	delete(s.accounts, account.GetId())
	return &chatpb.Empty{}, nil
}

func (s *Server) SendMessage(ctx context.Context, msg *chatpb.Message) (*chatpb.Response, error) {
	s.RLock()
	defer s.RUnlock()

	fmt.Println("[Server] Sending message")

	if channels, ok := s.chats[s.accounts[msg.GetSenderId()].channel]; ok {
		for _, channel := range channels {
			channel <- msg
		}

		return &chatpb.Response{
			Status: chatpb.Response_OK,
		}, nil

	} else {
		return &chatpb.Response{
			Status: chatpb.Response_ERROR,
		}, fmt.Errorf("User: %v is undefined or not in a channel\n", msg.GetSenderId())
	}
}
