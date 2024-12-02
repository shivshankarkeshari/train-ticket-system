package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	pb "train-ticket-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type InMemoryStorage struct {
	mu       sync.Mutex
	sections map[string][]*Seat 
	users    map[string]*pb.PurchaseResponse 
	nextSeat int
}

type Seat struct {
	SeatNumber int
	Email      string
}

var storage = &InMemoryStorage{
	sections: map[string][]*Seat{
		"A": {},
		"B": {},
	},
	users:    make(map[string]*pb.PurchaseResponse),
	nextSeat: 1,
}


// AddUser allocates a seat and stores user information.
func (s *InMemoryStorage) AddUser(user *pb.User, section, from, to string) *pb.PurchaseResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Assign a seat
	seatNumber := s.nextSeat
	s.nextSeat++

	seat := &Seat{
		SeatNumber: seatNumber,
		Email:      user.Email,
	}

	// Add seat to the section
	s.sections[section] = append(s.sections[section], seat)

	// Create receipt
	receipt := &pb.PurchaseResponse{
		ReceiptId: fmt.Sprintf("R-%d", seatNumber),
		From:      from,
		To:        to,
		User:      user,
		Section:   section,
		Seat:      fmt.Sprintf("%s-%d", section, seatNumber),
		PricePaid: 20.0, // Fixed price
	}

	// Store user receipt
	s.users[user.Email] = receipt
	return receipt
}

// GetReceipt retrieves a user's receipt by their email.
func (s *InMemoryStorage) GetReceipt(receiptID string) (*pb.GetReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, receipt := range s.users {
		if receipt.ReceiptId == receiptID {
			response := &pb.GetReceiptResponse{
				From:      receipt.From,
				To:        receipt.To,
				User:      receipt.User,
				Section:   receipt.Section,
				Seat:      receipt.Seat,
				PricePaid: receipt.PricePaid,
			}
		
			return response, nil
		}
	}
	return nil, fmt.Errorf("receipt not found")
}


// GetUsersBySection retrieves all users in a given section.
func (s *InMemoryStorage) GetUsersBySection(section string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sections[section]; !ok {
		return nil, fmt.Errorf("section not found")
	}

	users := []string{}
	for _, seat := range s.sections[section] {
		users = append(users, seat.Email)
	}

	return users, nil
}

type server struct {
	pb.UnsafeTrainServiceServer
}

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
	// Build the response.
	response := &pb.GetUsersBySectionResponse{
		Users: users,
	}
	return response, nil
}
func (s *server) RemoveUser(context.Context, *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUser not implemented")
}
func (s *server) ModifyUserSeat(context.Context, *pb.ModifyUserSeatRequest) (*pb.ModifyUserSeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyUserSeat not implemented")
}

// InMemoryStorage storage := NewInMemoryStorage()


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
