package main

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"grpcTest/appserve2/entity"

	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"grpcTest/pb"
)

type Serve1 struct {
}

func (c *Serve1) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return nil, nil
}
func (c *Serve1) Pipe(stream pb.Bus_PipeServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		pmu := &entity.ShoesRequest{}
		err = ptypes.UnmarshalAny(in.GetParams(), pmu)
		if err != nil {
			return err
		}

		rsp := &pb.PipeResponse{
			Code: true,
		}
		object := &any.Any{}
		if pmu.Size == 42 {
			pm := &entity.Shoes{
				Logo:  "NIKE",
				Size:  42,
				Price: 999.98,
			}
			object, err = ptypes.MarshalAny(pm)
			if err != nil {
				return err
			}
		} else {
			return errors.New("this size doesn't exist")
		}
		rsp.Params = object
		if err := stream.Send(rsp); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:10002")
	log.Println("server2 start")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	/* 	creds, err := credentials.NewServerT	LSFromFile(ca.CertFile, ca.KeyFile)
	   	if err != nil {
	   		log.Fatalf("could not load keys: %s", err)
	   	}
	   	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)} */

	s := grpc.NewServer()
	pb.RegisterBusServer(s, &Serve1{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
