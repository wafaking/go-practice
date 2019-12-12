package main

import (
	"context"
	"fmt"
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
	FuncReplyMany(client)
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
	client.ReplyMany(context.Background())
}

//func (c *greeterClient) ReplyMany(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Greeter_ReplyManyClient, error) {
//	stream, err := c.cc.NewStream(ctx, &_Greeter_serviceDesc.Streams[1], "/hello.Greeter/ReplyMany", opts...)
//	if err != nil {
//		return nil, err
//	}
//	x := &greeterReplyManyClient{stream}
//	if err := x.ClientStream.SendMsg(in); err != nil {
//		return nil, err
//	}
//	if err := x.ClientStream.CloseSend(); err != nil {
//		return nil, err
//	}
//	return x, nil
//}
//
//func (c *greeterClient) Talking(ctx context.Context, opts ...grpc.CallOption) (Greeter_TalkingClient, error) {
//	stream, err := c.cc.NewStream(ctx, &_Greeter_serviceDesc.Streams[2], "/hello.Greeter/Talking", opts...)
//	if err != nil {
//		return nil, err
//	}
//	x := &greeterTalkingClient{stream}
//	return x, nil
//}
