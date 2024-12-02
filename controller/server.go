package controller

import (
	"fmt"
	"sync"
	pb "train-ticket-system/proto"

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

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		sections: map[string][]*Seat{
			"A": {},
			"B": {},
		},
		users:    make(map[string]*pb.PurchaseResponse),
		nextSeat: 1,
	}
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

// RemoveUser removes a user from the system by their email.
func (s *InMemoryStorage) RemoveUser(email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.users[email]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Remove the seat
	sectionSeats := s.sections[receipt.Section]
	for i, seat := range sectionSeats {
		if seat.Email == email {
			s.sections[receipt.Section] = append(sectionSeats[:i], sectionSeats[i+1:]...)
			break
		}
	}
	// Remove the user
	delete(s.users, email)
	return nil
}

// ModifyUserSeat modifies a user's seat and section.
func (s *InMemoryStorage) ModifyUserSeat(email, newSection string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.users[email]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Remove the seat from the current section
	currentSectionSeats := s.sections[receipt.Section]
	for i, seat := range currentSectionSeats {
		if seat.Email == email {
			s.sections[receipt.Section] = append(currentSectionSeats[:i], currentSectionSeats[i+1:]...)
			break
		}
	}

	// Assign a new seat
	seatNumber := s.nextSeat
	s.nextSeat++
	newSeat := &Seat{
		SeatNumber: seatNumber,
		Email:      email,
	}

	// Add the seat to the new section
	s.sections[newSection] = append(s.sections[newSection], newSeat)

	// Update the user's receipt
	receipt.Section = newSection
	receipt.Seat = fmt.Sprintf("%s-%d", newSection, seatNumber)

	return nil
}
