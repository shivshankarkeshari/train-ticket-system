package main

import (
	"context"
	"log"
	pb "train-ticket-system/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err!=nil {
		log.Fatalf("failed to connect to gRPC server")
	}
	defer conn.Close()

	client:=pb.NewTrainServiceClient(conn)

	// Purchase a ticket
	log.Println("Purchasing a ticket...")
	user := &pb.User{
		FirstName: "Shiv",
		LastName:  "Shankar",
		Email:     "shiv.s.keshari@.com",
	}
	purchaseResp, err := client.PurchaseTicket(context.Background(), &pb.PurchaseRequest{
		From:    "London",
		To:      "France",
		User:    user,
		Section: "A",
	})
	if err != nil {
		log.Fatalf("Error purchasing ticket: %v", err)
	}
	log.Printf("Ticket purchased: %+v\n", purchaseResp)

	// View receipt
	log.Println("Fetching receipt...")
	receiptResp, err := client.GetReceipt(context.Background(), &pb.GetReceiptRequest{
		ReceiptId: purchaseResp.ReceiptId,
	})
	if err != nil {
		log.Fatalf("Error fetching receipt: %v", err)
	}
	log.Printf("Receipt details: %+v\n", receiptResp)

	// Get users by section
	log.Println("Getting users in section A...")
	usersResp, err := client.GetUsersBySection(context.Background(), &pb.GetUsersBySectionRequest{
		Section: "A",
	})
	if err != nil {
		log.Fatalf("Error fetching users by section: %v", err)
	}
	log.Printf("Users in section A: %v\n", usersResp.Users)

	// Modify user's seat
	log.Println("Modifying user's seat to section B...")
	modifyResp, err := client.ModifyUserSeat(context.Background(), &pb.ModifyUserSeatRequest{
		Email:      user.Email,
		NewSection: "B",
	})
	if err != nil {
		log.Fatalf("Error modifying user seat: %v", err)
	}
	log.Printf("Modify seat response: %s\n", modifyResp.Message)

	// Verify user is now in section B
	log.Println("Getting users in section B...")
	sectionBResp, err := client.GetUsersBySection(context.Background(), &pb.GetUsersBySectionRequest{
		Section: "B",
	})
	if err != nil {
		log.Fatalf("Error fetching users by section: %v", err)
	}
	log.Printf("Users in section B: %v\n", sectionBResp.Users)

	// Remove user from train
	log.Println("Removing user from train...")
	removeResp, err := client.RemoveUser(context.Background(), &pb.RemoveUserRequest{
		Email: user.Email,
	})
	if err != nil {
		log.Fatalf("Error removing user: %v", err)
	}
	log.Printf("Remove user response: %s\n", removeResp.Message)
	
}
