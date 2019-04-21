package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"grpcTest/midserve/entity"
	"grpcTest/pb"
	"log"
	"time"
)

func main() {
	app1Client()
	app2Client()
}
func app1Client() {
	conn, err := grpc.Dial("127.0.0.1:10001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewBusClient(conn)
	object, err := ptypes.MarshalAny(&entity.UserRequest{
		UserId: "id-1",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	preq := &pb.PipeRequest{
		Qid:    "1",
		Params: object,
	}
	rsp, err := pipe(c, preq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Params.String())
}
func app2Client() {
	conn, err := grpc.Dial("127.0.0.1:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewBusClient(conn)
	object, err := ptypes.MarshalAny(&entity.ShoesRequest{
		Size: 42,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	preq := &pb.PipeRequest{
		Qid:    "1",
		Params: object,
	}
	rsp, err := pipe(c, preq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp.Params.String())
}

func pipe(client pb.BusClient, request *pb.PipeRequest) (*pb.PipeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.Pipe(ctx)
	if err != nil {
		return nil, err
	}
	cli.Send(request)
	return cli.Recv()

}
