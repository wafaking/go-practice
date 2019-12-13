package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	hello "go_practice/grpc_demo/hello/client/proto"
)

const (
	Address = "127.0.0.1:3030"
)

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Errorf("dial err: %s\n", err)
	}
	defer conn.Close()
	log.Printf("dial succ on %s...\n", Address)

	client := hello.NewGreeterClient(conn)

	//FuncSayHello(client)
	//FuncSayMany(client)
	//FuncReplyMany(client)
	FuncTalking(client)
}

func FuncSayHello(client hello.GreeterClient) {
	resp, err := client.SayHello(context.Background(),
		&hello.HelloRequest{
			Name: "wafa",
		},
	)
	if err != nil {
		grpclog.Error("sayhello err: %s", err)
	}

	log.Println("SayHello resp: ", resp.Reply)
}

func FunSayMany(client hello.GreeterClient) {
	var sli []int
	for i := 0; i < 10; i++ {
		sli = append(sli, i)
	}
	stream, err := client.SayMany(context.Background())
	for _, v := range sli {
		var req = &hello.HelloRequest{
			Name: fmt.Sprintf("i am %d", v),
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("send err: %s", err)
			return
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		grpclog.Fatalf("CloseAndRecv: %s", err)
	}
	log.Printf("resp: %s", resp.Reply)
}

func FuncReplyMany(client hello.GreeterClient) {
	req := &hello.HelloRequest{
		Name: "I am client, I send just one time",
	}

	srv, err := client.ReplyMany(context.Background(), req)
	if err != nil {
		grpclog.Fatalf("reply many failed, err: %s", err)
		return
	}
	for {
		resp, err := srv.Recv()
		if err == io.EOF {
			log.Println("resp recv end...")
			break
		} else if err != nil {
			grpclog.Fatalf("stream recv err: %s\n", err)
			break
		}
		log.Println("resp recv processing, resp:  ", resp.Reply)
	}

}

func FuncTalking(client hello.GreeterClient) {
	waitc := make(chan struct{})
	srv, err := client.Talking(context.Background())
	if err != nil {
		grpclog.Fatalf("Taling err: %s\n", err)
		return
	}
	go func() {
		//接收
		for {
			resp, err := srv.Recv()
			if err == io.EOF {
				log.Println("resp recv end...")
				close(waitc) //读取完毕
				break
			} else if err != nil {
				grpclog.Fatalf("stream recv err: %s\n", err)
				break
			}
			log.Println("resp recv processing, resp:  ", resp.Reply)
		}
	}()

	var sli []int
	for i := 0; i < 10; i++ {
		sli = append(sli, i)
	}

	for _, v := range sli {
		//发送
		var req = &hello.HelloRequest{
			Name: fmt.Sprintf("i am %d", v),
		}
		if err := srv.Send(req); err != nil {
			log.Fatalf("send err: %s", err)
		}
	}
	srv.CloseSend()

	<-waitc
}
