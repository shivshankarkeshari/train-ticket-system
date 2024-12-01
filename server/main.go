package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "train-ticket-system/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type server struct {
	pb.UnsafeTrainServiceServer
}



func (s *server) PurchaseTicket(context.Context, *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PurchaseTicket not implemented")
}
func (s *server) GetReceipt(context.Context, *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReceipt not implemented")
}
func (s *server) GetUsersBySection(context.Context, *pb.GetUsersBySectionRequest) (*pb.GetUsersBySectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersBySection not implemented")
}
func (s *server) RemoveUser(context.Context, *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUser not implemented")
}
func (s *server) ModifyUserSeat(context.Context, *pb.ModifyUserSeatRequest) (*pb.ModifyUserSeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyUserSeat not implemented")
}

func main() {

	lis, err:=net.Listen("tcp", ":9001")
	if err!=nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer :=grpc.NewServer()
	pb.RegisterTrainServiceServer(grpcServer, &server{})
	fmt.Println("->>--->")

	if err:=grpcServer.Serve(lis); err!=nil {
		log.Fatalf("failed to serve :%s", err)
	}
	fmt.Println("->>--->")

	
}