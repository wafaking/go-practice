package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	hello "go_practice/grpc_demo/hello/server/proto"
)

const (
	//Address = "http://www.wafa.com"
	Address = "127.0.0.1:3030"
)

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		panic(err)
	}
	log.Printf("listen on: %s\n", Address)

	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, &Controller{})

	if err := s.Serve(listen); err != nil {
		grpclog.Fatalf("serve failed, err: %s\n", err)
	}
}

type Controller struct{}

func (c *Controller) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Printf("hello %s\n", req.Name)
	resp := &hello.HelloResponse{
		Reply: fmt.Sprintf("hello %s, I am %s", req.Name, "nancy"),
	}
	return resp, nil
}

func (c *Controller) SayMany(srv hello.Greeter_SayManyServer) error {
	var resp string
	for {
		req, err := srv.Recv()
		if err == nil {
			resp = fmt.Sprintf("%s, %s", resp, req.Name)
		} else if err == io.EOF {
			log.Printf("I have recieved all data: %s\n", resp)
			return srv.SendAndClose(
				&hello.HelloResponse{
					Reply: fmt.Sprintf("ok, I am receive all data"),
				},
			)
		} else {
			return srv.SendAndClose(
				&hello.HelloResponse{
					Reply: fmt.Sprintf("Sorry, Server is invalid..."),
				},
			)
		}
	}
}

func (c *Controller) ReplyMany(req *hello.HelloRequest, srv hello.Greeter_ReplyManyServer) error {
	for i := 0; i < 10; i++ {
		if err := srv.Send(&hello.HelloResponse{
			Reply: fmt.Sprint("hello %s%d", req.Name, i),
		}); err != nil {
			log.Fatalf("Reply many err: %s", err)
			return err
		}
	}
	return nil
}

func (c *Controller) Talking(srv hello.Greeter_TalkingServer) error {
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			grpclog.Error("recv err: %s", err)
			return err
		}

		var mm sync.Map
		mm.Store(req.Name, fmt.Sprintf("%s ok", req.Name))
		for {
			if v, ok := mm.Load(req.Name); ok {
				if err = srv.Send(
					&hello.HelloResponse{
						Reply: fmt.Sprintf("%s", v),
					}); err != nil {
					grpclog.Error("send err: %s", err)
					return err
				}
			}
		}
	}
}
