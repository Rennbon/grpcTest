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
	conn1, err := grpc.Dial("127.0.0.1:10001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	conn2, err := grpc.Dial("127.0.0.1:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()
	defer conn2.Close()
	c1 := pb.NewBusClient(conn1)
	pipe(c1, "id-1")

	c2 := pb.NewBusClient(conn2)
	pipe(c2, "id-2")

}
func pipe(client pb.BusClient, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	object, err := ptypes.MarshalAny(&entity.UserRequest{
		UserId: userId,
	})
	if err != nil {
		return err
	}
	preq := &pb.PipeRequest{
		Qid:    "1",
		Params: object,
	}
	cli, err := client.Pipe(ctx)
	if err != nil {
		return err
	}
	cli.Send(preq)
	for {
		resp, err := cli.Recv()
		if err != nil {
			return err
		}
		pmu := &entity.UserInfo{}
		err = ptypes.UnmarshalAny(resp.Params, pmu)
		if err != nil {
			return err
		}
		fmt.Println(pmu.UserId, pmu.UserName)
		break
	}
	return nil
}
