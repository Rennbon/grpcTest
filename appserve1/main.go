package main

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"grpcTest/pb/entity"
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
		pmu := &entity.UserRequest{}
		err = ptypes.UnmarshalAny(in.GetParams(), pmu)
		if err != nil {
			return err
		}

		rsp := &pb.PipeResponse{
			Code: true,
		}
		object := &any.Any{}
		if pmu.UserId == "id-1" {
			pm := &entity.UserInfo{
				UserId:   "id-001",
				UserName: "Jack",
				Face:     "https://upload.jianshu.io/users/upload_avatars/7215015/1d634670-8438-4a49-9748-9fc36f5168ce?imageMogr2/auto-orient/strip|imageView2/1/w/120/h/120",
			}
			object, err = ptypes.MarshalAny(pm)
			if err != nil {
				return err
			}
		} else {
			pm := &entity.Shoes{
				Logo:  "NIKE",
				Size:  43,
				Price: 999.98,
			}
			object, err = ptypes.MarshalAny(pm)
			if err != nil {
				return err
			}
		}
		rsp.Params = object
		if err := stream.Send(rsp); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:10001")
	log.Println("server start")
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
