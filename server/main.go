package main


import (
	"context"
	"fmt"
	"net"
	"log"
	"google.golang.org/grpc"
	pb "train-ticket-system/proto"
	controller "train-ticket-system/controller"

)

type server struct {
	pb.UnsafeTrainServiceServer
}

var storage = controller.NewInMemoryStorage()

func (s *server) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	fmt.Println("->>---> PurchaseTicket")
	receipt := storage.AddUser(req.User, req.Section, req.From, req.To)
	return receipt, nil
}

func (s *server) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.GetReceiptResponse, error) {
	fmt.Println("->>---> GetReceipt")
	receiptID := req.ReceiptId
	return storage.GetReceipt(receiptID)
}

func (s *server) GetUsersBySection(ctx context.Context, req *pb.GetUsersBySectionRequest) (*pb.GetUsersBySectionResponse, error) {
	fmt.Println("->>---> GetUsersBySection")
	section := req.Section
	users, err := storage.GetUsersBySection(section)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users for section %s: %v", section, err)
	}
	return &pb.GetUsersBySectionResponse{Users: users}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	fmt.Println("->>---> RemoveUser")
	err:=storage.RemoveUser(req.Email)
	if err!=nil {
		return nil, err
	}
	return &pb.RemoveUserResponse{Message: "user removed sucessfully"}, nil
}

func (s *server) ModifyUserSeat(ctx context.Context, req *pb.ModifyUserSeatRequest) (*pb.ModifyUserSeatResponse, error) {
	fmt.Println("->>---> ModifyUserSeat")
	err := storage.ModifyUserSeat(req.Email, req.NewSection)
	if err!=nil {
		return nil, err
	}
	return &pb.ModifyUserSeatResponse{Message: "seat got changed"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterTrainServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve :%s", err)
	}
	fmt.Println("->>--->")

}
